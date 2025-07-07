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

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug authentication issues",
	Long: `Debug and troubleshoot authentication issues with PipeOps API.

Examples:
  pipeops auth debug
  pipeops auth debug --test-userinfo
  pipeops auth debug --verbose`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		verbose, _ := cmd.Flags().GetBool("verbose")
		testUserinfo, _ := cmd.Flags().GetBool("test-userinfo")
		testConsent, _ := cmd.Flags().GetBool("test-consent")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Failed to load configuration", opts)
			return
		}

		// Enable debug mode for this session
		if cfg.Settings == nil {
			cfg.Settings = &config.Settings{}
		}
		cfg.Settings.Debug = true

		// Create OAuth service
		authService := auth.NewPKCEOAuthService(cfg)

		// Show configuration details
		fmt.Printf("API Endpoint: %s\n", cfg.OAuth.BaseURL)
		fmt.Printf("Client ID: %s\n", cfg.OAuth.ClientID)
		fmt.Printf("Scopes: %v\n", cfg.OAuth.Scopes)

		if cfg.OAuth.AccessToken != "" {
			fmt.Printf("Token: %s... (%d chars)\n", cfg.OAuth.AccessToken[:min(20, len(cfg.OAuth.AccessToken))], len(cfg.OAuth.AccessToken))
			fmt.Printf("Expires: %s\n", cfg.OAuth.ExpiresAt.Format("2006-01-02 15:04:05 MST"))

			// Check token expiry
			timeUntilExpiry := time.Until(cfg.OAuth.ExpiresAt)
			if timeUntilExpiry > 0 {
				fmt.Printf("Status: Valid (%s remaining)\n", timeUntilExpiry.Truncate(time.Minute))
			} else {
				fmt.Printf("Status: Expired (%s ago)\n", (-timeUntilExpiry).Truncate(time.Minute))
			}
		} else {
			fmt.Println("Token: Not present")
		}

		// Check authentication status
		if authService.IsAuthenticated() {
			fmt.Println("Authentication: Valid")
		} else {
			fmt.Println("Authentication: Invalid - run 'pipeops auth login'")
			return
		}

		// Test endpoints if requested
		if testUserinfo || cmd.Flags().Changed("test-userinfo") {
			fmt.Println("\nTesting /oauth/userinfo:")
			testUserInfoEndpoint(cfg, authService, verbose)
		}

		if testConsent || cmd.Flags().Changed("test-consent") {
			fmt.Println("\nTesting /oauth/consent:")
			testConsentEndpoint(cfg, authService, verbose)
		}

		// If no specific tests requested, test both
		if !testUserinfo && !testConsent && !cmd.Flags().Changed("test-userinfo") && !cmd.Flags().Changed("test-consent") {
			fmt.Println("\nTesting /oauth/userinfo:")
			testUserInfoEndpoint(cfg, authService, verbose)

			fmt.Println("\nTesting /oauth/consent:")
			testConsentEndpoint(cfg, authService, verbose)
		}
	},
	Args: cobra.NoArgs,
}

// testUserInfoEndpoint tests the OAuth userinfo endpoint
func testUserInfoEndpoint(cfg *config.Config, authService *auth.PKCEOAuthService, verbose bool) {
	userInfoService := auth.NewUserInfoService(cfg)
	ctx := context.Background()

	userInfo, err := userInfoService.GetUserInfo(ctx, authService.GetAccessToken())
	if err != nil {
		fmt.Printf("Result: Failed - %v\n", err)
		return
	}

	fmt.Printf("Result: Success\n")
	fmt.Printf("User: %s (%s)\n", userInfo.GetIDString(), userInfo.Email)

	if verbose {
		fmt.Printf("Full Response: %+v\n", userInfo)
	}
}

// testConsentEndpoint tests the OAuth consent endpoint
func testConsentEndpoint(cfg *config.Config, authService *auth.PKCEOAuthService, verbose bool) {
	consentInfo, err := getConsentInfo(cfg, authService.GetAccessToken())
	if err != nil {
		// Check if this is likely an authentication method mismatch
		if isAuthenticationMismatch(err) {
			fmt.Printf("Result: Authentication method mismatch (expected)\n")
			fmt.Printf("Reason: Consent endpoint requires JWT session tokens\n")
			if verbose {
				fmt.Printf("Error: %v\n", err)
			}
			return
		}

		fmt.Printf("Result: Failed - %v\n", err)
		return
	}

	fmt.Printf("Result: Success\n")
	fmt.Printf("Client: %s\n", consentInfo.ClientName)
	fmt.Printf("Scopes: %v\n", consentInfo.Scopes)

	if verbose {
		fmt.Printf("Full Response: %+v\n", consentInfo)
	}
}

// Helper function to get minimum of two integers (avoid duplication)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (k *authModel) debug() {
	debugCmd.Flags().BoolP("verbose", "v", false, "Show detailed debug information")
	debugCmd.Flags().Bool("test-userinfo", false, "Test only the userinfo endpoint")
	debugCmd.Flags().Bool("test-consent", false, "Test only the consent endpoint")
	k.rootCmd.AddCommand(debugCmd)
}
