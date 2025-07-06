package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/proxy"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/models"
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
  - Start a proxy to a project service:
    pipeops proxy start proj-123 web-service --port 8080

  - Start a proxy to an addon service:
    pipeops proxy start proj-123 redis --addon addon-456 --port 6379

  - List active proxies:
    pipeops proxy list

  - Stop a proxy:
    pipeops proxy stop proxy-123456

  - Stop all proxies:
    pipeops proxy stop-all`,
	Aliases: []string{"port-forward", "pf"},
}

var proxyStartCmd = &cobra.Command{
	Use:   "start <project-id> <service-name>",
	Short: "🚀 Start a proxy to a service",
	Long: `🚀 Start a local proxy connection to a deployed service. The service will be
accessible on your local machine through the specified port.

Examples:
  - Proxy to a web service on local port 8080:
    pipeops proxy start proj-123 web-service --port 8080

  - Auto-assign local port:
    pipeops proxy start proj-123 api-service

  - Proxy to an addon service:
    pipeops proxy start proj-123 redis --addon addon-456 --port 6379`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("❌ Project ID and service name are required")
			fmt.Println("Usage: pipeops proxy start <project-id> <service-name>")
			return
		}

		projectID := args[0]
		serviceName := args[1]

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
		portStr, _ := cmd.Flags().GetString("port")
		targetPortStr, _ := cmd.Flags().GetString("target-port")
		addonID, _ := cmd.Flags().GetString("addon")
		daemon, _ := cmd.Flags().GetBool("daemon")

		// Parse local port
		localPort, err := proxy.GetPortFromString(portStr)
		if err != nil {
			fmt.Printf("❌ Invalid local port: %v\n", err)
			return
		}

		// Parse target port
		targetPort, err := proxy.GetPortFromString(targetPortStr)
		if err != nil {
			fmt.Printf("❌ Invalid target port: %v\n", err)
			return
		}

		// Check if local port is available
		if localPort > 0 && !proxy.IsPortAvailable(localPort) {
			fmt.Printf("❌ Local port %d is already in use\n", localPort)
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
		fmt.Printf("🔍 Getting connection details for service '%s'", serviceName)
		if addonID != "" {
			fmt.Printf(" (addon: %s)", addonID)
		}
		fmt.Println("...")

		proxyResp, err := client.StartProxy(req)
		if err != nil {
			fmt.Printf("❌ Error starting proxy: %v\n", err)
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
			fmt.Printf("❌ Error starting local proxy: %v\n", err)
			return
		}

		fmt.Printf("✅ Proxy started successfully!\n")
		fmt.Printf("🆔 Proxy ID: %s\n", localProxyResp.ProxyID)
		fmt.Printf("🌐 Local: http://localhost:%d\n", localProxyResp.LocalPort)
		fmt.Printf("🎯 Remote: %s:%d\n", localProxyResp.RemoteHost, localProxyResp.RemotePort)

		if daemon {
			fmt.Println("🔄 Running in daemon mode. Use 'pipeops proxy stop' to stop.")
			return
		}

		// Run in foreground mode
		fmt.Println("🔄 Proxy is running... (Press Ctrl+C to stop)")

		// Set up signal handling
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Wait for signal
		<-sigChan

		// Stop the proxy
		fmt.Println("\n🛑 Stopping proxy...")
		if err := proxyManager.StopProxy(localProxyResp.ProxyID); err != nil {
			fmt.Printf("❌ Error stopping proxy: %v\n", err)
		} else {
			fmt.Println("✅ Proxy stopped successfully.")
		}
	},
	Args: cobra.ExactArgs(2),
}

var proxyListCmd = &cobra.Command{
	Use:   "list",
	Short: "📋 List active proxy connections",
	Long: `📋 List all currently active proxy connections, showing their status,
local and remote endpoints, and connection statistics.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxies := proxyManager.ListProxies()

		if len(proxies.Proxies) == 0 {
			fmt.Println("📭 No active proxy connections.")
			return
		}

		// Display header
		fmt.Printf("%-15s | %-10s | %-15s | %-20s | %-10s | %-15s\n",
			"PROXY ID", "STATUS", "LOCAL PORT", "REMOTE", "CONNECTIONS", "STARTED")
		fmt.Println(strings.Repeat("-", 100))

		// Display proxy details
		for _, proxy := range proxies.Proxies {
			startTime, _ := time.Parse(time.RFC3339, proxy.StartedAt)
			timeAgo := time.Since(startTime).Round(time.Second)

			fmt.Printf("%-15s | %-10s | %-15d | %-20s | %-10d | %-15s\n",
				proxy.ProxyID,
				proxy.Status,
				proxy.LocalPort,
				fmt.Sprintf("%s:%d", proxy.RemoteHost, proxy.RemotePort),
				proxy.ConnectionsIn,
				timeAgo.String()+" ago")
		}

		fmt.Printf("\n✅ Found %d active proxies.\n", len(proxies.Proxies))
	},
	Args: cobra.NoArgs,
}

var proxyStopCmd = &cobra.Command{
	Use:   "stop <proxy-id>",
	Short: "🛑 Stop a proxy connection",
	Long: `🛑 Stop a specific proxy connection by its ID. You can get proxy IDs
using the 'pipeops proxy list' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ Proxy ID is required")
			fmt.Println("Usage: pipeops proxy stop <proxy-id>")
			return
		}

		proxyID := args[0]

		if err := proxyManager.StopProxy(proxyID); err != nil {
			fmt.Printf("❌ Error stopping proxy: %v\n", err)
			return
		}

		fmt.Printf("✅ Proxy %s stopped successfully.\n", proxyID)
	},
	Args: cobra.ExactArgs(1),
}

var proxyStopAllCmd = &cobra.Command{
	Use:   "stop-all",
	Short: "🛑 Stop all proxy connections",
	Long:  `🛑 Stop all currently active proxy connections.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxies := proxyManager.ListProxies()

		if len(proxies.Proxies) == 0 {
			fmt.Println("📭 No active proxy connections to stop.")
			return
		}

		if err := proxyManager.StopAllProxies(); err != nil {
			fmt.Printf("❌ Error stopping proxies: %v\n", err)
			return
		}

		fmt.Printf("✅ Stopped %d proxy connections.\n", len(proxies.Proxies))
	},
	Args: cobra.NoArgs,
}

var proxyServicesCmd = &cobra.Command{
	Use:   "services <project-id>",
	Short: "📋 List available services for a project",
	Long: `📋 List all services available for proxying in a specific project or addon.
This shows you what services you can connect to using the proxy command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ Project ID is required")
			fmt.Println("Usage: pipeops proxy services <project-id>")
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

		// Get services
		fmt.Printf("🔍 Fetching services for project %s", projectID)
		if addonID != "" {
			fmt.Printf(" (addon: %s)", addonID)
		}
		fmt.Println("...")

		services, err := client.GetServices(projectID, addonID)
		if err != nil {
			fmt.Printf("❌ Error fetching services: %v\n", err)
			return
		}

		if len(services.Services) == 0 {
			fmt.Println("📭 No services found.")
			return
		}

		// Display header
		fmt.Printf("%-20s | %-10s | %-8s | %-10s | %-10s | %-30s\n",
			"SERVICE NAME", "TYPE", "PORT", "PROTOCOL", "HEALTH", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 110))

		// Display service details
		for _, service := range services.Services {
			description := service.Description
			if len(description) > 30 {
				description = description[:27] + "..."
			}

			fmt.Printf("%-20s | %-10s | %-8d | %-10s | %-10s | %-30s\n",
				service.Name,
				service.Type,
				service.Port,
				service.Protocol,
				service.Health,
				description)
		}

		fmt.Printf("\n✅ Found %d services.\n", len(services.Services))
	},
	Args: cobra.ExactArgs(1),
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
