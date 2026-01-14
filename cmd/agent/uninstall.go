package agent

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall PipeOps agent and destroy the cluster",
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

		executeUninstall()
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

func executeUninstall() {
	log.Println("Uninstalling PipeOps agent and destroying cluster...")

	uninstallScript := "curl -fsSL https://get.pipeops.dev/k8-uninstall.sh | bash -s -- --force"

	// Pass PIPEOPS_TOKEN if present in environment, though the script might not strictly need it for local destruction
	token := os.Getenv("PIPEOPS_TOKEN")
	var envVars []string
	if token != "" {
		envVars = append(envVars, fmt.Sprintf("PIPEOPS_TOKEN=%s", token))
	}

	_, err := utils.RunShellCommandWithEnvStreaming(uninstallScript, envVars)
	if err != nil {
		log.Fatalf("Error uninstalling PipeOps agent: %v", err)
	}

	log.Println("PipeOps agent uninstalled and cluster destroyed successfully!")
}

func (a *agentModel) uninstall() {
	uninstallCmd.Flags().Bool("force", false, "Skip confirmation prompt")
	a.rootCmd.AddCommand(uninstallCmd)
}
