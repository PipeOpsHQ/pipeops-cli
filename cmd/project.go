package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/project"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project-related operations.",
	Long: `The server command provides a set of subcommands for managing
project-related operations, such as deploying, configurations and interactions
with projects on PipeOps`,
}

func init() {
	// Add the server command as a subcommand of the root command
	rootCmd.AddCommand(projectCmd)

	// Register subcommands under the server command
	registerProjectSubcommands()
}

// registerProjectSubcommands initializes and registers subcommands for the server command
func registerProjectSubcommands() {
	// Initialize K3s-related commands under the server command
	k3sCmd := project.NewProject(projectCmd)
	k3sCmd.Register()
}
