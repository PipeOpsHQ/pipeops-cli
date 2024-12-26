// cmd/restart.go
package cmd

import (
	"log"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Restart the k3s service",
	// GroupID: "server",
	Long: `Restarts the k3s service, allowing the cluster to recover from any issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Restarting k3s service...")
		output, err := utils.RunCommand("systemctl", "restart", "k3s")
		if err != nil {
			log.Fatalf("Error restarting k3s: %v\nOutput: %s", err, output)
		}
		log.Println("k3s service restarted.")
	},
}

func (k *k3sModel) restart() {
	k.rootCmd.AddCommand(restartCmd)
}
