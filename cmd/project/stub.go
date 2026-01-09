package project

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
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

var deployCmd = &cobra.Command{
	Use:   "deploy <project-id>",
	Short: "Deploy a project",
	Long: `Deploy a project to trigger a new deployment.

Examples:
  pipeops project deploy proj-123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		projectID := args[0]

		client := pipeops.NewClient()
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		if !client.IsAuthenticated() {
			utils.HandleError(nil, "You are not logged in. Please run 'pipeops auth login' first.", opts)
			return
		}

		utils.PrintInfo(fmt.Sprintf("Deploying project %s...", projectID), opts)

		if err := client.DeployProject(projectID); err != nil {
			// Check if it's a 404 error (API not implemented)
			if strings.Contains(err.Error(), "404") {
				utils.PrintWarning("The deploy API is not yet available. Please use the PipeOps dashboard to deploy projects.", opts)
				return
			}
			utils.HandleError(err, "Error deploying project", opts)
			return
		}

		utils.PrintSuccess(fmt.Sprintf("Deployment triggered for project %s", projectID), opts)
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

// GetDeployCmd returns the deploy command for registration
func GetDeployCmd() *cobra.Command {
	return deployCmd
}
