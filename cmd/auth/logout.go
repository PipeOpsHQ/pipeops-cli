package auth

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "🚪 Logout from your PipeOps account",
	Long: `🚪 The "logout" command removes your authentication token and logs you out
from your PipeOps account. You will need to login again to use other CLI commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if user is authenticated
		if !client.IsAuthenticated() {
			fmt.Println("ℹ️  You are not currently logged in.")
			return
		}

		// Clear the token
		client.SetToken("")

		// Save the configuration
		if err := client.SaveConfig(); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			return
		}

		fmt.Println("✅ Successfully logged out!")
		fmt.Println("🔐 Your authentication token has been removed.")
		fmt.Println("💡 Run 'pipeops auth login' to authenticate again.")
	},
	Args: cobra.NoArgs,
}

func (a *authModel) logout() {
	a.rootCmd.AddCommand(logoutCmd)
}
