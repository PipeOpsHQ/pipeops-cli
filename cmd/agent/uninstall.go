package agent

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Aliases: []string{"remove", "rm"},
	Short:   "Uninstall PipeOps agent and destroy the cluster",
	Long: `The "uninstall" command removes the PipeOps agent and destroys the Kubernetes cluster created by PipeOps.

WARNING: This action is irreversible. It will remove the PipeOps agent and delete the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			if !confirmUninstall() {
				fmt.Println("Uninstall cancelled.")
				return
			}
		}

		executeUninstall(cmd)
	},
}

func confirmUninstall() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("WARNING: This will destroy the PipeOps agent and the Kubernetes cluster. Are you sure? (y/N): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func executeUninstall(cmd *cobra.Command) {
	log.Println("Uninstalling PipeOps agent and destroying cluster...")

	uninstallScript := "curl -fsSL https://get.pipeops.dev/k8-uninstall.sh | bash -s -- --force"

	// Gather environment variables
	var envVars []string

	// 1. Token
	token := os.Getenv("PIPEOPS_TOKEN")
	if token == "" {
		// Try to load from config
		cfg, err := config.Load()
		if err == nil && cfg.IsAuthenticated() {
			token = cfg.OAuth.AccessToken
		}
	}
	if token != "" {
		envVars = append(envVars, fmt.Sprintf("PIPEOPS_TOKEN=%s", token))
	}

	// 2. Cluster Name
	clusterName, _ := cmd.Flags().GetString("cluster-name")
	if clusterName == "" {
		clusterName = os.Getenv("CLUSTER_NAME")
	}
	if clusterName != "" {
		envVars = append(envVars, fmt.Sprintf("CLUSTER_NAME=%s", clusterName))
	}

	// 3. Cluster Type
	clusterType, _ := cmd.Flags().GetString("cluster-type")
	if clusterType == "" {
		clusterType = os.Getenv("CLUSTER_TYPE")
	}
	if clusterType != "" {
		envVars = append(envVars, fmt.Sprintf("CLUSTER_TYPE=%s", clusterType))
	}

	_, err := utils.RunShellCommandWithEnvStreaming(uninstallScript, envVars)
	if err != nil {
		log.Fatalf("Error uninstalling PipeOps agent: %v", err)
	}

	// Run k3s-uninstall.sh if it exists to fully cleanup
	k3sUninstallPath := "/usr/local/bin/k3s-uninstall.sh"
	if _, err := os.Stat(k3sUninstallPath); err == nil {
		log.Println("Running k3s cleanup script...")
		_, err := utils.RunShellCommandWithEnvStreaming(k3sUninstallPath, nil)
		if err != nil {
			log.Printf("Warning: Failed to run k3s uninstall script: %v", err)
		} else {
			log.Println("k3s cleanup completed successfully.")
		}
	}

	log.Println("PipeOps agent uninstalled and cluster destroyed successfully!")
}

func (a *agentModel) uninstall() {
	uninstallCmd.Flags().Bool("force", false, "Skip confirmation prompt")
	uninstallCmd.Flags().String("cluster-name", "", "Name of the cluster to destroy")
	uninstallCmd.Flags().String("cluster-type", "", "Type of the cluster (k3s|minikube|k3d|kind|auto)")
	a.rootCmd.AddCommand(uninstallCmd)
}