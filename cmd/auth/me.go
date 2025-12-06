package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:     "me",
	Aliases: []string{"whoami", "userinfo"},
	Short:   "Show current user information",
	Long: `Display information about the currently authenticated user.

Examples:
  pipeops auth me
  pipeops auth me --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Failed to load configuration", opts)
			return
		}

		// Create OAuth service
		authService := auth.NewPKCEOAuthService(cfg)

		// Check authentication
		if !authService.IsAuthenticated() {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON(map[string]interface{}{
					"authenticated": false,
					"error":         "not authenticated",
				})
			} else {
				fmt.Println("Not authenticated")
				fmt.Println()
				fmt.Println("You need to log in first:")
				fmt.Println("   pipeops auth login")
				fmt.Println()
				fmt.Println("Need help? Run: pipeops auth --help")
			}
			return
		}

		// Fetch user info from server
		userInfoService := auth.NewUserInfoService(cfg)
		ctx := context.Background()

		userInfo, err := userInfoService.GetUserInfo(ctx, authService.GetAccessToken())
		if err != nil {
			// If userinfo fails, fallback to token information
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON(map[string]interface{}{
					"authenticated": true,
					"error":         fmt.Sprintf("failed to fetch user info: %v", err),
					"fallback":      true,
					"token_info": map[string]interface{}{
						"client_id":    cfg.OAuth.ClientID,
						"api_endpoint": cfg.OAuth.BaseURL,
						"expires_at":   cfg.OAuth.ExpiresAt.Format(time.RFC3339),
						"scopes":       cfg.OAuth.Scopes,
					},
				})
			} else {
				fmt.Printf("Unable to fetch user details: %v\n", err)
				fmt.Println()
				fmt.Println("This might help:")
				fmt.Println("   • Check your internet connection")
				fmt.Println("   • Try: pipeops auth login")
				fmt.Println("   • For debugging: pipeops auth debug")
				showTokenInfo(cfg, authService, opts)
			}
			return
		}

		// Calculate token expiry info
		expiresAt := cfg.OAuth.ExpiresAt
		timeUntilExpiry := time.Until(expiresAt)

		// Format remaining time
		var expiryStatus string
		if timeUntilExpiry > 24*time.Hour {
			days := int(timeUntilExpiry.Hours() / 24)
			expiryStatus = fmt.Sprintf("%d days", days)
		} else if timeUntilExpiry > time.Hour {
			hours := int(timeUntilExpiry.Hours())
			expiryStatus = fmt.Sprintf("%d hours", hours)
		} else if timeUntilExpiry > 0 {
			minutes := int(timeUntilExpiry.Minutes())
			expiryStatus = fmt.Sprintf("%d minutes", minutes)
		} else {
			expiryStatus = "expired"
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			result := map[string]interface{}{
				"authenticated": true,
				"user": map[string]interface{}{
					"id":                  userInfo.ID,
					"uuid":                userInfo.UUID,
					"username":            userInfo.Username,
					"email":               userInfo.Email,
					"name":                userInfo.GetFullName(),
					"first_name":          userInfo.FirstName,
					"last_name":           userInfo.LastName,
					"avatar":              userInfo.Avatar,
					"email_verified":      userInfo.Verified,
					"subscription_active": userInfo.SubscriptionActive,
					"display_name":        userInfo.GetDisplayName(),
					"created_at":          userInfo.CreatedAt.Format(time.RFC3339),
					"updated_at":          userInfo.UpdatedAt.Format(time.RFC3339),
					"last_login_at":       userInfo.LastLoginAt.Format(time.RFC3339),
					"roles":               userInfo.Roles,
					"permissions":         userInfo.Permissions,
				},
				"authentication": map[string]interface{}{
					"client_id":    cfg.OAuth.ClientID,
					"api_endpoint": cfg.OAuth.BaseURL,
					"expires_at":   expiresAt.Format(time.RFC3339),
					"expires_in":   timeUntilExpiry.String(),
					"scopes":       cfg.OAuth.Scopes,
				},
			}
			utils.PrintJSON(result)
		} else {
			// Friendly, informative output
			fmt.Printf("Hello, %s!\n", userInfo.GetDisplayName())
			fmt.Printf("%s", userInfo.Email)
			if userInfo.Verified {
				fmt.Printf(" (verified)\n")
			} else {
				fmt.Printf(" (unverified)\n")
			}

			fmt.Printf("Session expires in %s\n", expiryStatus)

			// Show warnings and next steps based on token status
			if timeUntilExpiry <= 0 {
				fmt.Println()
				fmt.Println("Your session has expired")
				fmt.Println("Please login again: pipeops auth login")
			} else if timeUntilExpiry < 1*time.Hour {
				fmt.Println()
				fmt.Println("Your session expires soon")
				fmt.Println("Refresh it: pipeops auth login")
			} else {
				fmt.Println()
				fmt.Println("Ready to go! Try these commands:")
				fmt.Println("   pipeops project list    # View your projects")
				fmt.Println("   pipeops auth status     # Check auth details")
			}
		}
	},
	Args: cobra.NoArgs,
}

// showTokenInfo displays fallback token information when userinfo fails
func showTokenInfo(cfg *config.Config, authService *auth.PKCEOAuthService, opts utils.OutputOptions) {
	// Calculate time until expiration
	expiresAt := cfg.OAuth.ExpiresAt
	timeUntilExpiry := time.Until(expiresAt)

	var expiryStatus string
	if timeUntilExpiry > 24*time.Hour {
		days := int(timeUntilExpiry.Hours() / 24)
		expiryStatus = fmt.Sprintf("%d days", days)
	} else if timeUntilExpiry > time.Hour {
		hours := int(timeUntilExpiry.Hours())
		expiryStatus = fmt.Sprintf("%d hours", hours)
	} else if timeUntilExpiry > 0 {
		minutes := int(timeUntilExpiry.Minutes())
		expiryStatus = fmt.Sprintf("%d minutes", minutes)
	} else {
		expiryStatus = "expired"
	}

	fmt.Printf("Token expires in %s\n", expiryStatus)
}

func (k *authModel) me() {
	k.rootCmd.AddCommand(meCmd)
}
