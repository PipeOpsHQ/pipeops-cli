package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/proxy"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var proxyManager = proxy.NewManager()

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "🌐 Manage local proxy connections to deployed services",
	Long: `🌐 The proxy command allows you to create local port forwards to your deployed
services, making them accessible on your local machine. This is useful for debugging,
development, and accessing services that aren't publicly exposed.

Examples:
  - Start a proxy to a service in linked project:
    pipeops proxy start web-service --port 8080

  - Start a proxy to a specific project service:
    pipeops proxy start proj-123 web-service --port 8080

  - Start a proxy to an addon service:
    pipeops proxy start web-service --addon addon-456 --port 6379

  - List active proxies:
    pipeops proxy list

  - Stop a proxy:
    pipeops proxy stop proxy-123456

  - Stop all proxies:
    pipeops proxy stop-all`,
	Aliases: []string{"port-forward", "pf"},
}

var proxyStartCmd = &cobra.Command{
	Use:   "start [project-id] <service-name>",
	Short: "🚀 Start a proxy to a service",
	Long: `🚀 Start a local proxy connection to a deployed service. The service will be
accessible on your local machine through the specified port.

If no project ID is provided, uses the linked project from the current directory.

Examples:
  - Proxy to a service in linked project:
    pipeops proxy start web-service --port 8080

  - Proxy to a service in specific project:
    pipeops proxy start proj-123 web-service --port 8080

  - Auto-assign local port:
    pipeops proxy start api-service

  - Proxy to an addon service:
    pipeops proxy start redis --addon addon-456 --port 6379`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		var projectID, serviceName string
		var err error

		switch len(args) {
		case 1:
			// Only service name provided, try to get linked project
			serviceName = args[0]
			projectID, err = utils.GetProjectIDOrLinked("")
			if err != nil {
				utils.PrintError(err.Error(), opts)
				utils.PrintInfo("Use 'pipeops link <project-id>' to link a project to this directory", opts)
				utils.PrintInfo("Or provide both: pipeops proxy start <project-id> <service-name>", opts)
				return
			}
		case 2:
			// Both project ID and service name provided
			projectID = args[0]
			serviceName = args[1]
		default:
			utils.PrintError("Service name is required", opts)
			utils.PrintInfo("Usage: pipeops proxy start [project-id] <service-name>", opts)
			return
		}

		// Validate project ID
		if err := validation.ValidateProjectID(projectID); err != nil {
			utils.PrintError(fmt.Sprintf("Invalid project ID: %v", err), opts)
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

		// Show project context
		utils.PrintProjectContextWithOptions(projectID, opts)

		// Parse flags
		portStr, _ := cmd.Flags().GetString("port")
		targetPortStr, _ := cmd.Flags().GetString("target-port")
		addonID, _ := cmd.Flags().GetString("addon")
		daemon, _ := cmd.Flags().GetBool("daemon")

		// Parse local port
		localPort, err := proxy.GetPortFromString(portStr)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Invalid local port: %v", err), opts)
			return
		}

		// Parse target port
		targetPort, err := proxy.GetPortFromString(targetPortStr)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Invalid target port: %v", err), opts)
			return
		}

		// Check if local port is available
		if localPort > 0 && !proxy.IsPortAvailable(localPort) {
			utils.PrintError(fmt.Sprintf("Local port %d is already in use", localPort), opts)
			return
		}

		// Build proxy request
		req := &models.ProxyRequest{
			Target: models.ProxyTarget{
				ProjectID:   projectID,
				AddonID:     addonID,
				ServiceName: serviceName,
				Port:        targetPort,
			},
			LocalPort: localPort,
		}

		// Get proxy connection details from API
		message := fmt.Sprintf("Getting connection details for service '%s'", serviceName)
		if addonID != "" {
			message += fmt.Sprintf(" (addon: %s)", addonID)
		}
		utils.PrintInfo(message, opts)

		proxyResp, err := client.StartProxy(req)
		if err != nil {
			utils.HandleError(err, "Error starting proxy", opts)
			return
		}

		// Start local proxy
		localProxyResp, err := proxyManager.StartProxy(
			req.Target,
			proxyResp.LocalPort,
			proxyResp.RemoteHost,
			proxyResp.RemotePort,
		)
		if err != nil {
			utils.HandleError(err, "Error starting local proxy", opts)
			return
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(localProxyResp)
			return
		}

		utils.PrintSuccess("Proxy started successfully!", opts)
		fmt.Printf("🆔 Proxy ID: %s\n", localProxyResp.ProxyID)
		fmt.Printf("🌐 Local: http://localhost:%d\n", localProxyResp.LocalPort)
		fmt.Printf("🎯 Remote: %s:%d\n", localProxyResp.RemoteHost, localProxyResp.RemotePort)

		if daemon {
			utils.PrintInfo("Running in daemon mode. Use 'pipeops proxy stop' to stop.", opts)
			return
		}

		// Run in foreground mode
		utils.PrintInfo("Proxy is running... (Press Ctrl+C to stop)", opts)

		// Set up signal handling
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Wait for signal
		<-sigChan

		// Stop the proxy
		utils.PrintInfo("Stopping proxy...", opts)
		if err := proxyManager.StopProxy(localProxyResp.ProxyID); err != nil {
			utils.PrintError(fmt.Sprintf("Error stopping proxy: %v", err), opts)
		} else {
			utils.PrintSuccess("Proxy stopped successfully.", opts)
		}
	},
	Args: cobra.RangeArgs(1, 2),
}

var proxyListCmd = &cobra.Command{
	Use:   "list",
	Short: "📋 List active proxy connections",
	Long: `📋 List all currently active proxy connections, showing their status,
local and remote endpoints, and connection statistics.

Examples:
  - List all active proxies:
    pipeops proxy list

  - List proxies in JSON format:
    pipeops proxy list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		proxies := proxyManager.ListProxies()

		if len(proxies.Proxies) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON([]interface{}{})
			} else {
				utils.PrintWarning("No active proxy connections.", opts)
			}
			return
		}

		// Format output
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(proxies.Proxies)
		} else {
			// Prepare table data
			headers := []string{"PROXY ID", "STATUS", "LOCAL PORT", "REMOTE", "CONNECTIONS", "STARTED"}
			var rows [][]string

			for _, proxy := range proxies.Proxies {
				startTime, _ := time.Parse(time.RFC3339, proxy.StartedAt)
				timeAgo := time.Since(startTime).Round(time.Second)

				rows = append(rows, []string{
					proxy.ProxyID,
					utils.GetStatusIcon(proxy.Status) + " " + proxy.Status,
					fmt.Sprintf("%d", proxy.LocalPort),
					fmt.Sprintf("%s:%d", proxy.RemoteHost, proxy.RemotePort),
					fmt.Sprintf("%d", proxy.ConnectionsIn),
					timeAgo.String() + " ago",
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d active proxies", len(proxies.Proxies)), opts)
		}
	},
	Args: cobra.NoArgs,
}

var proxyStopCmd = &cobra.Command{
	Use:   "stop <proxy-id>",
	Short: "🛑 Stop a proxy connection",
	Long: `🛑 Stop a specific proxy connection by its ID. You can get proxy IDs
using the 'pipeops proxy list' command.

Examples:
  - Stop a specific proxy:
    pipeops proxy stop proxy-123456

  - Stop with JSON output:
    pipeops proxy stop proxy-123456 --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		if len(args) != 1 {
			utils.PrintError("Proxy ID is required", opts)
			utils.PrintInfo("Usage: pipeops proxy stop <proxy-id>", opts)
			return
		}

		proxyID := args[0]

		if err := proxyManager.StopProxy(proxyID); err != nil {
			utils.HandleError(err, "Error stopping proxy", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			result := map[string]interface{}{
				"proxy_id": proxyID,
				"status":   "stopped",
			}
			utils.PrintJSON(result)
		} else {
			utils.PrintSuccess(fmt.Sprintf("Proxy %s stopped successfully", proxyID), opts)
		}
	},
	Args: cobra.ExactArgs(1),
}

var proxyStopAllCmd = &cobra.Command{
	Use:   "stop-all",
	Short: "🛑 Stop all proxy connections",
	Long: `🛑 Stop all currently active proxy connections.

Examples:
  - Stop all proxies:
    pipeops proxy stop-all

  - Stop all with JSON output:
    pipeops proxy stop-all --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		proxies := proxyManager.ListProxies()

		if len(proxies.Proxies) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				result := map[string]interface{}{
					"stopped_count": 0,
					"message":       "No active proxy connections to stop",
				}
				utils.PrintJSON(result)
			} else {
				utils.PrintWarning("No active proxy connections to stop.", opts)
			}
			return
		}

		if err := proxyManager.StopAllProxies(); err != nil {
			utils.HandleError(err, "Error stopping proxies", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			result := map[string]interface{}{
				"stopped_count": len(proxies.Proxies),
				"status":        "success",
			}
			utils.PrintJSON(result)
		} else {
			utils.PrintSuccess(fmt.Sprintf("Stopped %d proxy connections", len(proxies.Proxies)), opts)
		}
	},
	Args: cobra.NoArgs,
}

var proxyServicesCmd = &cobra.Command{
	Use:   "services [project-id]",
	Short: "📋 List available services for a project",
	Long: `📋 List all services available for proxying in a specific project or addon.
This shows you what services you can connect to using the proxy command.

If no project ID is provided, uses the linked project from the current directory.

Examples:
  - List services for linked project:
    pipeops proxy services

  - List services for specific project:
    pipeops proxy services proj-123

  - List addon services:
    pipeops proxy services --addon addon-456

  - List services in JSON format:
    pipeops proxy services --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		var projectID string
		var err error

		if len(args) == 1 {
			projectID = args[0]
		} else {
			projectID, err = utils.GetProjectIDOrLinked("")
			if err != nil {
				utils.PrintError(err.Error(), opts)
				utils.PrintInfo("Use 'pipeops link <project-id>' to link a project to this directory", opts)
				utils.PrintInfo("Or provide: pipeops proxy services <project-id>", opts)
				return
			}
		}

		// Validate project ID
		if err := validation.ValidateProjectID(projectID); err != nil {
			utils.PrintError(fmt.Sprintf("Invalid project ID: %v", err), opts)
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

		// Show project context
		utils.PrintProjectContextWithOptions(projectID, opts)

		// Parse flags
		addonID, _ := cmd.Flags().GetString("addon")

		// Get services
		message := fmt.Sprintf("Fetching services for project %s", projectID)
		if addonID != "" {
			message += fmt.Sprintf(" (addon: %s)", addonID)
		}
		utils.PrintInfo(message, opts)

		services, err := client.GetServices(projectID, addonID)
		if err != nil {
			utils.HandleError(err, "Error fetching services", opts)
			return
		}

		if len(services.Services) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON([]interface{}{})
			} else {
				utils.PrintWarning("No services found.", opts)
			}
			return
		}

		// Format output
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(services.Services)
		} else {
			// Prepare table data
			headers := []string{"SERVICE NAME", "TYPE", "PORT", "PROTOCOL", "HEALTH", "DESCRIPTION"}
			var rows [][]string

			for _, service := range services.Services {
				description := utils.TruncateString(service.Description, 30)
				health := utils.GetStatusIcon(service.Health) + " " + service.Health

				rows = append(rows, []string{
					service.Name,
					service.Type,
					fmt.Sprintf("%d", service.Port),
					service.Protocol,
					health,
					description,
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d services", len(services.Services)), opts)

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\n💡 TIPS\n")
				fmt.Printf("├─ Start proxy: pipeops proxy start <service-name> --port <local-port>\n")
				fmt.Printf("└─ List proxies: pipeops proxy list\n")
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	// Add proxy command to root
	rootCmd.AddCommand(proxyCmd)

	// Add subcommands
	proxyCmd.AddCommand(proxyStartCmd)
	proxyCmd.AddCommand(proxyListCmd)
	proxyCmd.AddCommand(proxyStopCmd)
	proxyCmd.AddCommand(proxyStopAllCmd)
	proxyCmd.AddCommand(proxyServicesCmd)

	// Add flags to start command
	proxyStartCmd.Flags().StringP("port", "p", "", "Local port to bind to (auto-assign if not specified)")
	proxyStartCmd.Flags().String("target-port", "", "Target port on the service (default: service's default port)")
	proxyStartCmd.Flags().StringP("addon", "a", "", "Addon ID if connecting to an addon service")
	proxyStartCmd.Flags().BoolP("daemon", "d", false, "Run in daemon mode (background)")

	// Add flags to services command
	proxyServicesCmd.Flags().StringP("addon", "a", "", "List services for a specific addon")
}
