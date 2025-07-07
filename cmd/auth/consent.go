package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// ConsentInfo represents OAuth consent information
type ConsentInfo struct {
	ClientID     string    `json:"client_id"`
	ClientName   string    `json:"client_name"`
	Scopes       []string  `json:"scopes"`
	GrantedAt    time.Time `json:"granted_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Permissions  []string  `json:"permissions"`
	Description  string    `json:"description"`
	RedirectURIs []string  `json:"redirect_uris"`
}

// consentCmd represents the consent command
var consentCmd = &cobra.Command{
	Use:   "consent",
	Short: "🛡️ Manage OAuth consent and permissions",
	Long: `🛡️ Display and manage OAuth consent and permissions for PipeOps CLI.

This command shows:
- Current OAuth permissions granted to the CLI
- Scope details and descriptions
- Consent grant and expiration dates
- Available actions for managing consent

Examples:
  - Show current consent information:
    pipeops auth consent

  - Show consent in JSON format:
    pipeops auth consent --json

  - View consent details:
    pipeops auth consent --verbose`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		verbose, _ := cmd.Flags().GetBool("verbose")

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
				fmt.Println("❌ You must be authenticated to view consent information")
				fmt.Println()
				fmt.Println("🚀 Get started:")
				fmt.Println("   pipeops auth login")
				fmt.Println()
			}
			return
		}

		// Fetch consent info
		consentInfo, err := getConsentInfo(cfg, authService.GetAccessToken())
		if err != nil {
			utils.HandleError(err, "Failed to fetch consent information", opts)
			return
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(consentInfo)
		} else {
			displayConsentInfo(consentInfo, verbose)
		}
	},
	Args: cobra.NoArgs,
}

// getConsentInfo fetches consent information from the OAuth consent endpoint
func getConsentInfo(cfg *config.Config, accessToken string) (*ConsentInfo, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	// Create request to consent endpoint
	req, err := http.NewRequest("GET", cfg.OAuth.BaseURL+"/oauth/consent", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consent request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PipeOps-CLI/1.0")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("consent request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication expired - please run 'pipeops auth login'")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("consent request failed with status %d", resp.StatusCode)
	}

	// Parse response
	var consentInfo ConsentInfo
	if err := json.NewDecoder(resp.Body).Decode(&consentInfo); err != nil {
		return nil, fmt.Errorf("failed to parse consent response: %w", err)
	}

	return &consentInfo, nil
}

// displayConsentInfo displays consent information in a formatted way
func displayConsentInfo(consent *ConsentInfo, verbose bool) {
	fmt.Println()
	fmt.Println("🛡️ OAuth Consent & Permissions")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🏷️  Application: %s\n", consent.ClientName)
	fmt.Printf("🆔 Client ID: %s\n", consent.ClientID)
	fmt.Printf("📅 Granted: %s\n", consent.GrantedAt.Format("January 2, 2006 15:04 MST"))

	if !consent.ExpiresAt.IsZero() {
		timeUntilExpiry := time.Until(consent.ExpiresAt)
		if timeUntilExpiry > 0 {
			fmt.Printf("⏰ Expires: %s (%s remaining)\n", consent.ExpiresAt.Format("January 2, 2006 15:04 MST"), formatConsentDuration(timeUntilExpiry))
		} else {
			fmt.Printf("⏰ Expires: %s (⚠️ Expired)\n", consent.ExpiresAt.Format("January 2, 2006 15:04 MST"))
		}
	}

	if consent.Description != "" {
		fmt.Printf("📝 Description: %s\n", consent.Description)
	}

	fmt.Println()
	fmt.Println("🎯 Granted Scopes:")
	for _, scope := range consent.Scopes {
		scopeIcon := getScopeIcon(scope)
		scopeDesc := getScopeDescription(scope)
		fmt.Printf("   %s %s", scopeIcon, scope)
		if verbose && scopeDesc != "" {
			fmt.Printf(" - %s", scopeDesc)
		}
		fmt.Println()
	}

	if len(consent.Permissions) > 0 {
		fmt.Println()
		fmt.Println("🔑 Specific Permissions:")
		for _, perm := range consent.Permissions {
			fmt.Printf("   ✅ %s\n", perm)
		}
	}

	if verbose && len(consent.RedirectURIs) > 0 {
		fmt.Println()
		fmt.Println("🔄 Redirect URIs:")
		for _, uri := range consent.RedirectURIs {
			fmt.Printf("   🌐 %s\n", uri)
		}
	}

	fmt.Println()
	fmt.Println("🚀 Available Actions:")
	fmt.Println("   📋 pipeops project list      - Use your permissions")
	fmt.Println("   👤 pipeops auth me           - View user profile")
	fmt.Println("   🔄 pipeops auth login        - Refresh authentication")
	fmt.Println("   🚪 pipeops auth logout       - Revoke access")
	fmt.Println()

	fmt.Println("💡 TIP: To revoke consent, run 'pipeops auth logout' and re-authenticate")
	fmt.Println()
}

// getScopeIcon returns an appropriate icon for a scope
func getScopeIcon(scope string) string {
	switch scope {
	case "read:user", "user:read":
		return "👤"
	case "read:projects", "projects:read":
		return "📋"
	case "write:projects", "projects:write":
		return "✏️"
	case "read:deployments", "deployments:read":
		return "🚀"
	case "write:deployments", "deployments:write":
		return "🔧"
	case "read:servers", "servers:read":
		return "🖥️"
	case "write:servers", "servers:write":
		return "⚙️"
	default:
		return "🔹"
	}
}

// getScopeDescription returns a description for a scope
func getScopeDescription(scope string) string {
	switch scope {
	case "read:user", "user:read":
		return "View your profile information"
	case "read:projects", "projects:read":
		return "View your projects"
	case "write:projects", "projects:write":
		return "Create and modify projects"
	case "read:deployments", "deployments:read":
		return "View deployment information"
	case "write:deployments", "deployments:write":
		return "Create and manage deployments"
	case "read:servers", "servers:read":
		return "View server information"
	case "write:servers", "servers:write":
		return "Manage servers"
	default:
		return ""
	}
}

// formatConsentDuration formats a duration for display in consent context
func formatConsentDuration(d time.Duration) string {
	if d < 0 {
		return "expired"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

func (k *authModel) consent() {
	consentCmd.Flags().BoolP("verbose", "v", false, "Show detailed consent information")
	k.rootCmd.AddCommand(consentCmd)
}
