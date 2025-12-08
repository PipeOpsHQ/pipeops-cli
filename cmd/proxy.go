package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/internal/proxy"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var proxyManager = proxy.NewManager()

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Manage local proxy connections to deployed services",
	Long: `The proxy command allows you to create local port forwards to your deployed
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
	Short: "Start a proxy to a service",
	Long: `Start a proxy to a service in your project.

This command creates a local proxy connection to a service, allowing you to access it as if it were running locally.

Examples:
  - Start a proxy to a project service:
    pipeops proxy start web-service --project proj-123 --port 8080

  - Start a proxy to a specific service:
    pipeops proxy start database --project proj-123 --port 5432`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'proxy start' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
	Args: cobra.ExactArgs(1),
}

var proxyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active proxy connections",
	Long: `List all currently active proxy connections, showing their status,
local and remote endpoints, and connection statistics.

Examples:
  - List all active proxies:
    pipeops proxy list

  - List proxies in JSON format:
    pipeops proxy list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'proxy list' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
	Args: cobra.NoArgs,
}

var proxyStopCmd = &cobra.Command{
	Use:   "stop <proxy-id>",
	Short: "Stop a proxy connection",
	Long: `Stop a specific proxy connection by its ID. You can get proxy IDs
using the 'pipeops proxy list' command.

Examples:
  - Stop a specific proxy:
    pipeops proxy stop proxy-123456

  - Stop with JSON output:
    pipeops proxy stop proxy-123456 --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'proxy stop' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
	Args: cobra.ExactArgs(1),
}

var proxyStopAllCmd = &cobra.Command{
	Use:   "stop-all",
	Short: "Stop all proxy connections",
	Long: `Stop all currently active proxy connections.

Examples:
  - Stop all proxies:
    pipeops proxy stop-all

  - Stop all with JSON output:
    pipeops proxy stop-all --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'proxy stop-all' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
	Args: cobra.NoArgs,
}

var proxyServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List available services for proxying",
	Long: `List all services available for proxying in a specific project.

This command shows all the services you can proxy to in your project.

Examples:
  - List project services:
    pipeops proxy services --project proj-123

  - List services (with linked project):
    pipeops proxy services`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'proxy services' command is coming soon! Please check our documentation for updates.", opts)
		return
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
