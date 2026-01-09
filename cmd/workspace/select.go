package workspace

import (
	"context"
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// selectCmd represents the command to select a default workspace
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select a default workspace",
	Long: `Select a default workspace to use for all commands.

Examples:
  pipeops workspace select`,
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
			utils.PrintWarning("No workspaces found. Create a workspace first.", opts)
			return
		}

		if len(workspaces) == 1 {
			// Auto-select if only one workspace
			ws := workspaces[0]
			cfg.Settings.DefaultWorkspaceUUID = ws.UUID
			if err := config.Save(cfg); err != nil {
				utils.HandleError(err, "Error saving configuration", opts)
				return
			}
			utils.PrintSuccess(fmt.Sprintf("Selected workspace: %s (%s)", ws.Name, ws.UUID), opts)
			return
		}

		// Build options for selection
		options := make([]string, len(workspaces))
		for i, ws := range workspaces {
			marker := "  "
			if cfg.Settings != nil && cfg.Settings.DefaultWorkspaceUUID == ws.UUID {
				marker = "✓ "
			}
			options[i] = fmt.Sprintf("%s%s (%s)", marker, ws.Name, ws.UUID)
		}

		// Prompt user to select
		idx, _, err := utils.SelectOption("Select a workspace", options)
		if err != nil {
			utils.HandleError(err, "Selection cancelled", opts)
			return
		}

		selectedWS := workspaces[idx]
		cfg.Settings.DefaultWorkspaceUUID = selectedWS.UUID
		if err := config.Save(cfg); err != nil {
			utils.HandleError(err, "Error saving configuration", opts)
			return
		}

		utils.PrintSuccess(fmt.Sprintf("Selected workspace: %s (%s)", selectedWS.Name, selectedWS.UUID), opts)
		fmt.Println()
		fmt.Println("You can now run commands like:")
		fmt.Println("   pipeops project list")
		fmt.Println("   pipeops server list")
	},
	Args: cobra.NoArgs,
}

func (w *workspaceModel) selectWorkspace() {
	w.rootCmd.AddCommand(selectCmd)
}
