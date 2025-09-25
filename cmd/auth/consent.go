package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	Short: "View OAuth consent and permissions",
	Long: `View OAuth consent and permissions for your PipeOps CLI authentication.

Examples:
  pipeops auth consent
  pipeops auth consent --verbose
  pipeops auth consent --json`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		opts := utils.GetOutputOptions(cmd)

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Create auth service
		authService := auth.NewPKCEOAuthService(cfg)

		// Check if authenticated
		if !authService.IsAuthenticated() {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON(map[string]string{
					"error":   "not_authenticated",
					"message": "You must be authenticated to view consent information",
				})
			} else {
				fmt.Println("Not authenticated - run 'pipeops auth login'")
			}
			return
		}

		// Attempt to fetch consent info
		consentInfo, err := getConsentInfo(cfg, authService.GetAccessToken())
		if err != nil {
			// Check if this is an authentication method mismatch
			if isAuthenticationMismatch(err) {
				displayConsentUnavailableMessage(cfg, opts)
				return
			}
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

	// Set authorization header (try OAuth bearer token first)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PipeOps-CLI/1.0")

	// Debug information
	if cfg.Settings != nil && cfg.Settings.Debug {
		fmt.Printf("🔍 Debug: Making consent request to %s\n", req.URL.String())
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("consent request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for better error messages
	bodyBytes, readErr := io.ReadAll(resp.Body)
	var responseBody string
	if readErr == nil {
		responseBody = string(bodyBytes)
	}

	// Debug response information
	if cfg.Settings != nil && cfg.Settings.Debug {
		fmt.Printf("🔍 Debug: Consent response status: %d\n", resp.StatusCode)
		fmt.Printf("🔍 Debug: Consent response body: %s\n", responseBody)
	}

	// Check response status with detailed error messages
	if resp.StatusCode == http.StatusUnauthorized {
		if responseBody != "" {
			return nil, fmt.Errorf("consent access denied - server response: %s", responseBody)
		}
		return nil, fmt.Errorf("authentication expired - please run 'pipeops auth login'")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("consent endpoint not found - the API might not support this endpoint yet")
	}

	if resp.StatusCode != http.StatusOK {
		if responseBody != "" {
			return nil, fmt.Errorf("consent request failed with status %d: %s", resp.StatusCode, responseBody)
		}
		return nil, fmt.Errorf("consent request failed with status %d", resp.StatusCode)
	}

	// Parse response
	var consentInfo ConsentInfo
	if readErr != nil {
		return nil, fmt.Errorf("failed to read consent response: %w", readErr)
	}

	if err := json.Unmarshal(bodyBytes, &consentInfo); err != nil {
		return nil, fmt.Errorf("failed to parse consent response: %w (response: %s)", err, responseBody)
	}

	return &consentInfo, nil
}

// displayConsentInfo displays consent information in a formatted way
func displayConsentInfo(consent *ConsentInfo, verbose bool) {
	fmt.Printf("Application: %s\n", consent.ClientName)
	fmt.Printf("Client ID: %s\n", consent.ClientID)
	fmt.Printf("Granted: %s\n", consent.GrantedAt.Format("2006-01-02 15:04"))

	if !consent.ExpiresAt.IsZero() {
		timeUntilExpiry := time.Until(consent.ExpiresAt)
		if timeUntilExpiry > 0 {
			fmt.Printf("Expires: %s (%s remaining)\n", consent.ExpiresAt.Format("2006-01-02 15:04"), formatConsentDuration(timeUntilExpiry))
		} else {
			fmt.Printf("Expires: %s (expired)\n", consent.ExpiresAt.Format("2006-01-02 15:04"))
		}
	}

	if consent.Description != "" {
		fmt.Printf("Description: %s\n", consent.Description)
	}

	if len(consent.Scopes) > 0 {
		fmt.Printf("Scopes: %s\n", strings.Join(consent.Scopes, ", "))
	}

	if verbose && len(consent.Permissions) > 0 {
		fmt.Printf("Permissions: %s\n", strings.Join(consent.Permissions, ", "))
	}

	if verbose && len(consent.RedirectURIs) > 0 {
		fmt.Printf("Redirect URIs: %s\n", strings.Join(consent.RedirectURIs, ", "))
	}
}

// isAuthenticationMismatch checks if the error indicates an authentication method mismatch
func isAuthenticationMismatch(err error) bool {
	errorStr := err.Error()
	return strings.Contains(errorStr, "authentication failed") ||
		strings.Contains(errorStr, "consent access denied") ||
		strings.Contains(errorStr, "invalid token") ||
		strings.Contains(errorStr, "unauthorized") ||
		strings.Contains(errorStr, "403") ||
		strings.Contains(errorStr, "401")
}

// displayConsentUnavailableMessage shows when consent endpoint is not available via CLI
func displayConsentUnavailableMessage(cfg *config.Config, opts utils.OutputOptions) {
	if opts.Format == utils.OutputFormatJSON {
		utils.PrintJSON(map[string]interface{}{
			"error":         "consent_unavailable_cli",
			"message":       "Consent information requires web session authentication",
			"web_url":       cfg.OAuth.BaseURL + "/oauth/consent",
			"available_via": "web_interface",
		})
	} else {
		fmt.Println("Consent information requires web session authentication")
		fmt.Printf("Visit: %s/oauth/consent\n", cfg.OAuth.BaseURL)
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
