package auth

import (
	"context"
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "🔐 Login to your PipeOps account",
	Long: `🔐 Login to your PipeOps account using OAuth 2.0 authentication.

This command will:
1. Open your default browser to the PipeOps login page
2. After you authenticate, redirect back to the CLI
3. Store your authentication token securely
4. Verify the login by retrieving your user information

The authentication uses OAuth 2.0 with PKCE (Proof Key for Code Exchange) for
maximum security. Your credentials are never stored locally - only access tokens.

Examples:
  - Login to PipeOps:
    pipeops auth login

  - Login with custom client ID:
    PIPEOPS_CLIENT_ID=your-client-id pipeops auth login

  - Login with custom API URL:
    PIPEOPS_API_URL=https://api.example.com pipeops auth login`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		fmt.Println("🚀 Starting PipeOps CLI authentication...")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Failed to load configuration", opts)
			return
		}

		// Check if already authenticated
		if cfg.IsAuthenticated() {
			utils.PrintWarning("You are already logged in!", opts)

			// Get current user info to display
			authService := auth.NewOAuthService(cfg.OAuth)
			userInfo, err := authService.GetUserInfo(context.Background())
			if err == nil {
				fmt.Printf("Current user: %s (%s)\n", userInfo.Username, userInfo.Email)
			}

			fmt.Println("Use 'pipeops auth logout' to logout first if you want to login as a different user.")
			return
		}

		// Create OAuth service
		authService := auth.NewOAuthService(cfg.OAuth)

		// Perform login
		ctx := context.Background()
		if err := authService.Login(ctx); err != nil {
			utils.HandleError(err, "Login failed", opts)
			return
		}

		// Save updated configuration
		if err := config.Save(cfg); err != nil {
			utils.PrintWarning("Failed to save authentication: "+err.Error(), opts)
			fmt.Println("You may need to login again when you restart the CLI.")
		}

		// Verify login by getting user info
		userInfo, err := authService.GetUserInfo(ctx)
		if err != nil {
			utils.PrintWarning("Failed to verify login: "+err.Error(), opts)
			return
		}

		// Display success message
		utils.PrintSuccess("Successfully logged in!", opts)
		fmt.Printf("👤 Welcome, %s!\n", userInfo.FirstName)
		fmt.Printf("📧 Email: %s\n", userInfo.Email)
		fmt.Printf("🆔 Username: %s\n", userInfo.Username)

		if !opts.Quiet {
			fmt.Printf("\n💡 TIPS\n")
			fmt.Printf("├─ View your profile: pipeops auth me\n")
			fmt.Printf("├─ List your projects: pipeops list\n")
			fmt.Printf("├─ Check auth status: pipeops auth status\n")
			fmt.Printf("└─ Get help: pipeops --help\n")
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) login() {
	k.rootCmd.AddCommand(loginCmd)

	// Add flags
	loginCmd.Flags().Bool("json", false, "Output in JSON format")
}
