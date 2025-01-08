// cmd/restart.go
package k3s

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the k3s service",
	// GroupID: "server",
	Long: `Restarts the k3s service, allowing the cluster to recover from any issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Restarting k3s service...")

		output, err := utils.RunCommand("systemctl", "restart", "k3s")
		if err != nil {
			log.Fatalf("Error restarting k3s: %v\nOutput: %s", err, output)
		}

		log.Info("k3s service restarted.")
	},
}

func (k *k3sModel) restart() {
	k.rootCmd.AddCommand(restartCmd)
}
