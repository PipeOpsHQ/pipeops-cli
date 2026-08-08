package project

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

// buildLogsCmd fetches Firebase pipeops-build-logs via the control-plane API
// (same source as the dashboard Build Logs tab). Distinct from runtime
// "project logs" (application stdout).
var buildLogsCmd = &cobra.Command{
	Use:   "build-logs <project-id>",
	Short: "View Firebase deployment build logs for a project",
	Long: `Fetch deployment build logs from Firebase (pipeops-build-logs), the same
source as the PipeOps dashboard Build Logs tab.

Defaults to the latest deployment when --deployment / --build-sha are omitted.

Examples:
  - Latest deployment build logs:
    pipeops project build-logs proj-uuid

  - Specific deployment:
    pipeops project build-logs proj-uuid --deployment dep-uuid

  - Filter to the build stage only:
    pipeops project build-logs proj-uuid --stage build

  - JSON for scripting:
    pipeops project build-logs proj-uuid --json
`,
	Aliases: []string{"buildlogs", "build_logs"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}

		projectID := strings.TrimSpace(args[0])
		if projectID == "" {
			return fmt.Errorf("project-id is required")
		}

		limit := intFlag(cmd, "limit", 2000)
		if limit <= 0 {
			limit = 2000
		}
		if limit > 5000 {
			limit = 5000
		}

		stage, _ := cmd.Flags().GetString("stage")
		stage = strings.ToLower(strings.TrimSpace(stage))
		if stage != "" && stage != "git" && stage != "build" && stage != "deploy" {
			return fmt.Errorf("invalid --stage %q (want git, build, or deploy)", stage)
		}

		deploymentUUID, _ := cmd.Flags().GetString("deployment")
		buildSha, _ := cmd.Flags().GetString("build-sha")

		// WorkspaceUUID is filled by the client only when --workspace / env / default is set.
		resp, err := client.GetBuildLogs(projectID, &sdk.BuildLogsOptions{
			DeploymentUUID: strings.TrimSpace(deploymentUUID),
			BuildSha:       strings.TrimSpace(buildSha),
			Stage:          stage,
			Limit:          limit,
		})
		if err != nil {
			return fmt.Errorf("get build logs: %w", err)
		}

		return printBuildLogs(resp, opts)
	},
}

func printBuildLogs(resp *sdk.BuildLogsResponse, opts utils.OutputOptions) error {
	if resp == nil {
		return fmt.Errorf("empty build logs response")
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(resp)
	}

	data := resp.Data
	utils.PrintInfo(fmt.Sprintf(
		"Build logs project=%s deployment=%s build_sha=%s status=%s stage=%s source=%s count=%d",
		data.ProjectUUID,
		data.DeploymentUUID,
		data.BuildSha,
		data.Status,
		data.CurrentStage,
		data.Source,
		data.Count,
	), opts)

	if len(data.Logs) == 0 {
		utils.PrintWarning("No build log lines returned", opts)
		return nil
	}

	for _, line := range data.Logs {
		fmt.Println(formatBuildLogLine(line))
	}
	utils.PrintSuccess(fmt.Sprintf("Printed %d build log line(s)", len(data.Logs)), opts)
	return nil
}

func formatBuildLogLine(line map[string]interface{}) string {
	if line == nil {
		return ""
	}
	msg := firstString(line, "message", "log", "msg", "text")
	if msg == "" {
		msg = fmt.Sprintf("%v", line)
	}
	ts := firstString(line, "ts", "timestamp", "time", "created_at", "createdAt")
	stage := firstString(line, "stage", "current-stage", "currentStage")
	level := firstString(line, "level", "severity")

	var b strings.Builder
	if ts != "" {
		b.WriteString(ts)
		b.WriteString(" ")
	}
	if stage != "" {
		b.WriteString("[")
		b.WriteString(stage)
		b.WriteString("] ")
	}
	if level != "" {
		b.WriteString(strings.ToUpper(level))
		b.WriteString(" ")
	}
	b.WriteString(msg)
	return b.String()
}

func firstString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			s := strings.TrimSpace(fmt.Sprint(v))
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

func init() {
	buildLogsCmd.Flags().String("deployment", "", "Deployment UUID (default: latest for project)")
	buildLogsCmd.Flags().String("build-sha", "", "Build SHA override (default: from deployment)")
	buildLogsCmd.Flags().String("stage", "", "Filter stage: git, build, or deploy")
	buildLogsCmd.Flags().Int("limit", 2000, "Max log lines (max 5000)")
	buildLogsCmd.Flags().String("workspace", "", workspaceFlagHelp)
}
