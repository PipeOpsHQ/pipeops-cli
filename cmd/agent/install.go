package cmd

import (
	// "fmt"
	"log"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
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
		log.Println("Installing k3s...")
		installCmd := "curl -sfL https://get.k3s.io | sh -s -"
		output, err := utils.RunCommand("sh", "-c", installCmd)
		if err != nil {
			log.Fatalf("Error installing k3s: %v\nOutput: %s", err, output)
		}
		log.Println("k3s installed successfully.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// Ensure token is provided
		if len(args) < 1 {
			// return fmt.Errorf("token is required")
		}
		return nil
	},
}

func (a *agentModel) install() {
	a.rootCmd.AddCommand(installCmd)
}
