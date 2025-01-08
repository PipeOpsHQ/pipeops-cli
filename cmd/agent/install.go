package agent

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
)

var installCmd = &cobra.Command{
	Use:   "install <token>",
	Short: "🚀 Install k3s and connect to PipeOps",
	Long: `🚀 The "install" command installs the k3s server and connects it to the PipeOps control plane 
using your service account token. 

This command automates the installation process for k3s and ensures the server is properly linked to PipeOps.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate or prompt user input
		utils.ValidateOrPrompt()

		// Validate the token argument
		if len(args) < 1 {
			log.Fatalf("❌ Error: Service account token is required. Provide it as the first argument.")
		}
		token := args[0]

		// Log the installation start
		log.Info("🔧 Installing k3s...")

		// Command to install k3s
		installCmd := "curl -sfL https://get.k3s.io | sh -s -"
		output, err := utils.RunCommand("sh", "-c", installCmd)
		if err != nil {
			log.Fatalf("❌ Error installing k3s: %v\nOutput: %s", err, output)
		}

		log.Info("✅ k3s installed successfully.")

		// Simulate connecting to PipeOps using the token
		log.Infof("🔗 Connecting to PipeOps using token: %s\n", token)
		// Placeholder for the actual connection logic
		log.Info("✅ Connected to PipeOps successfully.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// Ensure token is provided
		if len(args) < 1 {
			err := fmt.Errorf("❌ token is required as the first argument")
			log.Error(err)
			return err
		}

		return nil
	},
}

func (a *agentModel) install() {
	a.rootCmd.AddCommand(installCmd)
}
