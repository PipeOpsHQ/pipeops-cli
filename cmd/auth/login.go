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

		// Check if already authenticated
		if oauthService.IsAuthenticated() {
			fmt.Println("Already authenticated")
			return
		}

		// Perform authentication
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := oauthService.Login(ctx); err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		// Save updated configuration
		if err := config.Save(cfg); err != nil {
			fmt.Printf("Failed to save credentials: %v\n", err)
			fmt.Println("You may need to authenticate again on next use.")
			return
		}
	},
}

func (k *authModel) login() {
	k.rootCmd.AddCommand(loginCmd)

	// Add flags
	loginCmd.Flags().Bool("json", false, "Output in JSON format")
}
