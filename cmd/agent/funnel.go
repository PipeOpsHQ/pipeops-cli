package agent

import (
	"fmt"
	"log"
	"strconv"

	"github.com/PipeOpsHQ/pipeops-cli/internal/tailscale"
	"github.com/spf13/cobra"
)

var funnelCmd = &cobra.Command{
	Use:   "funnel",
	Short: "🌐 Manage Tailscale Funnel for port 80 exposure",
	Long: `🌐 The "funnel" command manages Tailscale Funnel to expose your Kubernetes cluster's ingress on port 80 to the public internet.

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
	Short: "🌐 Enable Tailscale Funnel for port exposure",
	Long: `🌐 Enable Tailscale Funnel to expose a port to the public internet.

Examples:
  # Enable Funnel for port 80 (default)
  pipeops agent funnel enable

  # Enable Funnel for specific port
  pipeops agent funnel enable 8080`,
	Run: func(cmd *cobra.Command, args []string) {
		port := 80 // Default port
		if len(args) > 0 {
			if p, err := strconv.Atoi(args[0]); err == nil && p > 0 && p <= 65535 {
				port = p
			}
		}

		tsClient := tailscale.NewClient()

		// Check if Tailscale is installed
		if !tsClient.IsInstalled() {
			log.Fatal("❌ Tailscale is not installed. Please run 'pipeops agent install' first.")
		}

		// Check if Tailscale is connected
		connected, err := tsClient.IsConnected()
		if err != nil {
			log.Fatalf("❌ Error checking Tailscale connection: %v", err)
		}
		if !connected {
			log.Fatal("❌ Tailscale is not connected. Please run 'pipeops agent install' first.")
		}

		// Enable Funnel
		log.Printf("🌐 Enabling Tailscale Funnel for port %d...", port)
		if err := tsClient.EnableFunnel(port); err != nil {
			log.Fatalf("❌ Error enabling Funnel: %v", err)
		}

		log.Println("✅ Tailscale Funnel enabled successfully!")

		// Get and display the Funnel URL
		if url, err := tsClient.GetFunnelURL(); err == nil {
			log.Printf("🌍 Your service is now accessible at: %s", url)
		} else {
			log.Println("ℹ️ Run 'pipeops agent funnel url' to get your Funnel URL")
		}
	},
}

var funnelDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "🚫 Disable Tailscale Funnel",
	Long: `🚫 Disable Tailscale Funnel to stop public internet access.

This will remove public access to your services while keeping Tailscale VPN functionality intact.`,
	Run: func(cmd *cobra.Command, args []string) {
		tsClient := tailscale.NewClient()

		// Check if Tailscale is installed
		if !tsClient.IsInstalled() {
			log.Fatal("❌ Tailscale is not installed.")
		}

		// Disable Funnel
		log.Println("🚫 Disabling Tailscale Funnel...")
		if err := tsClient.DisableFunnel(); err != nil {
			log.Fatalf("❌ Error disabling Funnel: %v", err)
		}

		log.Println("✅ Tailscale Funnel disabled successfully!")
		log.Println("ℹ️ Your services are no longer accessible from the public internet")
	},
}

var funnelStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "📊 Check Tailscale Funnel status",
	Long: `📊 Check the current status of Tailscale Funnel and display detailed information.

This command shows:
- Funnel configuration status
- Active services and ports
- Connection information`,
	Run: func(cmd *cobra.Command, args []string) {
		tsClient := tailscale.NewClient()

		// Check if Tailscale is installed
		if !tsClient.IsInstalled() {
			log.Fatal("❌ Tailscale is not installed.")
		}

		// Get Funnel status
		log.Println("📊 Checking Tailscale Funnel status...")
		status, err := tsClient.GetFunnelStatus()
		if err != nil {
			log.Fatalf("❌ Error getting Funnel status: %v", err)
		}

		fmt.Println("Tailscale Funnel Status:")
		fmt.Println("========================")
		fmt.Println(status)
	},
}

var funnelUrlCmd = &cobra.Command{
	Use:   "url",
	Short: "🔗 Get Tailscale Funnel URL",
	Long: `🔗 Get the public URL for your Tailscale Funnel service.

This command displays the public internet URL where your service is accessible.`,
	Run: func(cmd *cobra.Command, args []string) {
		tsClient := tailscale.NewClient()

		// Check if Tailscale is installed
		if !tsClient.IsInstalled() {
			log.Fatal("❌ Tailscale is not installed.")
		}

		// Get Funnel URL
		log.Println("🔗 Getting Tailscale Funnel URL...")
		url, err := tsClient.GetFunnelURL()
		if err != nil {
			log.Fatalf("❌ Error getting Funnel URL: %v", err)
		}

		fmt.Printf("🌍 Your service is accessible at: %s\n", url)
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
