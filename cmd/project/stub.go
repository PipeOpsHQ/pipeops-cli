package project

import (
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update project configuration",
	Long: `Update the configuration of an existing project.

Examples:
  pipeops project update --id project-123 --name new-name`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'project update' command is coming soon! Please check our documentation for updates.", opts)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project",
	Long: `Delete a project from your PipeOps account.

Examples:
  pipeops project delete --id project-123`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'project delete' command is coming soon! Please check our documentation for updates.", opts)
	},
}

// GetUpdateCmd returns the update command for registration
func GetUpdateCmd() *cobra.Command {
	return updateCmd
}

// GetDeleteCmd returns the delete command for registration
func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}
