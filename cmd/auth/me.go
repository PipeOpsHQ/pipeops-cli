package auth

import (
	"context"
	"fmt"
	"strings"
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
	Short:   "👤 Show current user information",
	Long: `👤 Display detailed information about the currently authenticated user.

This command fetches your profile information from PipeOps including:
- User ID, username, and display name
- Email address and verification status
- Account creation and last login dates
- User roles and permissions
- Authentication token details

Examples:
  - Show user info:
    pipeops auth me

  - Show user info in JSON format:
    pipeops auth me --json

  - Alternative commands:
    pipeops auth whoami
    pipeops auth userinfo`,
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
				fmt.Println()
				fmt.Println("🔒 Not Authenticated")
				fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				fmt.Println("❌ You are not authenticated with PipeOps")
				fmt.Println()
				fmt.Println("🚀 Get started:")
				fmt.Println("   pipeops auth login")
				fmt.Println()
				fmt.Println("💡 Need help?")
				fmt.Println("   pipeops auth --help")
				fmt.Println()
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
				fmt.Println()
				fmt.Println("⚠️  Failed to fetch user details from server")
				fmt.Printf("Error: %v\n", err)
				fmt.Println()
				fmt.Println("Falling back to local token information...")
				showTokenInfo(cfg, authService, opts)
			}
			return
		}

		// Calculate token expiry info
		expiresAt := cfg.OAuth.ExpiresAt
		timeUntilExpiry := time.Until(expiresAt)

		// Format remaining time
		var expiryStatus string
		var expiryColor string
		if timeUntilExpiry > 24*time.Hour {
			days := int(timeUntilExpiry.Hours() / 24)
			expiryStatus = fmt.Sprintf("%d days", days)
			expiryColor = "🟢"
		} else if timeUntilExpiry > time.Hour {
			hours := int(timeUntilExpiry.Hours())
			expiryStatus = fmt.Sprintf("%d hours", hours)
			expiryColor = "🟡"
		} else if timeUntilExpiry > 0 {
			minutes := int(timeUntilExpiry.Minutes())
			expiryStatus = fmt.Sprintf("%d minutes", minutes)
			expiryColor = "🟠"
		} else {
			expiryStatus = "Expired"
			expiryColor = "🔴"
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			result := map[string]interface{}{
				"authenticated": true,
				"user": map[string]interface{}{
					"id":            userInfo.ID,
					"username":      userInfo.Username,
					"name":          userInfo.Name,
					"first_name":    userInfo.FirstName,
					"last_name":     userInfo.LastName,
					"email":         userInfo.Email,
					"verified":      userInfo.Verified,
					"avatar":        userInfo.Avatar,
					"created_at":    userInfo.CreatedAt.Format(time.RFC3339),
					"updated_at":    userInfo.UpdatedAt.Format(time.RFC3339),
					"last_login_at": userInfo.LastLoginAt.Format(time.RFC3339),
					"roles":         userInfo.Roles,
					"permissions":   userInfo.Permissions,
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
			fmt.Println()
			fmt.Println("👤 User Profile")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Print(userInfo.FormatUserInfo())
			fmt.Println()

			fmt.Println("🔐 Authentication Details")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("✅ Status: %s Authenticated\n", expiryColor)
			fmt.Printf("⏰ Token expires: %s (%s remaining)\n", expiresAt.Format("2006-01-02 15:04:05 MST"), expiryStatus)
			fmt.Printf("🌐 API Endpoint: %s\n", cfg.OAuth.BaseURL)
			fmt.Printf("🏷️  Client ID: %s\n", cfg.OAuth.ClientID)
			fmt.Printf("🎯 Scopes: %s\n", strings.Join(cfg.OAuth.Scopes, ", "))
			fmt.Println()

			// Show quick actions
			fmt.Println("🚀 Quick Actions:")
			fmt.Println("   📋 pipeops project list      - List your projects")
			fmt.Println("   📊 pipeops auth status        - Full authentication status")
			fmt.Println("   🔄 pipeops auth login         - Refresh authentication")
			fmt.Println("   🚪 pipeops auth logout        - Sign out")
			fmt.Println()

			// Show tips based on expiry
			if timeUntilExpiry < 24*time.Hour && timeUntilExpiry > 0 {
				fmt.Println("💡 TIP: Your token expires soon. Run 'pipeops auth login' to refresh it.")
			} else if timeUntilExpiry <= 0 {
				fmt.Println("❌ WARNING: Your token has expired. Run 'pipeops auth login' to authenticate again.")
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

	// Format remaining time
	var expiryStatus string
	var expiryColor string
	if timeUntilExpiry > 24*time.Hour {
		days := int(timeUntilExpiry.Hours() / 24)
		expiryStatus = fmt.Sprintf("%d days", days)
		expiryColor = "🟢"
	} else if timeUntilExpiry > time.Hour {
		hours := int(timeUntilExpiry.Hours())
		expiryStatus = fmt.Sprintf("%d hours", hours)
		expiryColor = "🟡"
	} else if timeUntilExpiry > 0 {
		minutes := int(timeUntilExpiry.Minutes())
		expiryStatus = fmt.Sprintf("%d minutes", minutes)
		expiryColor = "🟠"
	} else {
		expiryStatus = "Expired"
		expiryColor = "🔴"
	}

	fmt.Println()
	fmt.Println("🔐 Authentication Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("✅ Status: %s Authenticated\n", expiryColor)
	fmt.Printf("🔑 Token: %s...\n", authService.GetAccessToken()[:20])
	fmt.Printf("⏰ Expires: %s (%s remaining)\n", expiresAt.Format("2006-01-02 15:04:05 MST"), expiryStatus)
	fmt.Printf("🌐 API Endpoint: %s\n", cfg.OAuth.BaseURL)
	fmt.Printf("🏷️  Client ID: %s\n", cfg.OAuth.ClientID)
	fmt.Printf("🎯 Scopes: %s\n", strings.Join(cfg.OAuth.Scopes, ", "))
	fmt.Println()

	// Show quick actions
	fmt.Println("🚀 Quick Actions:")
	fmt.Println("   📋 pipeops project list      - List your projects")
	fmt.Println("   📊 pipeops auth status        - Full authentication status")
	fmt.Println("   🔄 pipeops auth login         - Refresh authentication")
	fmt.Println("   🚪 pipeops auth logout        - Sign out")
	fmt.Println()

	// Show tips based on expiry
	if timeUntilExpiry < 24*time.Hour {
		fmt.Println("💡 TIP: Your token expires soon. Run 'pipeops auth login' to refresh it.")
	}
}

func (k *authModel) me() {
	k.rootCmd.AddCommand(meCmd)
}

// formatTime formats a time string for display
func formatTime(timeStr string) string {
	if timeStr == "" {
		return "N/A"
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}

	return t.Format("2006-01-02 15:04:05 MST")
}
