package agent

import (
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var funnelCmd = &cobra.Command{
	Use:   "funnel",
	Short: "Manage Tailscale Funnel for port 80 exposure",
	Long: `The "funnel" command manages Tailscale Funnel to expose your Kubernetes cluster's ingress on port 80 to the public internet.

This command supports:
- Enable/disable Tailscale Funnel
- Check Funnel status and URLs
- Configure Funnel for specific services
- Manage Funnel permissions

Examples:
  # Enable Funnel for port 80
  pipeops agent funnel enable

  # Disable Funnel
  pipeops agent funnel disable

  # Check Funnel status
  pipeops agent funnel status

  # Get Funnel URL
  pipeops agent funnel url`,
}

var funnelEnableCmd = &cobra.Command{
	Use:   "enable [port]",
	Short: "Enable Tailscale Funnel for port exposure",
	Long: `Enable Tailscale Funnel to expose a port to the public internet.

Examples:
  # Enable Funnel for port 80 (default)
  pipeops agent funnel enable

  # Enable Funnel for specific port
  pipeops agent funnel enable 8080`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'agent funnel enable' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
}

var funnelDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable Tailscale Funnel",
	Long: `Disable Tailscale Funnel to stop public internet access.

This will remove public access to your services while keeping Tailscale VPN functionality intact.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'agent funnel disable' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
}

var funnelStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Tailscale Funnel status",
	Long: `Check the current status of Tailscale Funnel and display detailed information.

This command shows:
- Funnel configuration status
- Active services and ports
- Connection information`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'agent funnel status' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
}

var funnelUrlCmd = &cobra.Command{
	Use:   "url",
	Short: "Get Tailscale Funnel URL",
	Long: `Get the public URL for your Tailscale Funnel service.

This command displays the public internet URL where your service is accessible.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'agent funnel url' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
}

func (a *agentModel) funnel() {
	// Add subcommands
	funnelCmd.AddCommand(funnelEnableCmd)
	funnelCmd.AddCommand(funnelDisableCmd)
	funnelCmd.AddCommand(funnelStatusCmd)
	funnelCmd.AddCommand(funnelUrlCmd)

	a.rootCmd.AddCommand(funnelCmd)
}
