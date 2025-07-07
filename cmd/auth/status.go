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
	Short: "Show authentication status",
	Long: `Display current authentication status.

Examples:
  pipeops auth status
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

		if cfg.OAuth.AccessToken != "" {
			timeUntilExpiry = time.Until(cfg.OAuth.ExpiresAt)

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
			if cfg.IsAuthenticated() {
				if timeUntilExpiry > 24*time.Hour {
					fmt.Printf("✅ Authenticated (%s remaining)\n", expiryStatus)
					fmt.Println()
					fmt.Println("🚀 You're all set! Try these commands:")
					fmt.Println("   pipeops project list    # View your projects")
					fmt.Println("   pipeops auth me         # See your profile")
				} else if timeUntilExpiry > time.Hour {
					fmt.Printf("✅ Authenticated (%s remaining)\n", expiryStatus)
					fmt.Println("💡 Your session is active and ready to use")
				} else if timeUntilExpiry > 0 {
					fmt.Printf("⚠️  Authenticated (%s remaining)\n", expiryStatus)
					fmt.Println("🔄 Consider refreshing soon: pipeops auth login")
				} else {
					fmt.Println("❌ Token expired")
					fmt.Println("🔑 Please login again: pipeops auth login")
				}
			} else {
				fmt.Println("❌ Not authenticated")
				fmt.Println()
				fmt.Println("👋 Welcome to PipeOps! Let's get you started:")
				fmt.Println("   pipeops auth login      # Authenticate with your account")
				fmt.Println("   pipeops --help          # Explore available commands")
			}
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) status() {
	k.rootCmd.AddCommand(statusCmd)
}

// formatScopes formats scopes for display
func formatScopes(scopes []string) string {
	if len(scopes) == 0 {
		return "none"
	}
	return strings.Join(scopes, ", ")
}
