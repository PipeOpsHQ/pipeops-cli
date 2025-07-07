package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "🔑 Login to PipeOps",
	Long: `🔑 Login to PipeOps using OAuth2 authentication.

This command will:
1. Open your default browser
2. Redirect you to PipeOps authentication
3. Store your authentication tokens securely

Examples:
  pipeops auth login`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ Failed to load configuration: %v\n", err)
			return
		}

		// Create PKCE OAuth service
		oauthService := auth.NewPKCEOAuthService(cfg)

		// Check if already authenticated
		if oauthService.IsAuthenticated() {
			fmt.Println()
			fmt.Println("🔐 Already Authenticated")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println("✅ You are already authenticated with PipeOps!")
			fmt.Printf("🔑 Your session expires: %s\n", cfg.OAuth.ExpiresAt.Format("2006-01-02 15:04:05 MST"))
			fmt.Println()
			fmt.Println("💡 Available commands:")
			fmt.Println("   🔍 pipeops auth status    - Check authentication status")
			fmt.Println("   👤 pipeops auth me        - Show user information")
			fmt.Println("   📋 pipeops project list   - List your projects")
			fmt.Println("   🚪 pipeops auth logout    - Sign out")
			fmt.Println()
			return
		}

		// Perform authentication
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := oauthService.Login(ctx); err != nil {
			fmt.Printf("\n❌ Authentication failed: %v\n", err)
			fmt.Println()
			fmt.Println("🔧 Troubleshooting:")
			fmt.Println("   • Make sure your browser is working")
			fmt.Println("   • Check your internet connection")
			fmt.Println("   • Try running the command again")
			fmt.Println("   • Contact support if the problem persists")
			return
		}

		// Save updated configuration
		fmt.Print("💾 Saving credentials... ")
		if err := config.Save(cfg); err != nil {
			fmt.Println("❌ Failed")
			fmt.Printf("⚠️  Failed to save authentication tokens: %v\n", err)
			fmt.Println("You may need to authenticate again on next use.")
			return
		}
		fmt.Println("✅ Saved")
	},
}

func (k *authModel) login() {
	k.rootCmd.AddCommand(loginCmd)

	// Add flags
	loginCmd.Flags().Bool("json", false, "Output in JSON format")
}
