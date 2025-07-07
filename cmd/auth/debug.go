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
	Short: "🔍 Debug authentication issues",
	Long: `🔍 Debug and troubleshoot authentication issues with PipeOps API.

This command helps diagnose authentication problems by:
- Testing API endpoint connectivity
- Validating token format and expiration
- Testing OAuth endpoints with detailed logging
- Showing configuration details

Examples:
  - Run authentication debug:
    pipeops auth debug

  - Test specific endpoint:
    pipeops auth debug --test-userinfo

  - Show verbose debug information:
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

		fmt.Println()
		fmt.Println("🔍 PipeOps CLI Authentication Debug")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Show configuration details
		fmt.Println("📋 Configuration Details:")
		fmt.Printf("   API Endpoint: %s\n", cfg.OAuth.BaseURL)
		fmt.Printf("   Client ID: %s\n", cfg.OAuth.ClientID)
		fmt.Printf("   Scopes: %v\n", cfg.OAuth.Scopes)

		if cfg.OAuth.AccessToken != "" {
			fmt.Printf("   Token Present: ✅ Yes\n")
			fmt.Printf("   Token Length: %d characters\n", len(cfg.OAuth.AccessToken))
			fmt.Printf("   Token Preview: %s...\n", cfg.OAuth.AccessToken[:min(20, len(cfg.OAuth.AccessToken))])
			fmt.Printf("   Token Expires: %s\n", cfg.OAuth.ExpiresAt.Format("2006-01-02 15:04:05 MST"))

			// Check token expiry
			timeUntilExpiry := time.Until(cfg.OAuth.ExpiresAt)
			if timeUntilExpiry > 0 {
				fmt.Printf("   Token Status: ✅ Valid (%s remaining)\n", timeUntilExpiry.Truncate(time.Minute))
			} else {
				fmt.Printf("   Token Status: ❌ Expired (%s ago)\n", (-timeUntilExpiry).Truncate(time.Minute))
			}
		} else {
			fmt.Printf("   Token Present: ❌ No\n")
		}

		// Check authentication status
		fmt.Println()
		fmt.Println("🔐 Authentication Status:")
		if authService.IsAuthenticated() {
			fmt.Printf("   Local Check: ✅ Authenticated\n")
		} else {
			fmt.Printf("   Local Check: ❌ Not authenticated\n")
			fmt.Println("   💡 Run 'pipeops auth login' to authenticate")
			return
		}

		// Test endpoints if requested
		if testUserinfo || cmd.Flags().Changed("test-userinfo") {
			fmt.Println()
			fmt.Println("👤 Testing /oauth/userinfo endpoint:")
			testUserInfoEndpoint(cfg, authService, verbose)
		}

		if testConsent || cmd.Flags().Changed("test-consent") {
			fmt.Println()
			fmt.Println("🛡️ Testing /oauth/consent endpoint:")
			testConsentEndpoint(cfg, authService, verbose)
		}

		// If no specific tests requested, test both
		if !testUserinfo && !testConsent && !cmd.Flags().Changed("test-userinfo") && !cmd.Flags().Changed("test-consent") {
			fmt.Println()
			fmt.Println("👤 Testing /oauth/userinfo endpoint:")
			testUserInfoEndpoint(cfg, authService, verbose)

			fmt.Println()
			fmt.Println("🛡️ Testing /oauth/consent endpoint:")
			testConsentEndpoint(cfg, authService, verbose)
		}

		fmt.Println()
		fmt.Println("💡 Troubleshooting Tips:")
		fmt.Println("   • If endpoints return 404, the API might not support OAuth userinfo/consent yet")
		fmt.Println("   • If you get 401 errors, try: pipeops auth logout && pipeops auth login")
		fmt.Println("   • Check that your token has the required scopes: read:user")
		fmt.Println("   • Verify the API endpoint is correct in your configuration")
		fmt.Println()
	},
	Args: cobra.NoArgs,
}

// testUserInfoEndpoint tests the OAuth userinfo endpoint
func testUserInfoEndpoint(cfg *config.Config, authService *auth.PKCEOAuthService, verbose bool) {
	userInfoService := auth.NewUserInfoService(cfg)
	ctx := context.Background()

	userInfo, err := userInfoService.GetUserInfo(ctx, authService.GetAccessToken())
	if err != nil {
		fmt.Printf("   Result: ❌ Failed\n")
		fmt.Printf("   Error: %v\n", err)
		return
	}

	fmt.Printf("   Result: ✅ Success\n")
	fmt.Printf("   User ID: %s\n", userInfo.ID)
	fmt.Printf("   Username: %s\n", userInfo.Username)
	fmt.Printf("   Email: %s\n", userInfo.Email)

	if verbose {
		fmt.Printf("   Full Response: %+v\n", userInfo)
	}
}

// testConsentEndpoint tests the OAuth consent endpoint
func testConsentEndpoint(cfg *config.Config, authService *auth.PKCEOAuthService, verbose bool) {
	consentInfo, err := getConsentInfo(cfg, authService.GetAccessToken())
	if err != nil {
		fmt.Printf("   Result: ❌ Failed\n")
		fmt.Printf("   Error: %v\n", err)
		return
	}

	fmt.Printf("   Result: ✅ Success\n")
	fmt.Printf("   Client: %s\n", consentInfo.ClientName)
	fmt.Printf("   Scopes: %v\n", consentInfo.Scopes)

	if verbose {
		fmt.Printf("   Full Response: %+v\n", consentInfo)
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
