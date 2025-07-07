package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
)

// PKCEOAuthService handles OAuth2 authentication with PKCE
type PKCEOAuthService struct {
	config *config.Config
	client *http.Client
}

// NewPKCEOAuthService creates a new PKCE OAuth service
func NewPKCEOAuthService(cfg *config.Config) *PKCEOAuthService {
	return &PKCEOAuthService{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// OAuthCallbackResult represents the result of OAuth callback
type OAuthCallbackResult struct {
	Code  string
	Error error
}

// Login performs OAuth2 authentication with PKCE
func (s *PKCEOAuthService) Login(ctx context.Context) error {
	// Print stylized header
	fmt.Println("\n🔐 PipeOps CLI Authentication")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ Welcome! Let's get you authenticated with PipeOps")
	fmt.Println()

	// Step 1: Generate PKCE challenge
	fmt.Print("🔧 Generating security challenge... ")
	pkceChallenge, err := GeneratePKCEChallenge()
	if err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("failed to generate PKCE challenge: %w", err)
	}
	fmt.Println("✅ Done")

	// Step 2: Generate state parameter
	fmt.Print("🛡️  Generating security state... ")
	state, err := GenerateRandomState()
	if err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("failed to generate state: %w", err)
	}
	fmt.Println("✅ Done")

	// Build authorization URL with PKCE
	authURL := fmt.Sprintf("%s/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=%s",
		s.config.OAuth.BaseURL,
		s.config.OAuth.ClientID,
		url.QueryEscape("http://localhost:8085/callback"),
		url.QueryEscape(strings.Join(s.config.OAuth.Scopes, " ")),
		url.QueryEscape(state),
		url.QueryEscape(pkceChallenge.CodeChallenge),
		url.QueryEscape(pkceChallenge.Method),
	)

	fmt.Println()
	fmt.Println("🌐 Opening your browser for authentication...")
	fmt.Println("┌─────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                                                                                         │")
	fmt.Println("│  🚀 Your browser should open automatically                                              │")
	fmt.Println("│  📱 If not, please copy and paste the URL below:                                       │")
	fmt.Println("│                                                                                         │")
	fmt.Printf("│  🔗 %s\n", authURL)
	fmt.Println("│                                                                                         │")
	fmt.Println("│  ⏱️  You have 10 minutes to complete the authentication                                 │")
	fmt.Println("│  🔒 Your session will be securely saved locally                                        │")
	fmt.Println("│                                                                                         │")
	fmt.Println("└─────────────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Open browser
	if err := OpenBrowser(authURL); err != nil {
		fmt.Printf("⚠️  Could not open browser automatically: %v\n", err)
		fmt.Println("Please copy the URL above and paste it in your browser.")
	}

	// Start callback server
	callbackChan := make(chan OAuthCallbackResult, 1)
	server := s.startCallbackServer(callbackChan, state)
	defer server.Close()

	// Show waiting animation
	fmt.Print("⏳ Waiting for authentication")
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	dots := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Print(".")
				dots++
				if dots >= 3 {
					fmt.Print("\r⏳ Waiting for authentication   \r⏳ Waiting for authentication")
					dots = 0
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for callback
	select {
	case result := <-callbackChan:
		ticker.Stop()
		fmt.Print("\r                                        \r") // Clear waiting line
		if result.Error != nil {
			if result.Error.Error() == "callback handled" {
				return s.exchangeCodeForToken(ctx, result.Code, pkceChallenge.CodeVerifier)
			}
			return result.Error
		}
		return s.exchangeCodeForToken(ctx, result.Code, pkceChallenge.CodeVerifier)
	case <-time.After(10 * time.Minute):
		ticker.Stop()
		fmt.Print("\r                                        \r") // Clear waiting line
		fmt.Println("⏰ Authentication timeout after 10 minutes")
		fmt.Println("Please try again with 'pipeops auth login'")
		return fmt.Errorf("authentication timeout")
	case <-ctx.Done():
		ticker.Stop()
		fmt.Print("\r                                        \r") // Clear waiting line
		return ctx.Err()
	}
}

// startCallbackServer starts HTTP server for OAuth callback
func (s *PKCEOAuthService) startCallbackServer(resultChan chan<- OAuthCallbackResult, expectedState string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Close server after handling callback
			go func() {
				time.Sleep(3 * time.Second)
				resultChan <- OAuthCallbackResult{Error: fmt.Errorf("callback handled")}
			}()
		}()

		// Set content type to HTML
		w.Header().Set("Content-Type", "text/html")

		// Check for errors
		if errParam := r.URL.Query().Get("error"); errParam != "" {
			errDesc := r.URL.Query().Get("error_description")
			w.WriteHeader(400)
			errorPage := `
<!DOCTYPE html>
<html>
<head>
    <title>PipeOps Authentication - Error</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; }
        .container { max-width: 600px; margin: 0 auto; text-align: center; }
        .error-box { background: rgba(255, 255, 255, 0.1); padding: 40px; border-radius: 20px; backdrop-filter: blur(10px); }
        .error-icon { font-size: 80px; margin-bottom: 20px; }
        .error-title { font-size: 32px; margin-bottom: 20px; font-weight: 600; }
        .error-message { font-size: 18px; margin-bottom: 30px; opacity: 0.9; }
        .close-btn { background: rgba(255, 255, 255, 0.2); color: white; border: none; padding: 15px 30px; border-radius: 25px; font-size: 16px; cursor: pointer; transition: all 0.3s ease; }
        .close-btn:hover { background: rgba(255, 255, 255, 0.3); transform: translateY(-2px); }
    </style>
</head>
<body>
    <div class="container">
        <div class="error-box">
            <div class="error-icon">❌</div>
            <div class="error-title">Authentication Failed</div>
            <div class="error-message">` + errDesc + `</div>
            <button class="close-btn" onclick="window.close()">Close Window</button>
        </div>
    </div>
    <script>
        setTimeout(() => {
            window.close();
        }, 5000);
    </script>
</body>
</html>`
			w.Write([]byte(errorPage))
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("authorization error: %s - %s", errParam, errDesc)}
			return
		}

		// Verify state
		state := r.URL.Query().Get("state")
		if state != expectedState {
			w.WriteHeader(400)
			statePage := `
<!DOCTYPE html>
<html>
<head>
    <title>PipeOps Authentication - Security Error</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; }
        .container { max-width: 600px; margin: 0 auto; text-align: center; }
        .error-box { background: rgba(255, 255, 255, 0.1); padding: 40px; border-radius: 20px; backdrop-filter: blur(10px); }
        .error-icon { font-size: 80px; margin-bottom: 20px; }
        .error-title { font-size: 32px; margin-bottom: 20px; font-weight: 600; }
        .error-message { font-size: 18px; margin-bottom: 30px; opacity: 0.9; }
        .close-btn { background: rgba(255, 255, 255, 0.2); color: white; border: none; padding: 15px 30px; border-radius: 25px; font-size: 16px; cursor: pointer; transition: all 0.3s ease; }
        .close-btn:hover { background: rgba(255, 255, 255, 0.3); transform: translateY(-2px); }
    </style>
</head>
<body>
    <div class="container">
        <div class="error-box">
            <div class="error-icon">🛡️</div>
            <div class="error-title">Security Error</div>
            <div class="error-message">Invalid security state. Please try authenticating again.</div>
            <button class="close-btn" onclick="window.close()">Close Window</button>
        </div>
    </div>
    <script>
        setTimeout(() => {
            window.close();
        }, 5000);
    </script>
</body>
</html>`
			w.Write([]byte(statePage))
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("invalid state parameter")}
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			w.WriteHeader(400)
			noCodePage := `
<!DOCTYPE html>
<html>
<head>
    <title>PipeOps Authentication - No Code</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; }
        .container { max-width: 600px; margin: 0 auto; text-align: center; }
        .error-box { background: rgba(255, 255, 255, 0.1); padding: 40px; border-radius: 20px; backdrop-filter: blur(10px); }
        .error-icon { font-size: 80px; margin-bottom: 20px; }
        .error-title { font-size: 32px; margin-bottom: 20px; font-weight: 600; }
        .error-message { font-size: 18px; margin-bottom: 30px; opacity: 0.9; }
        .close-btn { background: rgba(255, 255, 255, 0.2); color: white; border: none; padding: 15px 30px; border-radius: 25px; font-size: 16px; cursor: pointer; transition: all 0.3s ease; }
        .close-btn:hover { background: rgba(255, 255, 255, 0.3); transform: translateY(-2px); }
    </style>
</head>
<body>
    <div class="container">
        <div class="error-box">
            <div class="error-icon">🔍</div>
            <div class="error-title">No Authorization Code</div>
            <div class="error-message">The authorization code was not received. Please try again.</div>
            <button class="close-btn" onclick="window.close()">Close Window</button>
        </div>
    </div>
    <script>
        setTimeout(() => {
            window.close();
        }, 5000);
    </script>
</body>
</html>`
			w.Write([]byte(noCodePage))
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("no authorization code received")}
			return
		}

		// Success response
		w.WriteHeader(200)
		successPage := `
<!DOCTYPE html>
<html>
<head>
    <title>PipeOps Authentication - Success</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 40px; background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); color: white; }
        .container { max-width: 600px; margin: 0 auto; text-align: center; }
        .success-box { background: rgba(255, 255, 255, 0.1); padding: 40px; border-radius: 20px; backdrop-filter: blur(10px); animation: slideIn 0.5s ease-out; }
        .success-icon { font-size: 80px; margin-bottom: 20px; animation: bounce 1s ease-in-out; }
        .success-title { font-size: 32px; margin-bottom: 20px; font-weight: 600; }
        .success-message { font-size: 18px; margin-bottom: 30px; opacity: 0.9; }
        .close-btn { background: rgba(255, 255, 255, 0.2); color: white; border: none; padding: 15px 30px; border-radius: 25px; font-size: 16px; cursor: pointer; transition: all 0.3s ease; }
        .close-btn:hover { background: rgba(255, 255, 255, 0.3); transform: translateY(-2px); }
        @keyframes slideIn { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
        @keyframes bounce { 0%, 20%, 50%, 80%, 100% { transform: translateY(0); } 40% { transform: translateY(-20px); } 60% { transform: translateY(-10px); } }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-box">
            <div class="success-icon">🎉</div>
            <div class="success-title">Authentication Successful!</div>
            <div class="success-message">You're now authenticated with PipeOps CLI. You can close this window and return to your terminal.</div>
            <button class="close-btn" onclick="window.close()">Close Window</button>
        </div>
    </div>
    <script>
        setTimeout(() => {
            window.close();
        }, 5000);
    </script>
</body>
</html>`
		w.Write([]byte(successPage))
		resultChan <- OAuthCallbackResult{Code: code}
	})

	server := &http.Server{
		Addr:    ":8085",
		Handler: mux,
	}

	go server.ListenAndServe()
	return server
}

// exchangeCodeForToken exchanges authorization code for access token using PKCE
func (s *PKCEOAuthService) exchangeCodeForToken(ctx context.Context, code, codeVerifier string) error {
	fmt.Println("🔗 Received authentication code!")
	fmt.Print("🔄 Exchanging code for access token... ")

	// Prepare token request with PKCE (no client secret needed for public clients)
	tokenReq := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  "http://localhost:8085/callback",
		"client_id":     s.config.OAuth.ClientID,
		"code_verifier": codeVerifier, // PKCE code verifier
	}

	jsonData, err := json.Marshal(tokenReq)
	if err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("failed to marshal token request: %w", err)
	}

	// Make token request
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.OAuth.BaseURL+"/oauth/token", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("❌ Failed")
		fmt.Printf("🚫 Token exchange failed (status %d): %s\n", resp.StatusCode, string(body))
		return fmt.Errorf("token exchange failed: %s", string(body))
	}

	// Parse token response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		fmt.Println("❌ Failed")
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	fmt.Println("✅ Success")

	// Save tokens
	s.config.OAuth.AccessToken = tokenResp.AccessToken
	s.config.OAuth.RefreshToken = tokenResp.RefreshToken
	s.config.OAuth.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Display success message
	fmt.Println()
	fmt.Println("🎉 Authentication Complete!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Successfully authenticated with PipeOps")
	fmt.Printf("🔑 Access token expires: %s\n", s.config.OAuth.ExpiresAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("⏰ Token valid for: %d hours\n", tokenResp.ExpiresIn/3600)
	fmt.Println("💾 Credentials saved securely to your local config")
	fmt.Println()
	fmt.Println("🚀 You can now use all PipeOps CLI commands!")
	fmt.Println("   Try: pipeops auth me")
	fmt.Println("   Or:  pipeops project list")
	fmt.Println()

	return nil
}

// IsAuthenticated checks if user is authenticated
func (s *PKCEOAuthService) IsAuthenticated() bool {
	return s.config.OAuth.AccessToken != "" && time.Now().Before(s.config.OAuth.ExpiresAt.Add(-5*time.Minute))
}

// GetAccessToken returns the current access token
func (s *PKCEOAuthService) GetAccessToken() string {
	return s.config.OAuth.AccessToken
}
