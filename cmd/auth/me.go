package auth

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "me",
	Short: "👤 Get details of the currently logged-in user",
	Long: `👤 The "me" command retrieves and displays details about the currently logged-in user
in the PipeOps platform. Use this to confirm your active user identity and associated account.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if user is authenticated
		if !client.IsAuthenticated() {
			fmt.Println("❌ You are not logged in. Please run 'pipeops auth login' first.")
			return
		}

		// Verify token and get user details
		fmt.Println("🔍 Fetching user details...")

		resp, err := client.VerifyToken()
		if err != nil {
			fmt.Printf("❌ Error verifying token: %v\n", err)
			fmt.Println("Your token may be expired. Please run 'pipeops auth login' to re-authenticate.")
			return
		}

		if !resp.Valid {
			fmt.Println("❌ Invalid token. Please run 'pipeops auth login' to re-authenticate.")
			return
		}

		// Display user information
		config := client.GetConfig()
		fmt.Println("✅ User Details:")
		if config.Username != "" {
			fmt.Printf("👤 Username: %s\n", config.Username)
		}
		if config.Email != "" {
			fmt.Printf("📧 Email: %s\n", config.Email)
		}
		if config.UserID != "" {
			fmt.Printf("🆔 User ID: %s\n", config.UserID)
		}
		fmt.Printf("🔗 API Base URL: %s\n", config.APIBaseURL)
		fmt.Println("🔐 Token: ✅ Valid")
	},
	Args: cobra.NoArgs, // This command does not accept arguments
}

func (k *authModel) me() {
	k.rootCmd.AddCommand(authCmd)
}
