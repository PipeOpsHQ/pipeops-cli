package server

import (
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new server",
	Long: `Create a new server in your PipeOps account.

Examples:
  pipeops server create --name my-server --type agent --region us-east`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'server create' command is coming soon! Please use the PipeOps web console at https://app.pipeops.io or 'pipeops agent install' to install on an existing server.", opts)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update server configuration",
	Long: `Update the configuration of an existing server.

Examples:
  pipeops server update --id server-123 --name new-name`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'server update' command is coming soon! Please check our documentation for updates.", opts)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a server",
	Long: `Delete a server from your PipeOps account.

Examples:
  pipeops server delete --id server-123`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'server delete' command is coming soon! Please check our documentation for updates.", opts)
	},
}

// GetCreateCmd returns the create command for registration
func GetCreateCmd() *cobra.Command {
	return createCmd
}

// GetUpdateCmd returns the update command for registration
func GetUpdateCmd() *cobra.Command {
	return updateCmd
}

// GetDeleteCmd returns the delete command for registration
func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}
