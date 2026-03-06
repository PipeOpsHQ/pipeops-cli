package agent

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the PipeOps agent",
	Long: `Restart the PipeOps agent running in the cluster.

This is useful if the agent gets stuck or if you need it to pick up new configurations immediately.
This command wraps 'kubectl rollout restart deployment/pipeops-agent -n pipeops-system'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for kubectl
		if _, err := exec.LookPath("kubectl"); err != nil {
			return fmt.Errorf("kubectl is required to restart the agent but was not found in PATH")
		}

		fmt.Printf("Restarting PipeOps agent...\n\n")

		execCmd := exec.Command("kubectl", "rollout", "restart", "deployment/pipeops-agent", "-n", "pipeops-system")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("failed to restart agent deployment: %w", err)
		}

		fmt.Println("\nAgent restart initiated successfully. You can run 'pipeops agent status' to check its progress.")
		return nil
	},
}

func (a *agentModel) restart() {
	a.rootCmd.AddCommand(restartCmd)
}
