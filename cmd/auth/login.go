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
	Short: "Login to PipeOps",
	Long: `Login to PipeOps using OAuth2 authentication.

Examples:
  pipeops auth login`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			return
		}

		// Create PKCE OAuth service
		oauthService := auth.NewPKCEOAuthService(cfg)

		// Check if already authenticated (local check)
		if oauthService.IsAuthenticated() {
			// Validate with server to ensure token is still valid
			userInfoService := auth.NewUserInfoService(cfg)
			ctx := context.Background()

			if _, err := userInfoService.GetUserInfo(ctx, oauthService.GetAccessToken()); err == nil {
				fmt.Println("✅ You're already authenticated!")
				fmt.Println("🚀 Ready to use PipeOps. Try: pipeops project list")
				return
			} else {
				// Token is invalid on server, clear it and proceed with login
				fmt.Println("⚠️  Your session has expired or been revoked")
				fmt.Println("🔄 Starting fresh authentication...")
				cfg.ClearAuth()
				config.Save(cfg)
			}
		}

		// Perform authentication
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := oauthService.Login(ctx); err != nil {
			fmt.Printf("❌ Authentication failed: %v\n", err)
			fmt.Println()
			fmt.Println("🔧 Troubleshooting tips:")
			fmt.Println("   • Check your internet connection")
			fmt.Println("   • Make sure you complete the login in your browser")
			fmt.Println("   • Try again: pipeops auth login")
			return
		}

		// Save updated configuration
		if err := config.Save(cfg); err != nil {
			fmt.Printf("⚠️  Failed to save credentials: %v\n", err)
			fmt.Println("   You may need to authenticate again next time")
			return
		}

		// Show helpful next steps
		fmt.Println()
		fmt.Println("🎯 What's next? Try these commands:")
		fmt.Println("   pipeops project list     # See your projects")
		fmt.Println("   pipeops auth me          # View your profile")
		fmt.Println("   pipeops --help           # Explore all commands")
	},
}

func (k *authModel) login() {
	k.rootCmd.AddCommand(loginCmd)

	// Add flags
	loginCmd.Flags().Bool("json", false, "Output in JSON format")
}
