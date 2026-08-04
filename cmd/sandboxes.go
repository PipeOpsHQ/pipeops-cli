package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

var sandboxesCmd = &cobra.Command{
	Use:     "sandboxes",
	Aliases: []string{"sandbox", "sb"},
	Short:   "Manage Rexec sandboxes (PipeOps BFF)",
	Long: `Manage sandboxes via the PipeOps controller BFF (proxied to Rexec).

Examples:
  pipeops sandboxes list --workspace <uuid>
  pipeops sandboxes get <id> --workspace <uuid>
  pipeops sandboxes create --name dev --image ubuntu --workspace <uuid>
  pipeops sandboxes start <id> --workspace <uuid>
  pipeops sandboxes stop <id> --workspace <uuid>
  pipeops sandboxes restart <id> --workspace <uuid>
  pipeops sandboxes delete <id> --yes --workspace <uuid>
  pipeops sandboxes session <id> --workspace <uuid>
  pipeops sandboxes usage --from 2026-08-01 --to 2026-08-04 --workspace <uuid>`,
}

func sandboxWorkspaceOpts(cmd *cobra.Command) *sdk.SandboxWorkspaceOptions {
	opts := &sdk.SandboxWorkspaceOptions{}
	if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
		opts.WorkspaceUUID = workspace
	}
	return opts
}

var sandboxesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List sandboxes",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.ListSandboxes(context.Background(), sandboxWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("list sandboxes: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		if len(resp.Data) == 0 {
			utils.PrintWarning("No sandboxes found", opts)
			return nil
		}
		rows := make([][]string, 0, len(resp.Data))
		for _, s := range resp.Data {
			rows = append(rows, []string{s.ID, s.Name, s.Image, s.Role, s.Status})
		}
		utils.PrintTable([]string{"ID", "NAME", "IMAGE", "ROLE", "STATUS"}, rows, opts)
		if !opts.Quiet {
			utils.PrintSuccess(fmt.Sprintf("Found %d sandboxes", len(resp.Data)), opts)
		}
		return nil
	},
	Args: cobra.NoArgs,
}

var sandboxesGetCmd = &cobra.Command{
	Use:   "get <sandbox-id>",
	Short: "Get sandbox details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		box, err := client.GetSandbox(context.Background(), args[0], sandboxWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("get sandbox: %w", err)
		}
		return printSandbox(box, opts)
	},
	Args: cobra.ExactArgs(1),
}

var sandboxesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a sandbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		role, _ := cmd.Flags().GetString("role")
		resp, err := client.CreateSandbox(context.Background(), sandboxWorkspaceOpts(cmd), &sdk.CreateSandboxRequest{
			Name:  name,
			Image: image,
			Role:  role,
		})
		if err != nil {
			return fmt.Errorf("create sandbox: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("Sandbox created", opts)
		return printSandbox(&resp.Data, opts)
	},
	Args: cobra.NoArgs,
}

var sandboxesStartCmd = &cobra.Command{
	Use:   "start <sandbox-id>",
	Short: "Start a sandbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sandboxAction(cmd, args[0], "start", func(c pipeops.ClientAPI, id string, o *sdk.SandboxWorkspaceOptions) error {
			return c.StartSandbox(context.Background(), id, o)
		})
	},
	Args: cobra.ExactArgs(1),
}

var sandboxesStopCmd = &cobra.Command{
	Use:   "stop <sandbox-id>",
	Short: "Stop a sandbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sandboxAction(cmd, args[0], "stop", func(c pipeops.ClientAPI, id string, o *sdk.SandboxWorkspaceOptions) error {
			return c.StopSandbox(context.Background(), id, o)
		})
	},
	Args: cobra.ExactArgs(1),
}

var sandboxesRestartCmd = &cobra.Command{
	Use:   "restart <sandbox-id>",
	Short: "Restart a sandbox (stop then start)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sandboxAction(cmd, args[0], "restart", func(c pipeops.ClientAPI, id string, o *sdk.SandboxWorkspaceOptions) error {
			return c.RestartSandbox(context.Background(), id, o)
		})
	},
	Args: cobra.ExactArgs(1),
}

var sandboxesDeleteCmd = &cobra.Command{
	Use:   "delete <sandbox-id>",
	Short: "Delete a sandbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			return fmt.Errorf("--yes is required to delete a sandbox")
		}
		client, err := rootClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteSandbox(context.Background(), args[0], sandboxWorkspaceOpts(cmd)); err != nil {
			return fmt.Errorf("delete sandbox: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "sandbox_id": args[0]})
		}
		utils.PrintSuccess("Sandbox deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var sandboxesSessionCmd = &cobra.Command{
	Use:   "session <sandbox-id>",
	Short: "Create a short-lived terminal session grant",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		sess, err := client.CreateSandboxSession(context.Background(), args[0], sandboxWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("create sandbox session: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(sess)
		}
		utils.PrintSuccess("Session grant created (store token securely; do not log)", opts)
		utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
			{"Container ID", sess.ContainerID},
			{"Base URL", sess.BaseURL},
			{"Token", sess.Token},
			{"Expires In (s)", fmt.Sprintf("%d", sess.ExpiresIn)},
			{"Token Source", sess.TokenSource},
			{"Grant ID", sess.GrantID},
		}, opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var sandboxesUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show daily sandbox usage for a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		fromStr, _ := cmd.Flags().GetString("from")
		toStr, _ := cmd.Flags().GetString("to")
		var from, to time.Time
		if fromStr != "" {
			from, err = time.Parse("2006-01-02", fromStr)
			if err != nil {
				return fmt.Errorf("invalid --from (want YYYY-MM-DD): %w", err)
			}
		}
		if toStr != "" {
			to, err = time.Parse("2006-01-02", toStr)
			if err != nil {
				return fmt.Errorf("invalid --to (want YYYY-MM-DD): %w", err)
			}
		}
		resp, err := client.SandboxUsageDaily(context.Background(), sandboxWorkspaceOpts(cmd), from, to)
		if err != nil {
			return fmt.Errorf("sandbox usage: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		if len(resp.Data) == 0 {
			utils.PrintWarning("No usage rows", opts)
			return nil
		}
		rows := make([][]string, 0, len(resp.Data))
		for _, u := range resp.Data {
			day := ""
			if u.Day != nil {
				day = u.Day.Format("2006-01-02")
			}
			rows = append(rows, []string{
				day,
				fmt.Sprintf("%d", u.CreatedCount),
				fmt.Sprintf("%d", u.StartedCount),
				fmt.Sprintf("%d", u.SessionCount),
				fmt.Sprintf("%d", u.TotalDurationSeconds),
				fmt.Sprintf("%d", u.UniqueContainers),
			})
		}
		utils.PrintTable([]string{"DAY", "CREATED", "STARTED", "SESSIONS", "DURATION_S", "UNIQUE"}, rows, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

func sandboxAction(cmd *cobra.Command, id, verb string, fn func(pipeops.ClientAPI, string, *sdk.SandboxWorkspaceOptions) error) error {
	opts := utils.GetOutputOptions(cmd)
	client, err := rootClient(cmd, opts)
	if err != nil || client == nil {
		return err
	}
	if err := fn(client, id, sandboxWorkspaceOpts(cmd)); err != nil {
		return fmt.Errorf("%s sandbox: %w", verb, err)
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(map[string]string{"status": verb + "ed", "sandbox_id": id})
	}
	utils.PrintSuccess(fmt.Sprintf("Sandbox %sed", verb), opts)
	return nil
}

func printSandbox(box *sdk.Sandbox, opts utils.OutputOptions) error {
	if box == nil {
		return fmt.Errorf("sandbox not found")
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(box)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"ID", box.ID},
		{"UUID", box.UUID},
		{"Name", box.Name},
		{"Image", box.Image},
		{"Role", box.Role},
		{"Status", box.Status},
	}, opts)
	return nil
}

func init() {
	workspaceFlag := "Workspace UUID (or set PIPEOPS_WORKSPACE_UUID / pipeops workspace select)"

	for _, c := range []*cobra.Command{
		sandboxesListCmd, sandboxesGetCmd, sandboxesCreateCmd,
		sandboxesStartCmd, sandboxesStopCmd, sandboxesRestartCmd,
		sandboxesDeleteCmd, sandboxesSessionCmd, sandboxesUsageCmd,
	} {
		c.Flags().String("workspace", "", workspaceFlag)
	}

	sandboxesCreateCmd.Flags().String("name", "", "Sandbox name")
	sandboxesCreateCmd.Flags().String("image", "ubuntu", "Container image")
	sandboxesCreateCmd.Flags().String("role", "standard", "Role/profile (e.g. standard)")

	sandboxesDeleteCmd.Flags().Bool("yes", false, "Confirm sandbox deletion")

	sandboxesUsageCmd.Flags().String("from", "", "Start day YYYY-MM-DD")
	sandboxesUsageCmd.Flags().String("to", "", "End day YYYY-MM-DD")

	sandboxesCmd.AddCommand(
		sandboxesListCmd,
		sandboxesGetCmd,
		sandboxesCreateCmd,
		sandboxesStartCmd,
		sandboxesStopCmd,
		sandboxesRestartCmd,
		sandboxesDeleteCmd,
		sandboxesSessionCmd,
		sandboxesUsageCmd,
	)
	rootCmd.AddCommand(sandboxesCmd)
}
