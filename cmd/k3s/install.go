package k3s

import (
	// "fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install k3s and connect to PipeOps",
	// GroupID: "server",
	Long: `Installs the k3s server and connects it to the PipeOps control plane 
using your service account token.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.ValidateOrPrompt()

		// Validate token argument
		// if len(args) < 1 {
		// 	log.Fatalf("Error: Token is required as an argument.")
		// }
		// token := args[0]

		// Install k3s
		log.Info("Installing k3s...")
		installCmd := "curl -sfL https://get.k3s.io | sh -s -"
		output, err := utils.RunCommand("sh", "-c", installCmd)
		if err != nil {
			log.Fatalf("Error installing k3s: %v\nOutput: %s", err, output)
		}

		log.Info("k3s installed successfully.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// Ensure token is provided
		if len(args) < 1 {
			// return fmt.Errorf("token is required")
		}
		return nil
	},
}

func (k *k3sModel) install() {
	k.rootCmd.AddCommand(installCmd)
}
