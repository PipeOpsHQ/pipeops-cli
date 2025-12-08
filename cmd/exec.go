package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/internal/terminal"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var terminalManager = terminal.NewManager()

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute commands in deployed containers",
	Long: `The exec command allows you to execute commands in your deployed containers.
This is useful for debugging, running maintenance tasks, or exploring your application environment.

Examples:
  - Execute a command in a project container:
    pipeops exec proj-123 web-service -- ls -la

  - Execute a command in an addon container:
    pipeops exec proj-123 redis --addon addon-456 -- redis-cli ping

  - Start an interactive shell:
    pipeops shell proj-123 web-service

  - List available containers:
    pipeops exec containers proj-123`,
	Aliases: []string{"execute", "run"},
}

var execRunCmd = &cobra.Command{
	Use:   "run [project-id] <container-name> -- <command>",
	Short: "Execute a command in a container",
	Long: `Execute a command in a container within your project.

This command allows you to run arbitrary commands inside containers, useful for debugging, maintenance, or data operations.

Examples:
  - Execute a command in a container:
    pipeops exec run proj-123 web-container -- ls -la

  - Run a script in a container:
    pipeops exec run proj-123 web-container -- node script.js`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'exec run' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
	Args: cobra.MinimumNArgs(3),
}

var shellCmd = &cobra.Command{
	Use:   "shell [project-id] <container-name>",
	Short: "Start an interactive shell in a container",
	Long: `Start an interactive shell session in a container within your project.

This provides direct shell access to containers for debugging, maintenance, or interactive operations.

Examples:
  - Start a shell in a container:
    pipeops shell proj-123 web-container

  - Start a shell (with linked project):
    pipeops shell web-container`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'shell' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
	Args: cobra.RangeArgs(1, 2),
}

var execContainersCmd = &cobra.Command{

	Use:   "containers [project-id]",

	Short: "List containers available for exec",

	Long: `List all containers available for exec access in a specific project.



This command shows all containers you can execute commands in or start shells within.



Examples:

  - List containers for linked project:

    pipeops exec containers



  - List containers for specific project:

    pipeops exec containers proj-123`,

	Run: func(cmd *cobra.Command, args []string) {

		opts := utils.GetOutputOptions(cmd)

		utils.PrintWarning("The 'exec containers' command is coming soon! Please check our documentation for updates.", opts)

		return

	},

	Args: cobra.MaximumNArgs(1),

}

func init() {
	// Add exec command to root
	rootCmd.AddCommand(execCmd)

	// Add shell command to root
	rootCmd.AddCommand(shellCmd)

	// Add subcommands to exec
	execCmd.AddCommand(execRunCmd)
	execCmd.AddCommand(execContainersCmd)

	// Add flags to exec run command
	execRunCmd.Flags().StringP("user", "u", "", "User to run command as")

	// Add flags to shell command
	shellCmd.Flags().StringP("user", "u", "", "User to run shell as")

	// Add flags to containers command
	execContainersCmd.Flags().StringP("project", "p", "", "Project ID")
}
