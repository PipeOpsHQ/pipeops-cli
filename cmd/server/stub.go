package server

import (
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new server",
	Long: `Create a new server in your PipeOps account.

Note: Server creation via CLI is coming soon. For now, please use:
  - PipeOps web console at https://app.pipeops.io
  - 'pipeops agent install' to add an existing server

Examples:
  pipeops server create`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'server create' command is coming soon! Please use the PipeOps web console at https://app.pipeops.io or 'pipeops agent install' to install on an existing server.", opts)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [server-id]",
	Short: "Delete a server",
	Long: `Delete a server from your PipeOps account.

Examples:
  pipeops server delete server-123`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'server delete' command is coming soon! Please check our documentation for updates.", opts)
	},
}

// GetCreateCmd returns the create command for registration
func GetCreateCmd() *cobra.Command {
	return createCmd
}

// GetDeleteCmd returns the delete command for registration
func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}
