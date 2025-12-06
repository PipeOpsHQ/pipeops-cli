package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/project"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project-related operations.",
	Long: `The project command provides a set of subcommands for managing 
project-related operations on PipeOps. These include deploying, configuring, 
and interacting with projects seamlessly.

Examples:
  - List all projects:
    pipeops project list

  - Deploy a project:
    pipeops project deploy --name my-project

  - Update project configurations:
    pipeops project update --id project-id`,
	Aliases: []string{"projects", "p"},
}

func init() {
	// Add the project command as a subcommand of the root command
	rootCmd.AddCommand(projectCmd)

	// Register subcommands under the project command
	registerProjectSubcommands()
}

// registerProjectSubcommands initializes and registers subcommands for the project command
func registerProjectSubcommands() {
	// Initialize and register project-related commands under the project command
	projectSub := project.NewProject(projectCmd)
	projectSub.Register()
}
