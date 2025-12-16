package agent

import (
	"fmt"
	"log"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show cluster information and join commands",
	Long: `The "info" command displays information about the current PipeOps cluster
including connection details and commands to join additional worker nodes.

This command is useful for:
- Getting the server URL and token for joining worker nodes
- Verifying cluster connectivity
- Displaying cluster status`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Retrieving cluster information...")

		// Run the cluster-info command from the installer
		infoCmd := "curl -fsSL https://raw.githubusercontent.com/PipeOpsHQ/pipeops-k8-agent/main/scripts/install.sh | bash -s -- cluster-info"

		output, err := utils.RunShellCommandWithEnv(infoCmd, nil)
		if err != nil {
			log.Fatalf("Error retrieving cluster information: %v", err)
		}

		log.Println("Cluster information retrieved successfully!")
		fmt.Println(output)

		log.Println("\nTo join additional worker nodes, use:")
		log.Println("  pipeops agent join <server-url> <token>")
	},
}

func (a *agentModel) info() {
	a.rootCmd.AddCommand(infoCmd)
}
