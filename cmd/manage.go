package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/manage"
	"github.com/spf13/cobra"
)

var manageCmd = &cobra.Command{
	Use:     "manage",
	Short:   "🔧 Manage operations effortlessly.",
	Long: `🔧 The manage command provides a set of subcommands for handling various operations, 
such as deploying, configuring, and interacting with projects on PipeOps. 🌐

Examples:
  - Configure project settings
  - Manage deployment pipelines
  - Monitor project status and logs`,
	Aliases: []string{"projects", "p"},
}

func init() {
	// Add the manage command as a subcommand of the root command
	rootCmd.AddCommand(manageCmd)

	// Register subcommands under the manage command
	registerManageSubcommands()
}

// registerManageSubcommands initializes and registers subcommands for the manage command
func registerManageSubcommands() {
	// Initialize and register management-related commands under the manage command
	manageSub := manage.NewManagement(manageCmd)
	manageSub.Register()
}
