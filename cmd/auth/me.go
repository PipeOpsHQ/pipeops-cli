package auth

import (
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
	Aliases: []string{"whoami"},
	Short:   "👤 Show current user information",
	Long: `👤 Display information about the currently authenticated user.

This command shows:
- User ID and username
- Email address and name
- Account creation date
- Authentication status

Examples:
  - Show user info:
    pipeops auth me

  - Show user info in JSON format:
    pipeops auth me --json

  - Alternative command:
    pipeops auth whoami`,
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

		// Get user info
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

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(map[string]interface{}{
				"authenticated": true,
				"client_id":     cfg.OAuth.ClientID,
				"api_endpoint":  cfg.OAuth.BaseURL,
				"expires_at":    expiresAt.Format(time.RFC3339),
				"expires_in":    timeUntilExpiry.String(),
				"scopes":        cfg.OAuth.Scopes,
				"token_preview": authService.GetAccessToken()[:20] + "...",
			})
		} else {
			fmt.Println()
			fmt.Println("👤 User Authentication Status")
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
			fmt.Println("   🔍 pipeops auth status        - Full authentication details")
			fmt.Println("   🔄 pipeops auth login         - Refresh authentication")
			fmt.Println("   🚪 pipeops auth logout        - Sign out")
			fmt.Println()

			// Show tips based on expiry
			if timeUntilExpiry < 24*time.Hour {
				fmt.Println("💡 TIP: Your token expires soon. Run 'pipeops auth login' to refresh it.")
			}
		}
	},
	Args: cobra.NoArgs,
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
