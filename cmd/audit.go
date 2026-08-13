package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:     "audit",
	Aliases: []string{"audit-logs", "activity"},
	Short:   "Query historical PipeOps actions (audit logs)",
	Long: `List project- and workspace-scoped audit history (who redeployed,
changed env, paused, agent/webhook deploys, etc.).

Uses the control-plane audit APIs:
  GET /project/audit-logs/:uuid
  GET /project/workspace-audit-logs?workspace_uuid=

Examples:
  pipeops audit project <project-uuid>
  pipeops audit project <project-uuid> --action project.redeploy --limit 50
  pipeops audit workspace
  pipeops audit workspace --workspace <ws-uuid> --from 2026-08-01T00:00:00Z
  pipeops audit workspace --project <project-uuid> --actor-type agent
  pipeops audit list --json
`,
}

func auditCommonFlags(cmd *cobra.Command) {
	cmd.Flags().String("action", "", "Action filter (comma-separated), e.g. project.redeploy,project.env.update")
	cmd.Flags().String("actor-type", "", "Actor type: user, webhook, system, service_account, agent")
	cmd.Flags().String("actor-user-uuid", "", "Filter by actor user UUID")
	cmd.Flags().String("category", "", "Category: lifecycle, settings, deployment, security, access")
	cmd.Flags().String("search", "", "Free-text search over summary / names")
	cmd.Flags().String("from", "", "Start time (RFC3339)")
	cmd.Flags().String("to", "", "End time (RFC3339)")
	cmd.Flags().Int("limit", 20, "Page size (default 20)")
	cmd.Flags().Int("offset", 0, "Pagination offset")
}

func auditProjectOpts(cmd *cobra.Command) *sdk.ProjectAuditLogListOptions {
	opts := &sdk.ProjectAuditLogListOptions{}
	opts.Action, _ = cmd.Flags().GetString("action")
	opts.ActorType, _ = cmd.Flags().GetString("actor-type")
	opts.ActorUserUUID, _ = cmd.Flags().GetString("actor-user-uuid")
	opts.Category, _ = cmd.Flags().GetString("category")
	opts.Search, _ = cmd.Flags().GetString("search")
	opts.From, _ = cmd.Flags().GetString("from")
	opts.To, _ = cmd.Flags().GetString("to")
	if n, err := cmd.Flags().GetInt("limit"); err == nil && n > 0 {
		opts.Limit = n
	}
	if n, err := cmd.Flags().GetInt("offset"); err == nil && n >= 0 {
		opts.Offset = n
	}
	return opts
}

func auditWorkspaceOpts(cmd *cobra.Command) *sdk.WorkspaceAuditLogListOptions {
	p := auditProjectOpts(cmd)
	opts := &sdk.WorkspaceAuditLogListOptions{
		Action:        p.Action,
		ActorType:     p.ActorType,
		ActorUserUUID: p.ActorUserUUID,
		Category:      p.Category,
		Search:        p.Search,
		From:          p.From,
		To:            p.To,
		Limit:         p.Limit,
		Offset:        p.Offset,
	}
	if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
		opts.WorkspaceUUID = strings.TrimSpace(workspace)
	}
	if project, _ := cmd.Flags().GetString("project"); project != "" {
		opts.ProjectUUID = strings.TrimSpace(project)
	}
	return opts
}

var auditProjectCmd = &cobra.Command{
	Use:     "project <project-uuid>",
	Aliases: []string{"proj"},
	Short:   "List audit logs for a project",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		out := utils.GetOutputOptions(cmd)
		client, err := rootClient(cmd, out)
		if err != nil || client == nil {
			return err
		}
		projectUUID := strings.TrimSpace(args[0])
		if projectUUID == "" {
			return fmt.Errorf("project-uuid is required")
		}
		resp, err := client.ListProjectAuditLogs(context.Background(), projectUUID, auditProjectOpts(cmd))
		if err != nil {
			return fmt.Errorf("list project audit logs: %w", err)
		}
		return printAuditLogs(resp.Data, resp.Pagination, out, "project")
	},
}

var auditWorkspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws", "list", "ls"},
	Short:   "List audit logs for a workspace (all projects)",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := utils.GetOutputOptions(cmd)
		client, err := rootClient(cmd, out)
		if err != nil || client == nil {
			return err
		}
		opts := auditWorkspaceOpts(cmd)
		if ws, _ := cmd.Flags().GetString("workspace"); ws != "" {
			client.SetWorkspaceOverride(strings.TrimSpace(ws))
		}
		resp, err := client.ListWorkspaceAuditLogs(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("list workspace audit logs: %w", err)
		}
		return printAuditLogs(resp.Data, resp.Pagination, out, "workspace")
	},
}

func printAuditLogs(logs []sdk.ProjectAuditLog, page sdk.AuditLogPagination, opts utils.OutputOptions, scope string) error {
	envelope := map[string]interface{}{
		"scope":      scope,
		"data":       logs,
		"pagination": page,
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(envelope)
	}
	if len(logs) == 0 {
		utils.PrintWarning("No audit log entries found", opts)
		return nil
	}
	rows := make([][]string, 0, len(logs))
	for _, log := range logs {
		when := formatAuditTime(log.CreatedAt)
		actor := log.Actor.Name
		if actor == "" {
			actor = log.Actor.Label
		}
		if actor == "" {
			actor = log.Actor.Type
		}
		label := log.ActionLabel
		if label == "" {
			label = log.Action
		}
		rows = append(rows, []string{
			when,
			label,
			log.Status,
			actor,
			displayOr(log.ProjectName, log.ProjectUUID),
			truncateAudit(log.Summary, 60),
		})
	}
	utils.PrintTable([]string{"WHEN", "ACTION", "STATUS", "ACTOR", "PROJECT", "SUMMARY"}, rows, opts)
	if !opts.Quiet {
		utils.PrintSuccess(fmt.Sprintf(
			"Showing %d of %d audit entries (limit=%d offset=%d)",
			len(logs), page.Total, page.Limit, page.Offset,
		), opts)
	}
	return nil
}

func formatAuditTime(ts *sdk.Timestamp) string {
	if ts == nil || ts.Time.IsZero() {
		return ""
	}
	return ts.Time.UTC().Format("2006-01-02 15:04")
}

func truncateAudit(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func init() {
	workspaceFlag := "Workspace UUID (or set PIPEOPS_WORKSPACE_UUID / pipeops workspace select)"

	auditCommonFlags(auditProjectCmd)
	auditCommonFlags(auditWorkspaceCmd)
	auditWorkspaceCmd.Flags().String("workspace", "", workspaceFlag)
	auditWorkspaceCmd.Flags().String("project", "", "Optional project UUID to filter within the workspace")

	auditCmd.AddCommand(auditProjectCmd, auditWorkspaceCmd)
	rootCmd.AddCommand(auditCmd)
}
