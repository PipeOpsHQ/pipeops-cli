package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/libs"
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
			clusterType = "k3s" // Default to k3s
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
			uninstallAgent(cmd, token)
			return
		}

		if update {
			updateAgent(cmd, token, clusterName)
			return
		}

		if existingCluster {
			installOnExistingCluster(cmd, token, clusterName, !noMonitoring)
		} else {
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

			return fmt.Errorf("❌ PipeOps authentication is required. Use 'pipeops auth login' or provide PIPEOPS_TOKEN environment variable")
		}
		return nil
	},
}

// getPipeOpsToken retrieves PipeOps token from args, environment, or config
func getPipeOpsToken(cmd *cobra.Command, args []string) string {
	// Check arguments first
	if len(args) > 0 {
		return args[0]
	}

	// Check environment variable
	if token := os.Getenv("PIPEOPS_TOKEN"); token != "" {
		return token
	}

	// Check OAuth config
	cfg, err := config.Load()
	if err == nil && cfg.IsAuthenticated() {
		return cfg.OAuth.AccessToken
	}

	return ""
}

// installNewCluster installs a new Kubernetes cluster with PipeOps agent
func installNewCluster(cmd *cobra.Command, token, clusterName, clusterType string, enableMonitoring bool) {
	    // Validate token
	    if err := validateToken(token); err != nil {
	        log.Printf("⚠️ Warning: Token validation skipped: %v", err)
	    }
	// Set environment variables for cluster installation
	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
		fmt.Sprintf("CLUSTER_TYPE=%s", clusterType),
		fmt.Sprintf("ENABLE_MONITORING=%t", enableMonitoring),
	}

	// Install Kubernetes cluster with PipeOps agent integration
	installCmd := "curl -fsSL https://get.pipeops.dev | bash"

	log.Printf("Installing cluster type: %s", clusterType)
	log.Printf("PipeOps monitoring: %s", map[bool]string{true: "enabled", false: "disabled"}[enableMonitoring])

	// Execute the installer with environment variables
	env := append(os.Environ(), envVars...)
	output, err := utils.RunCommandWithEnv("sh", []string{"-c", installCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error installing cluster with PipeOps agent: %v\nOutput: %s", err, output)
	}

	// Setup PipeOps Kubernetes agent
	log.Println("Setting up PipeOps Kubernetes agent...")
	if err := setupPipeOpsAgent(token, clusterName); err != nil {
		log.Printf("Warning: Failed to setup PipeOps agent: %v", err)
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

	// Validate token
	if err := validateToken(token); err != nil {
		log.Printf("⚠️ Warning: Token validation skipped: %v", err)
	}

	// The agent install script handles everything, including existing clusters
	installCmd := "curl -fsSL https://get.pipeops.dev | bash"
	
	// Set environment variables
	envVars := []string{
		fmt.Sprintf("PIPEOPS_TOKEN=%s", token),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
		fmt.Sprintf("ENABLE_MONITORING=%t", enableMonitoring),
	}
	
	env := append(os.Environ(), envVars...)
	output, err := utils.RunCommandWithEnv("sh", []string{"-c", installCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error installing PipeOps agent: %v\nOutput: %s", err, output)
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

	// Validate token
	if err := validateToken(token); err != nil {
		log.Printf("⚠️ Warning: Token validation skipped: %v", err)
	}

	// Update PipeOps agent
	updateCmd := "curl -fsSL https://get.pipeops.dev | bash"
	envVars := []string{fmt.Sprintf("PIPEOPS_TOKEN=%s", token)}
	env := append(os.Environ(), envVars...)

	output, err := utils.RunCommandWithEnv("sh", []string{"-c", updateCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error updating PipeOps agent: %v\nOutput: %s", err, output)
	}

	log.Println("PipeOps agent updated successfully!")
}

// uninstallAgent removes PipeOps agent and related components
func uninstallAgent(cmd *cobra.Command, token string) {
	log.Println("Uninstalling PipeOps agent...")

	// Validate token
	if err := validateToken(token); err != nil {
		log.Printf("⚠️ Warning: Token validation skipped: %v", err)
	}

	// Remove monitoring first
	log.Println("Removing PipeOps monitoring...")
	if err := removeMonitoring(); err != nil {
		log.Printf("Warning: Failed to remove monitoring: %v", err)
	}

	// Remove PipeOps agent
	log.Println("Removing PipeOps agent...")
	if err := removePipeOpsAgent(); err != nil {
		log.Printf("Warning: Failed to remove agent: %v", err)
	}

	// Uninstall PipeOps agent
	uninstallCmd := "curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-agent/main/scripts/uninstall.sh | bash"
	envVars := []string{fmt.Sprintf("PIPEOPS_TOKEN=%s", token)}
	env := append(os.Environ(), envVars...)

	output, err := utils.RunCommandWithEnv("sh", []string{"-c", uninstallCmd}, env)
	if err != nil {
		log.Fatalf("❌ Error uninstalling PipeOps agent: %v\nOutput: %s", err, output)
	}

	log.Println("PipeOps agent uninstalled successfully!")
}

// Helper functions

// validateToken validates the PipeOps token
func validateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}

	// Use the libs HTTP client to verify the token
	httpClient := libs.NewHttpClient()
	_, err := httpClient.VerifyToken(token, "")
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	return nil
}

// setupPipeOpsAgent sets up the PipeOps agent on the cluster
func setupPipeOpsAgent(token, clusterName string) error {
	log.Printf("Installing PipeOps agent for cluster: %s", clusterName)

	// Apply PipeOps agent manifests
	setupCmd := fmt.Sprintf(`
kubectl create namespace pipeops-system --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic pipeops-token -n pipeops-system --from-literal=token=%s --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f https://raw.githubusercontent.com/PipeOpsHQ/pipeops-agent/main/manifests/agent.yaml
`, token)

	output, err := utils.RunCommand("sh", "-c", setupCmd)
	if err != nil {
		return fmt.Errorf("failed to setup agent: %v\nOutput: %s", err, output)
	}

	log.Println("PipeOps agent setup completed")
	return nil
}

// setupMonitoring sets up monitoring components
func setupMonitoring(token, clusterName string) error {
	log.Printf("Installing monitoring for cluster: %s", clusterName)

	// Apply monitoring manifests
	monitoringCmd := fmt.Sprintf(`
kubectl create namespace pipeops-monitoring --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic pipeops-token -n pipeops-monitoring --from-literal=token=%s --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f https://raw.githubusercontent.com/PipeOpsHQ/pipeops-agent/main/manifests/monitoring.yaml
`, token)

	output, err := utils.RunCommand("sh", "-c", monitoringCmd)
	if err != nil {
		return fmt.Errorf("failed to setup monitoring: %v\nOutput: %s", err, output)
	}

	log.Println("Monitoring setup completed")
	return nil
}

// removePipeOpsAgent removes the PipeOps agent
func removePipeOpsAgent() error {
	removeCmd := `
kubectl delete -f https://raw.githubusercontent.com/PipeOpsHQ/pipeops-agent/main/manifests/agent.yaml --ignore-not-found=true
kubectl delete secret pipeops-token -n pipeops-system --ignore-not-found=true
kubectl delete namespace pipeops-system --ignore-not-found=true
`

	output, err := utils.RunCommand("sh", "-c", removeCmd)
	if err != nil {
		return fmt.Errorf("failed to remove agent: %v\nOutput: %s", err, output)
	}

	return nil
}

// removeMonitoring removes monitoring components
func removeMonitoring() error {
	removeCmd := `
kubectl delete -f https://raw.githubusercontent.com/PipeOpsHQ/pipeops-agent/main/manifests/monitoring.yaml --ignore-not-found=true
kubectl delete secret pipeops-token -n pipeops-monitoring --ignore-not-found=true
kubectl delete namespace pipeops-monitoring --ignore-not-found=true
`

	output, err := utils.RunCommand("sh", "-c", removeCmd)
	if err != nil {
		return fmt.Errorf("failed to remove monitoring: %v\nOutput: %s", err, output)
	}

	return nil
}

func (a *agentModel) install() {
	// Add flags to the install command
	installCmd.Flags().String("cluster-name", "", "Name for the cluster (default: pipeops-cluster)")
	installCmd.Flags().String("cluster-type", "", "Kubernetes distribution (k3s|minikube|k3d|kind) (default: k3s)")
	installCmd.Flags().Bool("existing-cluster", false, "Install PipeOps agent on existing Kubernetes cluster")
	installCmd.Flags().Bool("no-monitoring", false, "Skip monitoring setup (agent only)")
	installCmd.Flags().Bool("update", false, "Update PipeOps agent to the latest version")
	installCmd.Flags().Bool("uninstall", false, "Uninstall PipeOps agent and related components")

	a.rootCmd.AddCommand(installCmd)
}
