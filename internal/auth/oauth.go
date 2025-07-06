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

// OAuthService handles OAuth 2.0 authentication with PKCE
type OAuthService struct {
	config     *config.OAuthConfig
	httpClient *http.Client
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(cfg *config.OAuthConfig) *OAuthService {
	return &OAuthService{
		config:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

// UserInfo represents user information from the API
type UserInfo struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AuthCallbackResult represents the result of the OAuth callback
type AuthCallbackResult struct {
	Code  string
	Error error
}

// Login performs the complete OAuth login flow with PKCE
func (s *OAuthService) Login(ctx context.Context) error {
	fmt.Println("🔐 Starting OAuth login...")

	// Generate PKCE challenge
	pkce, err := GeneratePKCEChallenge()
	if err != nil {
		return fmt.Errorf("failed to generate PKCE challenge: %w", err)
	}

	// Generate state parameter for CSRF protection
	state, err := GenerateRandomState()
	if err != nil {
		return fmt.Errorf("failed to generate state parameter: %w", err)
	}

	// Build authorization URL
	authURL := s.buildAuthorizationURL(pkce, state)

	// Start local callback server
	callbackCh := make(chan AuthCallbackResult, 1)
	server := s.startCallbackServer(callbackCh, state)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	// Open browser
	fmt.Printf("📱 Opening browser for authentication...\n")
	fmt.Printf("🌐 If browser doesn't open automatically, visit:\n%s\n\n", authURL)

	if err := OpenBrowser(authURL); err != nil {
		fmt.Printf("⚠️  Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please copy and paste the URL above into your browser.\n\n")
	}

	// Wait for callback or timeout
	select {
	case result := <-callbackCh:
		if result.Error != nil {
			return fmt.Errorf("authentication failed: %w", result.Error)
		}
		return s.exchangeCodeForToken(ctx, result.Code, pkce.CodeVerifier)

	case <-time.After(10 * time.Minute):
		return fmt.Errorf("authentication timeout after 10 minutes")

	case <-ctx.Done():
		return fmt.Errorf("authentication cancelled: %w", ctx.Err())
	}
}

// buildAuthorizationURL constructs the OAuth authorization URL
func (s *OAuthService) buildAuthorizationURL(pkce *PKCEChallenge, state string) string {
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {s.config.ClientID},
		"redirect_uri":          {"http://localhost:8085/callback"},
		"scope":                 {strings.Join(s.config.Scopes, " ")},
		"code_challenge":        {pkce.CodeChallenge},
		"code_challenge_method": {pkce.Method},
		"state":                 {state},
	}

	return fmt.Sprintf("%s/oauth/authorize?%s", s.config.BaseURL, params.Encode())
}

// startCallbackServer starts a local HTTP server to handle the OAuth callback
func (s *OAuthService) startCallbackServer(resultCh chan<- AuthCallbackResult, expectedState string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Check for error parameter
		if errorParam := r.URL.Query().Get("error"); errorParam != "" {
			errorDesc := r.URL.Query().Get("error_description")
			resultCh <- AuthCallbackResult{
				Error: fmt.Errorf("authorization error: %s - %s", errorParam, errorDesc),
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(s.getErrorPage(errorParam, errorDesc)))
			return
		}

		// Verify state parameter (CSRF protection)
		state := r.URL.Query().Get("state")
		if state != expectedState {
			resultCh <- AuthCallbackResult{
				Error: fmt.Errorf("invalid state parameter - possible CSRF attack"),
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(s.getErrorPage("invalid_state", "State parameter mismatch")))
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			resultCh <- AuthCallbackResult{
				Error: fmt.Errorf("no authorization code received"),
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(s.getErrorPage("no_code", "No authorization code received")))
			return
		}

		// Send success response
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s.getSuccessPage()))

		// Send result
		resultCh <- AuthCallbackResult{Code: code}
	})

	server := &http.Server{
		Addr:         ":8085",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			resultCh <- AuthCallbackResult{
				Error: fmt.Errorf("callback server error: %w", err),
			}
		}
	}()

	return server
}

// exchangeCodeForToken exchanges the authorization code for an access token
func (s *OAuthService) exchangeCodeForToken(ctx context.Context, code, codeVerifier string) error {
	fmt.Println("🔄 Exchanging authorization code for access token...")

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {s.config.ClientID},
		"code":          {code},
		"redirect_uri":  {"http://localhost:8085/callback"},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		s.config.BaseURL+"/oauth/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	// Update configuration with new tokens
	s.config.AccessToken = tokenResp.AccessToken
	s.config.RefreshToken = tokenResp.RefreshToken
	s.config.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	fmt.Println("✅ Authentication successful!")
	return nil
}

// GetUserInfo retrieves the current user's information
func (s *OAuthService) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	if !s.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated - please run 'pipeops auth login'")
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		s.config.BaseURL+"/oauth/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.config.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user info request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle unauthorized response (try refresh if available)
	if resp.StatusCode == http.StatusUnauthorized {
		if s.config.RefreshToken != "" {
			if refreshErr := s.RefreshToken(ctx); refreshErr == nil {
				// Retry with new token
				return s.GetUserInfo(ctx)
			}
		}
		return nil, fmt.Errorf("authentication expired - please run 'pipeops auth login'")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// IsAuthenticated checks if the user has a valid access token
func (s *OAuthService) IsAuthenticated() bool {
	return s.config.AccessToken != "" && time.Now().Before(s.config.ExpiresAt.Add(-5*time.Minute))
}

// RefreshToken refreshes the access token using the refresh token
func (s *OAuthService) RefreshToken(ctx context.Context) error {
	if s.config.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	fmt.Println("🔄 Refreshing access token...")

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {s.config.ClientID},
		"refresh_token": {s.config.RefreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		s.config.BaseURL+"/oauth/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to decode refresh response: %w", err)
	}

	// Update stored tokens
	s.config.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		s.config.RefreshToken = tokenResp.RefreshToken
	}
	s.config.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	fmt.Println("✅ Token refreshed successfully!")
	return nil
}

// GetAccessToken returns the current access token
func (s *OAuthService) GetAccessToken() string {
	return s.config.AccessToken
}

// HTML page templates for OAuth callback
func (s *OAuthService) getSuccessPage() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>PipeOps CLI - Authentication Successful</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; text-align: center; padding: 50px; background: #f8f9fa; }
        .container { max-width: 500px; margin: 0 auto; background: white; padding: 40px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .success { color: #28a745; font-size: 48px; margin-bottom: 20px; }
        h1 { color: #333; margin-bottom: 20px; }
        p { color: #666; line-height: 1.6; }
        .close-instruction { background: #e9ecef; padding: 15px; border-radius: 4px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="success">✅</div>
        <h1>Authentication Successful!</h1>
        <p>You have successfully logged in to PipeOps CLI.</p>
        <div class="close-instruction">
            <strong>You can now close this browser tab and return to your terminal.</strong>
        </div>
    </div>
    <script>
        // Auto-close tab after 3 seconds
        setTimeout(() => {
            window.close();
        }, 3000);
    </script>
</body>
</html>`
}

func (s *OAuthService) getErrorPage(error, description string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>PipeOps CLI - Authentication Error</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; text-align: center; padding: 50px; background: #f8f9fa; }
        .container { max-width: 500px; margin: 0 auto; background: white; padding: 40px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .error { color: #dc3545; font-size: 48px; margin-bottom: 20px; }
        h1 { color: #333; margin-bottom: 20px; }
        p { color: #666; line-height: 1.6; }
        .error-details { background: #f8d7da; padding: 15px; border-radius: 4px; margin-top: 20px; border: 1px solid #f5c6cb; }
    </style>
</head>
<body>
    <div class="container">
        <div class="error">❌</div>
        <h1>Authentication Failed</h1>
        <p>There was an error during the authentication process.</p>
        <div class="error-details">
            <strong>Error:</strong> %s<br>
            <strong>Description:</strong> %s
        </div>
        <p>Please close this tab and try again in your terminal.</p>
    </div>
</body>
</html>`, error, description)
}
