package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "📊 Show authentication status",
	Long: `📊 Display current authentication status including token expiration.

This command shows:
- Whether you're currently logged in
- Token expiration time
- Configured API endpoint
- Available scopes

Examples:
  - Show authentication status:
    pipeops auth status

  - Show status in JSON format:
    pipeops auth status --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Failed to load configuration", opts)
			return
		}

		// Calculate time until expiration
		var timeUntilExpiry time.Duration
		var expiryStatus string
		var expiryColor string
		var statusIcon string

		if cfg.OAuth.AccessToken != "" {
			timeUntilExpiry = time.Until(cfg.OAuth.ExpiresAt)

			if timeUntilExpiry > 24*time.Hour {
				days := int(timeUntilExpiry.Hours() / 24)
				expiryStatus = fmt.Sprintf("%d days", days)
				expiryColor = "🟢"
				statusIcon = "✅"
			} else if timeUntilExpiry > time.Hour {
				hours := int(timeUntilExpiry.Hours())
				expiryStatus = fmt.Sprintf("%d hours", hours)
				expiryColor = "🟡"
				statusIcon = "⚠️"
			} else if timeUntilExpiry > 0 {
				minutes := int(timeUntilExpiry.Minutes())
				expiryStatus = fmt.Sprintf("%d minutes", minutes)
				expiryColor = "🟠"
				statusIcon = "⏰"
			} else {
				expiryStatus = "Expired"
				expiryColor = "🔴"
				statusIcon = "❌"
			}
		}

		// Prepare status information
		status := map[string]interface{}{
			"authenticated": cfg.IsAuthenticated(),
			"api_endpoint":  cfg.OAuth.BaseURL,
			"client_id":     cfg.OAuth.ClientID,
			"scopes":        cfg.OAuth.Scopes,
		}

		if cfg.OAuth.AccessToken != "" {
			status["token_expires_at"] = cfg.OAuth.ExpiresAt.Format(time.RFC3339)
			status["token_expires_in_seconds"] = int(timeUntilExpiry.Seconds())
			status["expires_in_human"] = expiryStatus
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(status)
		} else {
			fmt.Println()
			fmt.Printf("🔐 PipeOps Authentication Status\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

			if cfg.IsAuthenticated() {
				fmt.Printf("%s Status: %s Authenticated\n", statusIcon, expiryColor)
				fmt.Printf("🔑 Token: %s...\n", cfg.OAuth.AccessToken[:20])
				fmt.Printf("⏰ Expires: %s\n", cfg.OAuth.ExpiresAt.Format("2006-01-02 15:04:05 MST"))
				fmt.Printf("⌛ Time remaining: %s\n", expiryStatus)
				fmt.Printf("🌐 API Endpoint: %s\n", cfg.OAuth.BaseURL)
				fmt.Printf("🏷️  Client ID: %s\n", cfg.OAuth.ClientID)
				fmt.Printf("🎯 Scopes: %s\n", formatScopes(cfg.OAuth.Scopes))
				fmt.Println()

				// Show warnings and tips
				if timeUntilExpiry < 24*time.Hour && timeUntilExpiry > 0 {
					fmt.Println("⚠️  WARNING: Your token expires soon!")
					fmt.Println("   Run 'pipeops auth login' to refresh your session")
					fmt.Println()
				} else if timeUntilExpiry <= 0 {
					fmt.Println("❌ EXPIRED: Your token has expired!")
					fmt.Println("   Run 'pipeops auth login' to authenticate again")
					fmt.Println()
				}

				// Show available actions
				fmt.Println("🚀 Available Actions:")
				fmt.Println("   📋 pipeops project list      - List your projects")
				fmt.Println("   👤 pipeops auth me           - Show user info")
				fmt.Println("   🔄 pipeops auth login        - Refresh authentication")
				fmt.Println("   🚪 pipeops auth logout       - Sign out")

			} else {
				fmt.Printf("❌ Status: Not authenticated\n")
				fmt.Printf("🌐 API Endpoint: %s\n", cfg.OAuth.BaseURL)
				fmt.Printf("🏷️  Client ID: %s\n", cfg.OAuth.ClientID)
				fmt.Printf("🎯 Scopes: %s\n", formatScopes(cfg.OAuth.Scopes))
				fmt.Println()

				fmt.Println("🚀 Get Started:")
				fmt.Println("   pipeops auth login")
				fmt.Println()

				fmt.Println("💡 Need Help?")
				fmt.Println("   pipeops auth --help")
			}

			fmt.Println()
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) status() {
	k.rootCmd.AddCommand(statusCmd)
}

// formatDuration formats a duration for display
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "expired"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

// formatScopes formats scopes for display
func formatScopes(scopes []string) string {
	if len(scopes) == 0 {
		return "none"
	}

	// Add icons for common scopes
	var formatted []string
	for _, scope := range scopes {
		switch scope {
		case "read:user":
			formatted = append(formatted, "👤 "+scope)
		case "read:projects":
			formatted = append(formatted, "📋 "+scope)
		case "write:projects":
			formatted = append(formatted, "✏️ "+scope)
		default:
			formatted = append(formatted, "🔧 "+scope)
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(formatted, ", "))
}
