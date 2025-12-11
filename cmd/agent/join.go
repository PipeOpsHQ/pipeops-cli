package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join <server-url> <token>",
	Short: "Join worker node to existing PipeOps cluster",
	Long: `The "join" command joins a worker node to an existing PipeOps-managed Kubernetes cluster.

This command downloads and runs the join-worker.sh script to add the current machine
as a worker node to an existing cluster.

Examples:
  # Join with server URL and token
  pipeops agent join https://192.168.1.100:6443 abc123def456

  # Join using environment variables
  export K3S_URL="https://192.168.1.100:6443"
  export K3S_TOKEN="abc123def456"
  pipeops agent join`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get server URL and token from args or environment
		serverURL := getServerURL(cmd, args)
		token := getJoinToken(cmd, args)

		log.Println("Joining worker node to PipeOps cluster...")
		log.Printf("Server URL: %s", serverURL)

		// Set environment variables for the join script
		envVars := []string{
			fmt.Sprintf("K3S_URL=%s", serverURL),
			fmt.Sprintf("K3S_TOKEN=%s", token),
		}

		// Run the join-worker script from GitHub
			joinCmd := "curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-k8-agent/main/scripts/join-worker.sh | bash"

			env := append(os.Environ(), envVars...)
			_, err := utils.RunCommandWithEnvStreaming("sh", []string{"-c", joinCmd}, env)
			if err != nil {
				log.Fatalf("❌ Error joining worker node")
			}

		log.Println("Worker node joined successfully!")
		log.Println("This node is now part of the PipeOps cluster")

		// Show verification commands
		log.Println("\nVerification commands:")
		log.Println("  kubectl get nodes")
		log.Println("  kubectl get pods -n pipeops-system")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// Check if we have args or environment variables
		if len(args) == 0 {
			serverURL := os.Getenv("K3S_URL")
			token := os.Getenv("K3S_TOKEN")
			if serverURL == "" || token == "" {
				return fmt.Errorf("server URL and token are required either as arguments or K3S_URL/K3S_TOKEN environment variables")
			}
		} else if len(args) < 2 {
			return fmt.Errorf("server URL and token are required as arguments")
		}
		return nil
	},
}

// getServerURL retrieves server URL from args or environment
func getServerURL(cmd *cobra.Command, args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return os.Getenv("K3S_URL")
}

// getJoinToken retrieves join token from args or environment
func getJoinToken(cmd *cobra.Command, args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	return os.Getenv("K3S_TOKEN")
}

func (a *agentModel) join() {
	a.rootCmd.AddCommand(joinCmd)
}
