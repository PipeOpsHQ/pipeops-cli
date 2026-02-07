package agent

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View PipeOps agent logs",
	Long: `View and stream logs from the PipeOps agent running in your Kubernetes cluster.

This command wraps kubectl to fetch logs from the pipeops-agent pod in the pipeops-system namespace.
It automatically finds the correct pod and streams logs.

Examples:
  - View recent logs:
    pipeops agent logs

  - Stream logs in real-time:
    pipeops agent logs -f

  - View logs with tail:
    pipeops agent logs --tail=100`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check for kubectl
		if _, err := exec.LookPath("kubectl"); err != nil {
			log.Fatalf("Error: kubectl is required to view agent logs but was not found in PATH.")
		}

		follow, _ := cmd.Flags().GetBool("follow")
		tail, _ := cmd.Flags().GetInt("tail")

		// Construct kubectl command to get pod name
		// We use a shell command to handle the subshell execution nicely or we can do it in two steps in Go.
		// Doing it in two steps in Go is cleaner and safer.

		log.Println("Finding PipeOps agent pod...")
		
		// 1. Find the agent pod name
		// kubectl get pods -n pipeops-system -l app=pipeops-agent -o jsonpath="{.items[0].metadata.name}"
		findPodCmd := exec.Command("kubectl", "get", "pods", "-n", "pipeops-system", "-l", "app=pipeops-agent", "-o", "jsonpath={.items[0].metadata.name}")
		output, err := findPodCmd.Output()
		if err != nil {
			log.Fatalf("Failed to find PipeOps agent pod: %v. Is the agent installed and running in 'pipeops-system' namespace?", err)
		}
		podName := string(output)
		if podName == "" {
			log.Fatalf("No PipeOps agent pod found in 'pipeops-system' namespace with label 'app=pipeops-agent'.")
		}

		log.Printf("Found agent pod: %s", podName)

		// 2. Stream logs
		// kubectl logs -n pipeops-system <podName> [-f] [--tail=n]
		kubectlArgs := []string{"logs", "-n", "pipeops-system", podName}
		
		if follow {
			kubectlArgs = append(kubectlArgs, "-f")
		}
		
		if tail > 0 {
			kubectlArgs = append(kubectlArgs, fmt.Sprintf("--tail=%d", tail))
		} else if !follow {
			// Default tail if not following and no tail specified, to avoid dumping massive logs
			// But kubectl logs defaults to all logs. Let's keep kubectl default behavior unless specified.
		}

		// Use utils.RunCommandWithEnvStreaming to execute kubectl logs and stream output to user
		logCommand := fmt.Sprintf("kubectl %s", fmt.Sprintf("logs -n pipeops-system %s", podName))
		if follow {
			logCommand += " -f"
		}
		if tail > 0 {
			logCommand += fmt.Sprintf(" --tail=%d", tail)
		}
		
		// For the actual execution, we can just use the args directly with os/exec to connect streams
		// creating an interactive experience (Ctrl+C works properly)
		cmdLog := exec.Command("kubectl", kubectlArgs...)
		
		// Connect streams directly
		cmdLog.Stdout = cmd.OutOrStdout()
		cmdLog.Stderr = cmd.OutOrStderr()
		
		log.Printf("Fetching logs from %s...", podName)
		if err := cmdLog.Run(); err != nil {
			// Don't fatal here as Ctrl+C might cause a non-zero exit which is fine for -f
			if follow {
				return
			}
			log.Fatalf("Error streaming logs: %v", err)
		}
	},
}

func (a *agentModel) logs() {
	logsCmd.Flags().BoolP("follow", "f", false, "Stream logs in real-time")
	logsCmd.Flags().Int("tail", -1, "Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided.")
	a.rootCmd.AddCommand(logsCmd)
}
