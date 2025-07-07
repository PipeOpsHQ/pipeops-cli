package cmd

import (
	"fmt"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/proxy"
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
	Use:   "start <service-name>",
	Short: "🔗 Start a proxy to a service",
	Long: `🔗 Start a proxy to a service in your project.

This command creates a local proxy connection to a service, allowing you to access it as if it were running locally.

Examples:
  - Start a proxy to a project service:
    pipeops proxy start web-service --project proj-123 --port 8080

  - Start a proxy to a specific service:
    pipeops proxy start database --project proj-123 --port 5432`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		if len(args) == 0 {
			utils.HandleError(fmt.Errorf("service name is required"), "Service name is required", opts)
			return
		}

		serviceName := args[0]

		// Get project context
		projectContext, err := utils.LoadProjectContext()
		if err != nil {
			utils.HandleError(err, "Error loading project context", opts)
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

		// Get flags
		projectID := projectContext.ProjectID
		if projectID == "" {
			flagProjectID, _ := cmd.Flags().GetString("project")
			if flagProjectID == "" {
				utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use --project flag or link a project with 'pipeops link'", opts)
				return
			}
			projectID = flagProjectID
		}

		localPort, _ := cmd.Flags().GetInt("port")
		if localPort == 0 {
			localPort = 8080 // default port
		}

		// Create proxy request
		req := &models.ProxyRequest{
			Target: models.ProxyTarget{
				ProjectID:   projectID,
				ServiceName: serviceName,
				Port:        0, // Use service's default port
			},
			LocalPort: localPort,
		}

		// Start proxy
		utils.PrintInfo(fmt.Sprintf("Starting proxy to service '%s' on port %d...", serviceName, localPort), opts)

		message := fmt.Sprintf("Starting proxy to service '%s' in project '%s'", serviceName, projectID)

		utils.PrintInfo(message, opts)

		proxyResp, err := client.StartProxy(req)
		if err != nil {
			utils.HandleError(err, "Error starting proxy", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(proxyResp)
		} else {
			utils.PrintSuccess(fmt.Sprintf("Proxy started successfully on port %d", localPort), opts)
			utils.PrintInfo(fmt.Sprintf("Remote endpoint: %s:%d", proxyResp.RemoteHost, proxyResp.RemotePort), opts)
			utils.PrintInfo("Press Ctrl+C to stop the proxy", opts)
		}
	},
	Args: cobra.ExactArgs(1),
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
	Use:   "services",
	Short: "📋 List available services for proxying",
	Long: `📋 List all services available for proxying in a specific project.

This command shows all the services you can proxy to in your project.

Examples:
  - List project services:
    pipeops proxy services --project proj-123

  - List services (with linked project):
    pipeops proxy services`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Get project context
		projectContext, err := utils.LoadProjectContext()
		if err != nil {
			utils.HandleError(err, "Error loading project context", opts)
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

		// Get project ID
		projectID := projectContext.ProjectID
		if projectID == "" {
			flagProjectID, _ := cmd.Flags().GetString("project")
			if flagProjectID == "" {
				utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use --project flag or link a project with 'pipeops link'", opts)
				return
			}
			projectID = flagProjectID
		}

		message := fmt.Sprintf("Fetching services for project '%s'...", projectID)

		utils.PrintInfo(message, opts)

		services, err := client.GetServices(projectID, "")
		if err != nil {
			utils.HandleError(err, "Error fetching services", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(services)
		} else {
			if len(services.Services) == 0 {
				utils.PrintWarning("No services found for this project", opts)
				return
			}

			headers := []string{"SERVICE NAME", "TYPE", "PORT", "PROTOCOL", "HEALTH"}
			var rows [][]string

			for _, service := range services.Services {
				rows = append(rows, []string{
					service.Name,
					service.Type,
					fmt.Sprintf("%d", service.Port),
					service.Protocol,
					service.Health,
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d services", len(services.Services)), opts)
		}
	},
	Args: cobra.NoArgs,
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
	proxyStartCmd.Flags().StringP("project", "p", "", "Project ID")
	proxyStartCmd.Flags().IntP("port", "", 8080, "Local port for the proxy")

	// Add flags to services command
	proxyServicesCmd.Flags().StringP("project", "p", "", "Project ID")
}
