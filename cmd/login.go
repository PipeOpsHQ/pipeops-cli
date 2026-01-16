package cmd

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
  pipeops login`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			return
		}

		// Create PKCE OAuth service
		oauthService := auth.NewPKCEOAuthService(cfg)

		// Override config with flags if provided
		if clientID, _ := cmd.Flags().GetString("client-id"); clientID != "" {
			cfg.OAuth.ClientID = clientID
		}
		if authURL, _ := cmd.Flags().GetString("auth-url"); authURL != "" {
			cfg.OAuth.DashboardURL = authURL
		}
		if tokenURL, _ := cmd.Flags().GetString("token-url"); tokenURL != "" {
			cfg.OAuth.BaseURL = tokenURL
		}

		// Check if already authenticated (local check)
		if oauthService.IsAuthenticated() {
			// Validate with server to ensure token is still valid
			userInfoService := auth.NewUserInfoService(cfg)
			ctx := context.Background()

			if _, err := userInfoService.GetUserInfo(ctx, oauthService.GetAccessToken()); err == nil {
				fmt.Println("You're already authenticated!")
				fmt.Println("Ready to use PipeOps. Try: pipeops project list")
				return
			} else {
				// Token is invalid on server, clear it and proceed with login
				fmt.Println("Your session has expired or been revoked")
				fmt.Println("Starting fresh authentication...")
				cfg.ClearAuth()
				config.Save(cfg)
			}
		}

		// Perform authentication
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := oauthService.Login(ctx); err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			fmt.Println()
			fmt.Println("Troubleshooting tips:")
			fmt.Println("   • Check your internet connection")
			fmt.Println("   • Make sure you complete the login in your browser")
			fmt.Println("   • Try again: pipeops login")
			return
		}

		// Save updated configuration
		if err := config.Save(cfg); err != nil {
			fmt.Printf("Failed to save credentials: %v\n", err)
			fmt.Println("   You may need to authenticate again next time")
			return
		}

		// Show helpful next steps
		fmt.Println()
		fmt.Println("What's next? Try these commands:")
		fmt.Println("   pipeops project list     # See your projects")
		fmt.Println("   pipeops me          # View your profile")
		fmt.Println("   pipeops --help           # Explore all commands")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().Bool("json", false, "Output in JSON format")
	loginCmd.Flags().String("client-id", "", "OAuth2 client ID")
	loginCmd.Flags().String("auth-url", "", "OAuth2 authorization URL")
	loginCmd.Flags().String("token-url", "", "OAuth2 token URL")
}
