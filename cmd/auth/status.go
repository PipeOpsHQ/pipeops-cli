package auth

import (
	"fmt"
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

		// Prepare status information
		status := map[string]interface{}{
			"authenticated": cfg.IsAuthenticated(),
			"api_endpoint":  cfg.OAuth.BaseURL,
			"client_id":     cfg.OAuth.ClientID,
			"scopes":        cfg.OAuth.Scopes,
		}

		if cfg.OAuth.AccessToken != "" {
			status["token_expires_at"] = cfg.OAuth.ExpiresAt.Format(time.RFC3339)
			status["token_expires_in_seconds"] = int(time.Until(cfg.OAuth.ExpiresAt).Seconds())
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(status)
		} else {
			fmt.Printf("🔐 AUTHENTICATION STATUS\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

			if cfg.IsAuthenticated() {
				utils.PrintSuccess("Authenticated", opts)
				fmt.Printf("├─ Expires: %s\n", cfg.OAuth.ExpiresAt.Format("2006-01-02 15:04:05 MST"))

				timeUntilExpiry := time.Until(cfg.OAuth.ExpiresAt)
				if timeUntilExpiry > 0 {
					fmt.Printf("├─ Expires in: %s\n", formatDuration(timeUntilExpiry))
				} else {
					fmt.Printf("├─ Expires in: ⚠️  Expired\n")
				}
			} else {
				utils.PrintError("Not authenticated", opts)
				fmt.Printf("├─ Action: Run 'pipeops auth login' to authenticate\n")
			}

			fmt.Printf("├─ API Endpoint: %s\n", cfg.OAuth.BaseURL)
			fmt.Printf("├─ Client ID: %s\n", cfg.OAuth.ClientID)
			fmt.Printf("└─ Scopes: %v\n", cfg.OAuth.Scopes)
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
