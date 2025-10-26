package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [token]",
	Short: "🚀 Install PipeOps agent and Kubernetes cluster",
	Long: `🚀 The "install" command installs a Kubernetes cluster and the PipeOps agent using the intelligent installer.

This command supports:
- Automatic cluster type detection (k3s, minikube, k3d, kind, auto)
- Full monitoring stack installation (Prometheus, Loki, Grafana, OpenCost)
- Installation on existing Kubernetes clusters
- Environment variable configuration

Examples:
  # Install with token as argument
  pipeops agent install your-token-here

  # Install using environment variables
  export PIPEOPS_TOKEN="your-token"
  export CLUSTER_NAME="my-cluster"
  pipeops agent install

  # Install on existing cluster
  pipeops agent install --existing-cluster --cluster-name="my-existing-cluster"

  # Install without monitoring stack
  pipeops agent install --no-monitoring`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get token from args or environment
		token := getToken(cmd, args)

		// Get cluster name from flag or environment
		clusterName, _ := cmd.Flags().GetString("cluster-name")
		if clusterName == "" {
			clusterName = os.Getenv("CLUSTER_NAME")
		}
		if clusterName == "" {
			clusterName = "pipeops-cluster"
		}

		// Get cluster type from flag or environment
		clusterType, _ := cmd.Flags().GetString("cluster-type")
		if clusterType == "" {
			clusterType = os.Getenv("CLUSTER_TYPE")
		}
		if clusterType == "" {
			clusterType = "auto"
		}

		// Check if installing on existing cluster
		existingCluster, _ := cmd.Flags().GetBool("existing-cluster")

		// Check if monitoring should be disabled
		noMonitoring, _ := cmd.Flags().GetBool("no-monitoring")

		// Check if this is an update operation
		update, _ := cmd.Flags().GetBool("update")

		// Check if this is an uninstall operation
		uninstall, _ := cmd.Flags().GetBool("uninstall")

		if uninstall {
			uninstallAgent(cmd)
			return
		}

		if update {
			updateAgent(cmd, token, clusterName)
			return
		}

		if existingCluster {
			installOnExistingCluster(cmd, token, clusterName)
		} else {
			installNewCluster(cmd, token, clusterName, clusterType, noMonitoring)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// Token can be provided as argument or environment variable
		if len(args) == 0 {
			token := os.Getenv("PIPEOPS_TOKEN")
			if token == "" {
				return fmt.Errorf("❌ token is required either as argument or PIPEOPS_TOKEN environment variable")
			}
		}
		return nil
	},
}

// getToken retrieves token from args or environment
func getToken(cmd *cobra.Command, args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return os.Getenv("PIPEOPS_TOKEN")
}

// installNewCluster installs a new Kubernetes cluster with PipeOps agent
func installNewCluster(cmd *cobra.Command, token, clusterName, clusterType string, noMonitoring bool) {
	log.Println("🚀 Starting PipeOps agent installation...")

	// Set environment variables for the installer
	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
		fmt.Sprintf("CLUSTER_TYPE=%s", clusterType),
	}

	if noMonitoring {
		envVars = append(envVars, "INSTALL_MONITORING=false")
	}

	// Run the PipeOps installer from GitHub
	installCmd := "curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-k8-agent/main/scripts/install.sh | bash"

	log.Printf("🔧 Installing cluster type: %s", clusterType)
	log.Printf("📊 Monitoring stack: %s", map[bool]string{true: "disabled", false: "enabled"}[noMonitoring])

	// Execute the installer with environment variables
	env := append(os.Environ(), envVars...)
	output, err := utils.RunCommandWithEnv("sh", []string{"-c", installCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error installing PipeOps agent: %v\nOutput: %s", err, output)
	}

	log.Println("✅ PipeOps agent installed successfully!")
	log.Println("🔗 Your cluster is now connected to PipeOps control plane")

	// Show verification commands
	log.Println("\n📋 Verification commands:")
	log.Println("  kubectl get pods -n pipeops-system")
	log.Println("  kubectl get pods -n pipeops-monitoring")
}

// installOnExistingCluster installs PipeOps agent on an existing Kubernetes cluster
func installOnExistingCluster(cmd *cobra.Command, token, clusterName string) {
	log.Println("🚀 Installing PipeOps agent on existing cluster...")

	// Set environment variables
	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		fmt.Sprintf("PIPEOPS_CLUSTER_NAME=%s", clusterName),
		"INSTALL_MONITORING=false", // Skip monitoring for existing clusters by default
	}

	// Download and apply the agent manifest
	agentCmd := `curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-k8-agent/main/deployments/agent.yaml \
  | sed "s/PIPEOPS_TOKEN: \"your-token-here\"/PIPEOPS_TOKEN: \"${PIPEOPS_TOKEN}\"/" \
  | sed "s/token: \"your-token-here\"/token: \"${PIPEOPS_TOKEN}\"/" \
  | sed "s/cluster_name: \"default-cluster\"/cluster_name: \"${PIPEOPS_CLUSTER_NAME}\"/" \
  | kubectl apply -f -`

	log.Printf("🔧 Applying agent manifest to cluster: %s", clusterName)

	env := append(os.Environ(), envVars...)
	output, err := utils.RunCommandWithEnv("sh", []string{"-c", agentCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error installing agent on existing cluster: %v\nOutput: %s", err, output)
	}

	log.Println("✅ PipeOps agent installed on existing cluster!")
	log.Println("🔗 Your cluster is now connected to PipeOps control plane")

	// Show verification commands
	log.Println("\n📋 Verification commands:")
	log.Println("  kubectl rollout status deployment/pipeops-agent -n pipeops-system")
	log.Println("  kubectl logs deployment/pipeops-agent -n pipeops-system")
}

// updateAgent updates the PipeOps agent to the latest version
func updateAgent(cmd *cobra.Command, token, clusterName string) {
	log.Println("🔄 Updating PipeOps agent...")

	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
	}

	updateCmd := "curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-k8-agent/main/scripts/install.sh | bash -s -- update"

	env := append(os.Environ(), envVars...)
	output, err := utils.RunCommandWithEnv("sh", []string{"-c", updateCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error updating PipeOps agent: %v\nOutput: %s", err, output)
	}

	log.Println("✅ PipeOps agent updated successfully!")
}

// uninstallAgent removes the PipeOps agent and monitoring stack
func uninstallAgent(cmd *cobra.Command) {
	log.Println("🗑️ Uninstalling PipeOps agent...")

	uninstallCmd := "curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-k8-agent/main/scripts/install.sh | bash -s -- uninstall"

	output, err := utils.RunCommand("sh", "-c", uninstallCmd)
	if err != nil {
		log.Fatalf("❌ Error uninstalling PipeOps agent: %v\nOutput: %s", err, output)
	}

	log.Println("✅ PipeOps agent uninstalled successfully!")
}

func (a *agentModel) install() {
	// Add flags to the install command
	installCmd.Flags().String("cluster-name", "", "Name for the cluster (default: pipeops-cluster)")
	installCmd.Flags().String("cluster-type", "", "Kubernetes distribution (k3s|minikube|k3d|kind|auto) (default: auto)")
	installCmd.Flags().Bool("existing-cluster", false, "Install agent on existing Kubernetes cluster")
	installCmd.Flags().Bool("no-monitoring", false, "Skip installation of monitoring stack (Prometheus, Loki, Grafana)")
	installCmd.Flags().Bool("update", false, "Update the agent to the latest version")
	installCmd.Flags().Bool("uninstall", false, "Uninstall the agent and monitoring stack")

	a.rootCmd.AddCommand(installCmd)
}
