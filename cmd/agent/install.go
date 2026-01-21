package agent

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/sanitize"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [pipeops-token]",
	Short: "Install PipeOps agent and configure Kubernetes cluster",
	Long: `The "install" command installs the PipeOps agent on your Kubernetes cluster for monitoring and management.

This command supports:
	Automatic PipeOps agent installation
	Kubernetes cluster setup with PipeOps integration
	Agent configuration and monitoring setup
	Automatic cluster detection and registration


Examples:
  # Install with token as argument

  pipeops agent install your-pipeops-token-here

  # Install using environment variables
  export PIPEOPS_TOKEN="your-pipeops-token-here"

  export CLUSTER_NAME="my-cluster"
  pipeops agent install


  # Install on existing cluster
  pipeops agent install --existing-cluster --cluster-name="my-existing-cluster"

  # Install without monitoring (basic setup only)
  pipeops agent install --no-monitoring`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("[DEBUG] Starting install command execution")
		// Get PipeOps token from args, environment, or config
		token := getPipeOpsToken(cmd, args)

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
			clusterType = "auto" // Let installer pick the best option
		}

		// Check if installing on existing cluster
		existingCluster, _ := cmd.Flags().GetBool("existing-cluster")

		// Check if monitoring should be disabled
		noMonitoring, _ := cmd.Flags().GetBool("no-monitoring")

		// Check if this is an update operation
		update, _ := cmd.Flags().GetBool("update")

		if update {
			log.Println("[DEBUG] Running update agent")
			updateAgent(cmd, token, clusterName)
			return
		}

		if existingCluster {
			log.Println("[DEBUG] Running install on existing cluster")
			installOnExistingCluster(cmd, token, clusterName, !noMonitoring)
		} else {
			log.Println("[DEBUG] Running install on new cluster")
			installNewCluster(cmd, token, clusterName, clusterType, !noMonitoring)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// PipeOps token can be provided as argument, environment variable, or from config
		if len(args) == 0 {
			// Check environment variable
			if token := os.Getenv("PIPEOPS_TOKEN"); token != "" {
				return nil
			}

			// Check if user is authenticated via OAuth
			cfg, err := config.Load()
			if err == nil && cfg.IsAuthenticated() {
				return nil
			}

			return fmt.Errorf("Error: PipeOps token is required. Provide token as argument: 'pipeops agent install <token>' or set PIPEOPS_TOKEN environment variable")
		}
		return nil
	},
}

// getPipeOpsToken retrieves PipeOps token from args, environment, or config
func getPipeOpsToken(cmd *cobra.Command, args []string) string {
	// Check arguments first
	if len(args) > 0 {
		return strings.TrimSpace(args[0])
	}

	// Check environment variable
	if token := os.Getenv("PIPEOPS_TOKEN"); token != "" {
		return strings.TrimSpace(token)
	}

	// Check OAuth config
	cfg, err := config.Load()
	if err == nil && cfg.IsAuthenticated() {
		return strings.TrimSpace(cfg.OAuth.AccessToken)
	}

	return ""
}

// installNewCluster installs a new Kubernetes cluster with PipeOps agent
func installNewCluster(cmd *cobra.Command, token, clusterName, clusterType string, enableMonitoring bool) {
	if token == "" {
		log.Fatalf("Error: PipeOps token is required. Please provide it as an argument or set PIPEOPS_TOKEN environment variable.")
	}

	// Set environment variables for cluster installation
	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
		fmt.Sprintf("CLUSTER_TYPE=%s", clusterType),
		fmt.Sprintf("INSTALL_MONITORING=%t", enableMonitoring),
	}

	// Install Kubernetes cluster with PipeOps agent integration
	installCmd := "curl -fsSL https://get.pipeops.dev/k8-install.sh | bash"

	log.Printf("Installing cluster type: %s", sanitize.Log(clusterType))
	log.Printf("PipeOps monitoring: %s", map[bool]string{true: "enabled", false: "disabled"}[enableMonitoring])

	// Execute the installer with environment variables
	log.Println("[DEBUG] Executing installer script")
	_, err := utils.RunShellCommandWithEnvStreaming(installCmd, envVars)
	if err != nil {
		log.Fatalf("Error installing cluster with PipeOps agent: %v", err)
	}

	log.Println("PipeOps agent and cluster setup completed successfully!")
	log.Println("Your cluster is now connected to PipeOps")

	// Show verification commands
	log.Println("\nVerification commands:")
	log.Println("  kubectl get pods -n pipeops-system")
	log.Println("  pipeops server list")
	if enableMonitoring {
		log.Println("  kubectl get pods -n pipeops-monitoring")
	}
}

// installOnExistingCluster installs PipeOps agent on an existing Kubernetes cluster
func installOnExistingCluster(cmd *cobra.Command, token, clusterName string, enableMonitoring bool) {
	log.Println("Installing PipeOps agent on existing cluster...")

	if token == "" {
		log.Fatalf("Error: PipeOps token is required. Please provide it as an argument or set PIPEOPS_TOKEN environment variable.")
	}

	// The agent install script handles everything, including existing clusters
	installCmd := "curl -fsSL https://get.pipeops.dev/k8-install.sh | bash"

	// Set environment variables
	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
		fmt.Sprintf("INSTALL_MONITORING=%t", enableMonitoring),
	}

	log.Println("[DEBUG] Executing installer script")
	_, err := utils.RunShellCommandWithEnvStreaming(installCmd, envVars)
	if err != nil {
		log.Fatalf("Error installing PipeOps agent: %v", err)
	}

	log.Println("PipeOps agent installed on existing cluster!")
	log.Println("Your cluster is now connected to PipeOps")

	// Show verification commands
	log.Println("\nVerification commands:")
	log.Println("  kubectl get pods -n pipeops-system")
	log.Println("  pipeops server list")
	if enableMonitoring {
		log.Println("  kubectl get pods -n pipeops-monitoring")
	}
}

// updateAgent updates PipeOps agent to the latest version
func updateAgent(cmd *cobra.Command, token, clusterName string) {
	log.Println("Updating PipeOps agent...")

	// Update PipeOps agent
	updateCmd := "curl -fsSL https://get.pipeops.dev/k8-install.sh | bash"
	
	clusterType, _ := cmd.Flags().GetString("cluster-type")
	
	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		"UPDATE=true",
	}
	
	if clusterName != "" {
		envVars = append(envVars, fmt.Sprintf("CLUSTER_NAME=%s", clusterName))
	}
	if clusterType != "" {
		envVars = append(envVars, fmt.Sprintf("CLUSTER_TYPE=%s", clusterType))
	}

	_, err := utils.RunShellCommandWithEnvStreaming(updateCmd, envVars)
	if err != nil {
		log.Fatalf("Error updating PipeOps agent: %v", err)
	}

	log.Println("PipeOps agent updated successfully!")
}





func (a *agentModel) install() {
	// Add flags to the install command
	installCmd.Flags().String("cluster-name", "", "Name for the cluster (default: pipeops-cluster)")
	installCmd.Flags().String("cluster-type", "", "Kubernetes distribution (k3s|minikube|k3d|kind|auto) (default: auto)")
	installCmd.Flags().Bool("existing-cluster", false, "Install PipeOps agent on existing Kubernetes cluster")
	installCmd.Flags().Bool("no-monitoring", false, "Skip monitoring setup (agent only)")
	installCmd.Flags().Bool("update", false, "Update PipeOps agent to the latest version")
	a.rootCmd.AddCommand(installCmd)
}
