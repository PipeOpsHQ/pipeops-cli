package k3s

import (
	"fmt"
	"log"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var joinCmd = &cobra.Command{
	Use:   "join [server-url]",
	Short: "Join a worker node to the k3s cluster",
	// GroupID: "server",
	Long: `Joins the current node as a worker to an existing k3s cluster using the provided server URL.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverURL := args[0]
		if !utils.IsValidURL(serverURL) {
			log.Fatalf("Invalid server URL: %s", serverURL)
		}

		joinCommand := fmt.Sprintf("curl -sfL https://get.k3s.io | K3S_URL=%s K3S_TOKEN=%s sh -", serverURL, viper.Get("service_account_token"))
		log.Println("Joining the k3s cluster...")
		output, err := utils.RunCommand("sh", "-c", joinCommand)
		if err != nil {
			log.Fatalf("Error joining k3s cluster: %v\nOutput: %s", err, output)
		}
		log.Println("Successfully joined the k3s cluster.")
	},
}

func (k *k3sModel) join() {
	k.rootCmd.AddCommand(joinCmd)
}
