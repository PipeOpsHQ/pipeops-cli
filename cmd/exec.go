package cmd

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/terminal"
	"github.com/PipeOpsHQ/pipeops-cli/models"
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

		// Parse arguments
		var projectID, containerName string
		var command []string

		// Find the -- separator
		dashIndex := -1
		for i, arg := range args {
			if arg == "--" {
				dashIndex = i
				break
			}
		}

		if dashIndex == -1 {
			utils.HandleError(fmt.Errorf("command separator missing"), "Command separator '--' is required", opts)
			return
		}

		// Parse based on argument structure
		if dashIndex == 2 {
			// Both project ID and container name provided
			projectID = args[0]
			containerName = args[1]
			command = args[dashIndex+1:]
		} else if dashIndex == 1 {
			// Only container name provided, use linked project
			projectContext, err := utils.LoadProjectContext()
			if err != nil {
				utils.HandleError(err, "Error loading project context", opts)
				return
			}

			projectID = projectContext.ProjectID
			if projectID == "" {
				utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use format: pipeops exec run <project-id> <container> -- <command>", opts)
				return
			}

			containerName = args[0]
			command = args[dashIndex+1:]
		} else {
			utils.HandleError(fmt.Errorf("invalid arguments"), "Usage: pipeops exec run [project-id] <container> -- <command>", opts)
			return
		}

		if len(command) == 0 {
			utils.HandleError(fmt.Errorf("command is required"), "Command is required after '--'", opts)
			return
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Get user input
		user, _ := cmd.Flags().GetString("user")

		// Build exec request
		req := &models.ExecRequest{
			ProjectID:   projectID,
			ServiceName: containerName,
			Container:   containerName,
			Command:     command,
			Interactive: false,
			User:        user,
		}

		// Execute command
		utils.PrintInfo(fmt.Sprintf("Executing command in container '%s'...", containerName), opts)

		if containerName != "" {
			fmt.Printf(" (container: %s)", containerName)
		}

		execResp, err := client.StartExec(req)
		if err != nil {
			utils.HandleError(err, "Error executing command", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(execResp)
		} else {
			utils.PrintSuccess(fmt.Sprintf("Command executed successfully (ID: %s)", execResp.ExecID), opts)
			utils.PrintInfo(fmt.Sprintf("WebSocket URL: %s", execResp.WebSocketURL), opts)
		}
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

		// Parse arguments
		var projectID, containerName string

		switch len(args) {
		case 1:
			// Only container name provided, use linked project
			projectContext, err := utils.LoadProjectContext()
			if err != nil {
				utils.HandleError(err, "Error loading project context", opts)
				return
			}

			projectID = projectContext.ProjectID
			if projectID == "" {
				utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use format: pipeops shell <project-id> <container>", opts)
				return
			}

			containerName = args[0]
		case 2:
			// Both project ID and container name provided
			projectID = args[0]
			containerName = args[1]
		default:
			utils.HandleError(fmt.Errorf("invalid arguments"), "Usage: pipeops shell [project-id] <container>", opts)
			return
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Get user input
		user, _ := cmd.Flags().GetString("user")

		// Build shell request
		req := &models.ShellRequest{
			ProjectID:   projectID,
			ServiceName: containerName,
			Container:   containerName,
			User:        user,
		}

		// Start shell session
		utils.PrintInfo(fmt.Sprintf("Starting shell in container '%s'...", containerName), opts)

		if containerName != "" {
			fmt.Printf(" (container: %s)", containerName)
		}

		shellResp, err := client.StartShell(req)
		if err != nil {
			utils.HandleError(err, "Error starting shell", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(shellResp)
		} else {
			utils.PrintSuccess(fmt.Sprintf("Shell session started (ID: %s)", shellResp.SessionID), opts)
			utils.PrintInfo(fmt.Sprintf("WebSocket URL: %s", shellResp.WebSocketURL), opts)
		}
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

		// Get project ID
		var projectID string
		if len(args) == 1 {
			projectID = args[0]
		} else {
			projectContext, err := utils.LoadProjectContext()
			if err != nil {
				utils.HandleError(err, "Error loading project context", opts)
				return
			}

			projectID = projectContext.ProjectID
			if projectID == "" {
				utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use format: pipeops exec containers <project-id>", opts)
				return
			}
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Get containers
		utils.PrintInfo(fmt.Sprintf("Fetching containers for project '%s'...", projectID), opts)

		containers, err := client.GetContainers(projectID, "")
		if err != nil {
			utils.HandleError(err, "Error fetching containers", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(containers)
		} else {
			if len(containers.Containers) == 0 {
				utils.PrintWarning("No containers found for this project", opts)
				return
			}

			headers := []string{"CONTAINER NAME", "SERVICE", "STATUS", "IMAGE", "CREATED"}
			var rows [][]string

			for _, container := range containers.Containers {
				rows = append(rows, []string{
					container.Name,
					container.ServiceName,
					container.Status,
					container.Image,
					container.CreatedAt,
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d containers", len(containers.Containers)), opts)
		}
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
