package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/PipeOpsHQ/pipeops-cli/internal/tailscale"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [auth-key]",
	Short: "🚀 Install Tailscale and configure Funnel for port 80 exposure",
	Long: `🚀 The "install" command installs Tailscale and configures it to expose your Kubernetes cluster's ingress on port 80 using Tailscale Funnel.

This command supports:
- Automatic Tailscale installation
- Kubernetes cluster setup with Tailscale integration
- Tailscale Funnel configuration for public port 80 access
- Tailscale Kubernetes operator installation
- Ingress configuration with Funnel annotations

Examples:
  # Install with auth key as argument
  pipeops agent install tskey-auth-your-key-here

  # Install using environment variables
  export TAILSCALE_AUTH_KEY="tskey-auth-your-key-here"
  export CLUSTER_NAME="my-cluster"
  pipeops agent install

  # Install on existing cluster
  pipeops agent install --existing-cluster --cluster-name="my-existing-cluster"

  # Install without Funnel (private access only)
  pipeops agent install --no-funnel`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get auth key from args or environment
		authKey := getAuthKey(cmd, args)

		// Get cluster name from flag or environment
		clusterName, _ := cmd.Flags().GetString("cluster-name")
		if clusterName == "" {
			clusterName = os.Getenv("CLUSTER_NAME")
		}
		if clusterName == "" {
			clusterName = "tailscale-cluster"
		}

		// Get cluster type from flag or environment
		clusterType, _ := cmd.Flags().GetString("cluster-type")
		if clusterType == "" {
			clusterType = os.Getenv("CLUSTER_TYPE")
		}
		if clusterType == "" {
			clusterType = "k3s" // Default to k3s for Tailscale integration
		}

		// Check if installing on existing cluster
		existingCluster, _ := cmd.Flags().GetBool("existing-cluster")

		// Check if Funnel should be disabled
		noFunnel, _ := cmd.Flags().GetBool("no-funnel")

		// Check if this is an update operation
		update, _ := cmd.Flags().GetBool("update")

		// Check if this is an uninstall operation
		uninstall, _ := cmd.Flags().GetBool("uninstall")

		if uninstall {
			uninstallTailscale(cmd)
			return
		}

		if update {
			updateTailscale(cmd, authKey, clusterName)
			return
		}

		if existingCluster {
			installOnExistingCluster(cmd, authKey, clusterName, !noFunnel)
		} else {
			installNewCluster(cmd, authKey, clusterName, clusterType, !noFunnel)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// Auth key can be provided as argument or environment variable
		if len(args) == 0 {
			authKey := os.Getenv("TAILSCALE_AUTH_KEY")
			if authKey == "" {
				return fmt.Errorf("❌ Tailscale auth key is required either as argument or TAILSCALE_AUTH_KEY environment variable")
			}
		}
		return nil
	},
}

// getAuthKey retrieves auth key from args or environment
func getAuthKey(cmd *cobra.Command, args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return os.Getenv("TAILSCALE_AUTH_KEY")
}

// installNewCluster installs a new Kubernetes cluster with Tailscale and Funnel
func installNewCluster(cmd *cobra.Command, authKey, clusterName, clusterType string, enableFunnel bool) {
	log.Println("🚀 Starting Tailscale installation and cluster setup...")

	// Import Tailscale client
	tsClient := tailscale.NewClient()

	// Install Tailscale first
	log.Println("📦 Installing Tailscale...")
	if err := tsClient.InstallTailscale(); err != nil {
		log.Fatalf("❌ Error installing Tailscale: %v", err)
	}

	// Connect to Tailscale
	log.Println("🔗 Connecting to Tailscale...")
	if err := tsClient.Connect(authKey); err != nil {
		log.Fatalf("❌ Error connecting to Tailscale: %v", err)
	}

	// Set environment variables for cluster installation
	envVars := []string{
		fmt.Sprintf("TAILSCALE_AUTH_KEY=%s", authKey),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
		fmt.Sprintf("CLUSTER_TYPE=%s", clusterType),
		fmt.Sprintf("ENABLE_FUNNEL=%t", enableFunnel),
	}

	// Install Kubernetes cluster with Tailscale integration
	installCmd := "curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/tailscale-k8s-setup/main/scripts/install.sh | bash"

	log.Printf("🔧 Installing cluster type: %s", clusterType)
	log.Printf("🌐 Tailscale Funnel: %s", map[bool]string{true: "enabled", false: "disabled"}[enableFunnel])

	// Execute the installer with environment variables
	env := append(os.Environ(), envVars...)
	output, err := utils.RunCommandWithEnv("sh", []string{"-c", installCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error installing cluster with Tailscale: %v\nOutput: %s", err, output)
	}

	// Setup Tailscale Kubernetes operator
	log.Println("🔧 Setting up Tailscale Kubernetes operator...")
	if err := tsClient.SetupKubernetesOperator(); err != nil {
		log.Printf("⚠️ Warning: Failed to setup Tailscale Kubernetes operator: %v", err)
	}

	// Enable Funnel if requested
	if enableFunnel {
		log.Println("🌐 Enabling Tailscale Funnel for port 80...")
		if err := tsClient.EnableFunnel(80); err != nil {
			log.Printf("⚠️ Warning: Failed to enable Funnel: %v", err)
		} else {
			// Get the Funnel URL
			if url, err := tsClient.GetFunnelURL(); err == nil {
				log.Printf("🌍 Your service is now accessible at: %s", url)
			}
		}
	}

	log.Println("✅ Tailscale and cluster setup completed successfully!")
	log.Println("🔗 Your cluster is now connected via Tailscale")

	// Show verification commands
	log.Println("\n📋 Verification commands:")
	log.Println("  tailscale status")
	log.Println("  kubectl get pods -n tailscale-operator")
	if enableFunnel {
		log.Println("  tailscale serve status")
	}
}

// installOnExistingCluster installs Tailscale on an existing Kubernetes cluster
func installOnExistingCluster(cmd *cobra.Command, authKey, clusterName string, enableFunnel bool) {
	log.Println("🚀 Installing Tailscale on existing cluster...")

	// Import Tailscale client
	tsClient := tailscale.NewClient()

	// Install Tailscale first
	log.Println("📦 Installing Tailscale...")
	if err := tsClient.InstallTailscale(); err != nil {
		log.Fatalf("❌ Error installing Tailscale: %v", err)
	}

	// Connect to Tailscale
	log.Println("🔗 Connecting to Tailscale...")
	if err := tsClient.Connect(authKey); err != nil {
		log.Fatalf("❌ Error connecting to Tailscale: %v", err)
	}

	// Setup Tailscale Kubernetes operator
	log.Println("🔧 Setting up Tailscale Kubernetes operator...")
	if err := tsClient.SetupKubernetesOperator(); err != nil {
		log.Fatalf("❌ Error setting up Tailscale Kubernetes operator: %v", err)
	}

	// Enable Funnel if requested
	if enableFunnel {
		log.Println("🌐 Enabling Tailscale Funnel for port 80...")
		if err := tsClient.EnableFunnel(80); err != nil {
			log.Printf("⚠️ Warning: Failed to enable Funnel: %v", err)
		} else {
			// Get the Funnel URL
			if url, err := tsClient.GetFunnelURL(); err == nil {
				log.Printf("🌍 Your service is now accessible at: %s", url)
			}
		}
	}

	log.Println("✅ Tailscale installed on existing cluster!")
	log.Println("🔗 Your cluster is now connected via Tailscale")

	// Show verification commands
	log.Println("\n📋 Verification commands:")
	log.Println("  tailscale status")
	log.Println("  kubectl get pods -n tailscale-operator")
	if enableFunnel {
		log.Println("  tailscale serve status")
	}
}

// updateTailscale updates Tailscale to the latest version
func updateTailscale(cmd *cobra.Command, authKey, clusterName string) {
	log.Println("🔄 Updating Tailscale...")

	// Update Tailscale
	updateCmd := "curl -fsSL https://tailscale.com/install.sh | sh"
	output, err := utils.RunCommand("sh", "-c", updateCmd)
	if err != nil {
		log.Fatalf("❌ Error updating Tailscale: %v\nOutput: %s", err, output)
	}

	log.Println("✅ Tailscale updated successfully!")
}

// uninstallTailscale removes Tailscale and related components
func uninstallTailscale(cmd *cobra.Command) {
	log.Println("🗑️ Uninstalling Tailscale...")

	tsClient := tailscale.NewClient()

	// Disable Funnel first
	log.Println("🌐 Disabling Tailscale Funnel...")
	if err := tsClient.DisableFunnel(); err != nil {
		log.Printf("⚠️ Warning: Failed to disable Funnel: %v", err)
	}

	// Disconnect from Tailscale
	log.Println("🔗 Disconnecting from Tailscale...")
	if err := tsClient.Disconnect(); err != nil {
		log.Printf("⚠️ Warning: Failed to disconnect from Tailscale: %v", err)
	}

	// Uninstall Tailscale
	uninstallCmd := "tailscale uninstall"
	output, err := utils.RunCommand("sh", "-c", uninstallCmd)
	if err != nil {
		log.Fatalf("❌ Error uninstalling Tailscale: %v\nOutput: %s", err, output)
	}

	log.Println("✅ Tailscale uninstalled successfully!")
}

func (a *agentModel) install() {
	// Add flags to the install command
	installCmd.Flags().String("cluster-name", "", "Name for the cluster (default: tailscale-cluster)")
	installCmd.Flags().String("cluster-type", "", "Kubernetes distribution (k3s|minikube|k3d|kind) (default: k3s)")
	installCmd.Flags().Bool("existing-cluster", false, "Install Tailscale on existing Kubernetes cluster")
	installCmd.Flags().Bool("no-funnel", false, "Skip Tailscale Funnel setup (private access only)")
	installCmd.Flags().Bool("update", false, "Update Tailscale to the latest version")
	installCmd.Flags().Bool("uninstall", false, "Uninstall Tailscale and related components")

	a.rootCmd.AddCommand(installCmd)
}
