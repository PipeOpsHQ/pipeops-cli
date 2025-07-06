package cmd

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/terminal"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var terminalManager = terminal.NewManager()

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "💻 Execute commands in deployed containers",
	Long: `💻 The exec command allows you to execute commands in your deployed containers.
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
	Use:   "run <project-id> <service-name> -- <command>",
	Short: "🚀 Execute a command in a container",
	Long: `🚀 Execute a specific command in a deployed container. Use -- to separate
the command from the flags.

Examples:
  - List files in a project container:
    pipeops exec run proj-123 web-service -- ls -la

  - Check application logs:
    pipeops exec run proj-123 web-service -- cat /app/logs/app.log

  - Run a script in an addon container:
    pipeops exec run proj-123 redis --addon addon-456 -- redis-cli info`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("❌ Project ID, service name, and command are required")
			fmt.Println("Usage: pipeops exec run <project-id> <service-name> -- <command>")
			return
		}

		projectID := args[0]
		serviceName := args[1]

		// Find the separator
		separatorIndex := -1
		for i, arg := range args {
			if arg == "--" {
				separatorIndex = i
				break
			}
		}

		if separatorIndex == -1 || separatorIndex == len(args)-1 {
			fmt.Println("❌ Command is required after --")
			fmt.Println("Usage: pipeops exec run <project-id> <service-name> -- <command>")
			return
		}

		command := args[separatorIndex+1:]

		// Validate project ID
		if err := validation.ValidateProjectID(projectID); err != nil {
			fmt.Printf("❌ Invalid project ID: %v\n", err)
			return
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if user is authenticated
		if !client.IsAuthenticated() {
			fmt.Println("❌ You are not logged in. Please run 'pipeops auth login' first.")
			return
		}

		// Parse flags
		addonID, _ := cmd.Flags().GetString("addon")
		container, _ := cmd.Flags().GetString("container")
		user, _ := cmd.Flags().GetString("user")
		workdir, _ := cmd.Flags().GetString("workdir")
		interactive, _ := cmd.Flags().GetBool("interactive")

		// Parse environment variables
		envFlags, _ := cmd.Flags().GetStringArray("env")
		environment := make(map[string]string)
		for _, env := range envFlags {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				environment[parts[0]] = parts[1]
			}
		}

		// Build exec request
		req := &models.ExecRequest{
			ProjectID:   projectID,
			AddonID:     addonID,
			ServiceName: serviceName,
			Container:   container,
			Command:     command,
			Interactive: interactive,
			Environment: environment,
			WorkingDir:  workdir,
			User:        user,
		}

		// Start exec session
		fmt.Printf("🚀 Starting exec session for service '%s'", serviceName)
		if addonID != "" {
			fmt.Printf(" (addon: %s)", addonID)
		}
		fmt.Printf(" with command: %s\n", strings.Join(command, " "))

		execResp, err := client.StartExec(req)
		if err != nil {
			fmt.Printf("❌ Error starting exec session: %v\n", err)
			return
		}

		fmt.Printf("💻 Exec session started (ID: %s)\n", execResp.ExecID)

		// Connect to terminal session
		if interactive {
			session, err := terminalManager.StartExecSession(execResp.ExecID, execResp.WebSocketURL, true)
			if err != nil {
				fmt.Printf("❌ Error connecting to terminal: %v\n", err)
				return
			}
			defer session.Close()

			fmt.Println("🔗 Connected to interactive terminal. Press Ctrl+C to exit.")
			session.WaitForCompletion()
		} else {
			// Non-interactive execution
			if err := terminalManager.ExecCommand(execResp.ExecID, execResp.WebSocketURL, command); err != nil {
				fmt.Printf("❌ Error executing command: %v\n", err)
				return
			}
		}
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell [project-id] <service-name>",
	Short: "🐚 Start an interactive shell in a container",
	Long: `🐚 Start an interactive shell session in a deployed container. This gives you
full access to the container's environment for debugging and exploration.

If no project ID is provided, uses the linked project from the current directory.

Examples:
  - Start a shell in linked project:
    pipeops shell web-service

  - Start a shell in specific project:
    pipeops shell proj-123 web-service

  - Start a shell in an addon container:
    pipeops shell web-service --addon addon-456

  - Start a specific shell:
    pipeops shell web-service --shell zsh`,
	Run: func(cmd *cobra.Command, args []string) {
		var projectID, serviceName string
		var err error

		if len(args) == 1 {
			// Only service name provided, try to get linked project
			serviceName = args[0]
			projectID, err = utils.GetLinkedProject()
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				fmt.Println("💡 Use 'pipeops link <project-id>' to link a project to this directory")
				fmt.Println("   Or provide both: pipeops shell <project-id> <service-name>")
				return
			}
		} else if len(args) == 2 {
			// Both project ID and service name provided
			projectID = args[0]
			serviceName = args[1]
		} else {
			fmt.Println("❌ Service name is required")
			fmt.Println("Usage: pipeops shell [project-id] <service-name>")
			return
		}

		// Validate project ID
		if err := validation.ValidateProjectID(projectID); err != nil {
			fmt.Printf("❌ Invalid project ID: %v\n", err)
			return
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if user is authenticated
		if !client.IsAuthenticated() {
			fmt.Println("❌ You are not logged in. Please run 'pipeops auth login' first.")
			return
		}

		// Show project context
		utils.PrintProjectContext(projectID)

		// Parse flags
		addonID, _ := cmd.Flags().GetString("addon")
		container, _ := cmd.Flags().GetString("container")
		user, _ := cmd.Flags().GetString("user")
		workdir, _ := cmd.Flags().GetString("workdir")
		shell, _ := cmd.Flags().GetString("shell")

		// Parse environment variables
		envFlags, _ := cmd.Flags().GetStringArray("env")
		environment := make(map[string]string)
		for _, env := range envFlags {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				environment[parts[0]] = parts[1]
			}
		}

		// Get terminal size
		cols, rows, err := terminal.GetTerminalSize()
		if err != nil {
			// Default size if we can't get it
			cols, rows = 80, 24
		}

		// Build shell request
		req := &models.ShellRequest{
			ProjectID:   projectID,
			AddonID:     addonID,
			ServiceName: serviceName,
			Container:   container,
			Shell:       shell,
			Environment: environment,
			WorkingDir:  workdir,
			User:        user,
			Cols:        cols,
			Rows:        rows,
		}

		// Start shell session
		fmt.Printf("🐚 Starting shell session for service '%s'", serviceName)
		if addonID != "" {
			fmt.Printf(" (addon: %s)", addonID)
		}
		fmt.Println("...")

		shellResp, err := client.StartShell(req)
		if err != nil {
			fmt.Printf("❌ Error starting shell session: %v\n", err)
			return
		}

		fmt.Printf("💻 Shell session started (ID: %s)\n", shellResp.SessionID)

		// Connect to terminal session
		session, err := terminalManager.StartShellSession(shellResp.SessionID, shellResp.WebSocketURL)
		if err != nil {
			fmt.Printf("❌ Error connecting to terminal: %v\n", err)
			return
		}
		defer session.Close()

		fmt.Println("🔗 Connected to interactive shell. Press Ctrl+C to exit.")
		session.WaitForCompletion()
	},
	Args: cobra.RangeArgs(1, 2),
}

var execContainersCmd = &cobra.Command{
	Use:   "containers <project-id>",
	Short: "📦 List available containers",
	Long: `📦 List all containers available for exec access in a specific project or addon.
This shows you what containers you can connect to.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ Project ID is required")
			fmt.Println("Usage: pipeops exec containers <project-id>")
			return
		}

		projectID := args[0]

		// Validate project ID
		if err := validation.ValidateProjectID(projectID); err != nil {
			fmt.Printf("❌ Invalid project ID: %v\n", err)
			return
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if user is authenticated
		if !client.IsAuthenticated() {
			fmt.Println("❌ You are not logged in. Please run 'pipeops auth login' first.")
			return
		}

		// Parse flags
		addonID, _ := cmd.Flags().GetString("addon")

		// Get containers
		fmt.Printf("🔍 Fetching containers for project %s", projectID)
		if addonID != "" {
			fmt.Printf(" (addon: %s)", addonID)
		}
		fmt.Println("...")

		containers, err := client.GetContainers(projectID, addonID)
		if err != nil {
			fmt.Printf("❌ Error fetching containers: %v\n", err)
			return
		}

		if len(containers.Containers) == 0 {
			fmt.Println("📭 No containers found.")
			return
		}

		// Display header
		fmt.Printf("%-20s | %-15s | %-20s | %-12s | %-8s | %-15s\n",
			"CONTAINER NAME", "SERVICE", "IMAGE", "STATUS", "RESTARTS", "STARTED")
		fmt.Println(strings.Repeat("-", 100))

		// Display container details
		for _, container := range containers.Containers {
			image := container.Image
			if len(image) > 20 {
				image = image[:17] + "..."
			}

			startedAt := "N/A"
			if container.StartedAt != "" {
				startedAt = container.StartedAt[:10] // Just the date part
			}

			fmt.Printf("%-20s | %-15s | %-20s | %-12s | %-8d | %-15s\n",
				container.Name,
				container.ServiceName,
				image,
				container.Status,
				container.RestartCount,
				startedAt)
		}

		fmt.Printf("\n✅ Found %d containers.\n", len(containers.Containers))
	},
	Args: cobra.ExactArgs(1),
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
	execRunCmd.Flags().StringP("addon", "a", "", "Addon ID if connecting to an addon service")
	execRunCmd.Flags().StringP("container", "c", "", "Specific container name")
	execRunCmd.Flags().StringP("user", "u", "", "User to run command as")
	execRunCmd.Flags().StringP("workdir", "w", "", "Working directory")
	execRunCmd.Flags().BoolP("interactive", "i", false, "Run in interactive mode")
	execRunCmd.Flags().StringArrayP("env", "e", nil, "Environment variables (KEY=VALUE)")

	// Add flags to shell command
	shellCmd.Flags().StringP("addon", "a", "", "Addon ID if connecting to an addon service")
	shellCmd.Flags().StringP("container", "c", "", "Specific container name")
	shellCmd.Flags().StringP("user", "u", "", "User to run shell as")
	shellCmd.Flags().StringP("workdir", "w", "", "Working directory")
	shellCmd.Flags().StringP("shell", "s", "", "Shell to use (bash, sh, zsh, etc.)")
	shellCmd.Flags().StringArrayP("env", "e", nil, "Environment variables (KEY=VALUE)")

	// Add flags to containers command
	execContainersCmd.Flags().StringP("addon", "a", "", "List containers for a specific addon")
}
