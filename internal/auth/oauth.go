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
	// Generate PKCE challenge
	pkceChallenge, err := GeneratePKCEChallenge()
	if err != nil {
		return fmt.Errorf("failed to generate PKCE challenge: %w", err)
	}

	// Generate state parameter
	state, err := GenerateRandomState()
	if err != nil {
		return fmt.Errorf("failed to generate state: %w", err)
	}

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

	fmt.Println("Opening browser for authentication...")
	fmt.Printf("If browser doesn't open, visit: %s\n", authURL)

	// Open browser
	if err := OpenBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser: %v\n", err)
		fmt.Println("Please visit the URL above manually.")
	}

	// Start callback server
	callbackChan := make(chan OAuthCallbackResult, 1)
	server := s.startCallbackServer(callbackChan, state)
	defer server.Close()

	// Wait for callback with simple status
	fmt.Print("Waiting for authentication...")

	select {
	case result := <-callbackChan:
		fmt.Print("\r                                \r") // Clear line
		if result.Error != nil {
			if result.Error.Error() == "callback handled" {
				return s.exchangeCodeForToken(ctx, result.Code, pkceChallenge.CodeVerifier)
			}
			return result.Error
		}
		return s.exchangeCodeForToken(ctx, result.Code, pkceChallenge.CodeVerifier)
	case <-time.After(10 * time.Minute):
		fmt.Print("\r                                \r") // Clear line
		return fmt.Errorf("authentication timeout - please try again")
	case <-ctx.Done():
		fmt.Print("\r                                \r") // Clear line
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
	// Prepare token request with PKCE (no client secret needed for public clients)
	tokenReq := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  "http://localhost:8085/callback",
		"client_id":     s.config.OAuth.ClientID,
		"code_verifier": codeVerifier, // PKCE code verifier
		"token_format":  "jwt",        // Request JWT tokens if server supports it
	}

	jsonData, err := json.Marshal(tokenReq)
	if err != nil {
		return fmt.Errorf("failed to marshal token request: %w", err)
	}

	// Make token request
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.OAuth.BaseURL+"/oauth/token", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != 200 {
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
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Save tokens
	s.config.OAuth.AccessToken = tokenResp.AccessToken
	s.config.OAuth.RefreshToken = tokenResp.RefreshToken
	s.config.OAuth.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	fmt.Println("Authentication successful")
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
