package agent

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the PipeOps agent",
	Long: `Check if the PipeOps agent is running and healthy in the cluster.

This command wraps 'kubectl get pods -n pipeops-system' to quickly view the agent's status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for kubectl
		if _, err := exec.LookPath("kubectl"); err != nil {
			return fmt.Errorf("kubectl is required to check agent status but was not found in PATH")
		}

		fmt.Printf("Checking PipeOps agent status...\n\n")

		execCmd := exec.Command("kubectl", "get", "pods", "-n", "pipeops-system", "-l", "app=pipeops-agent")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("failed to get agent status: %w", err)
		}

		return nil
	},
}

func (a *agentModel) status() {
	a.rootCmd.AddCommand(statusCmd)
}
