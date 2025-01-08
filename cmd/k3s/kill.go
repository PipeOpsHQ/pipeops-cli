// cmd/kill.go
package k3s

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill the k3s service",
	// GroupID: "server",
	Long: `Stops the k3s service gracefully, effectively killing the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Stopping k3s service...")

		output, err := utils.RunCommand("systemctl", "stop", "k3s")
		if err != nil {
			log.Fatalf("Error stopping k3s: %v\nOutput: %s", err, output)
		}

		log.Println("k3s service stopped.")
	},
}

func (k *k3sModel) kill() {
	k.rootCmd.AddCommand(killCmd)
}
