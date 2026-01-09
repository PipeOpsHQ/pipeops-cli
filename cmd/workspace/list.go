package workspace

import (
	"context"
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the command to list all workspaces
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all workspaces",
	Long: `List all workspaces in your PipeOps account.

Examples:
  pipeops workspace list
  pipeops workspace ls
  pipeops workspace list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		client := pipeops.NewClientWithConfig(cfg)

		if !utils.RequireAuth(client, opts) {
			return
		}

		utils.PrintInfo("Fetching workspaces...", opts)

		workspaces, err := client.GetWorkspaces(context.Background())
		if err != nil {
			if !utils.HandleAuthError(err, opts) {
				return
			}
			utils.HandleError(err, "Error fetching workspaces", opts)
			return
		}

		if len(workspaces) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON([]interface{}{})
			} else {
				fmt.Println("No workspaces found")
			}
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(workspaces)
		} else {
			headers := []string{"UUID", "NAME", "DESCRIPTION", "CREATED"}
			var rows [][]string

			for _, ws := range workspaces {
				name := utils.TruncateString(ws.Name, 30)
				desc := utils.TruncateString(ws.Description, 40)
				created := ""
				if ws.CreatedAt != nil {
					created = utils.FormatDateShort(ws.CreatedAt.Time)
				}

				// Mark current workspace
				if cfg.Settings != nil && cfg.Settings.DefaultWorkspaceUUID == ws.UUID {
					name = "✓ " + name
				}

				rows = append(rows, []string{
					ws.UUID,
					name,
					desc,
					created,
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d workspaces", len(workspaces)), opts)

			if cfg.Settings != nil && cfg.Settings.DefaultWorkspaceUUID != "" {
				fmt.Printf("\n💡 Current workspace: %s\n", cfg.Settings.DefaultWorkspaceUUID)
			} else {
				fmt.Printf("\n💡 TIP: Select a default workspace with: pipeops workspace select\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

func (w *workspaceModel) list() {
	listCmd.Flags().Bool("json", false, "Output in JSON format")
	w.rootCmd.AddCommand(listCmd)
}
