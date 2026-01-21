package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update PipeOps agent to the latest version",
	Long: `The "update" command updates the PipeOps agent installed on your Kubernetes cluster to the latest available version.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := getPipeOpsToken(cmd, args)
		if token == "" {
			log.Fatalf("Error: PipeOps token is required. Please login or set PIPEOPS_TOKEN environment variable.")
		}

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		if clusterName == "" {
			clusterName = os.Getenv("CLUSTER_NAME")
		}

		clusterType, _ := cmd.Flags().GetString("cluster-type")
		if clusterType == "" {
			clusterType = os.Getenv("CLUSTER_TYPE")
		}

		log.Println("Updating PipeOps agent...")

		updateScript := "curl -fsSL https://get.pipeops.dev/k8-install.sh | bash"
		
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

		_, err := utils.RunShellCommandWithEnvStreaming(updateScript, envVars)
		if err != nil {
			log.Fatalf("Error updating PipeOps agent: %v", err)
		}

		log.Println("PipeOps agent updated successfully!")
	},
}
