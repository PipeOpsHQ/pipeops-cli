package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
)

// PKCEOAuthService handles OAuth2 authentication with PKCE
type PKCEOAuthService struct {
	config       *config.Config
	client       *http.Client
	callbackPort int
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
	Token string // Direct token if server returns it
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

	// Find available port for callback server
	port, err := s.findAvailablePort()
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}
	s.callbackPort = port

	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	// Build authorization URL with PKCE
	authURL := fmt.Sprintf("%s/auth/signin?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=%s&oauth=true",
		s.config.OAuth.DashboardURL,
		s.config.OAuth.ClientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape(strings.Join(s.config.OAuth.Scopes, " ")),
		url.QueryEscape(state),
		url.QueryEscape(pkceChallenge.CodeChallenge),
		url.QueryEscape(pkceChallenge.Method),
	)

	fmt.Println("🔐 Starting secure authentication...")
	fmt.Println("→ Opening your browser for PipeOps login")
	fmt.Printf("  If it doesn't open automatically, visit:\n  %s\n", authURL)
	fmt.Println()

	// Open browser
	if err := OpenBrowser(authURL); err != nil {
		fmt.Printf("⚠️  Browser didn't open automatically: %v\n", err)
		fmt.Println("   No worries! Just copy the URL above and paste it in your browser")
		fmt.Println()
	}

	// Start callback server
	callbackChan := make(chan OAuthCallbackResult, 1)
	server, err := s.startCallbackServer(callbackChan, state)
	if err != nil {
		return fmt.Errorf("failed to start callback server: %w", err)
	}
	defer server.Close()

	// Wait for callback with encouraging feedback
	fmt.Print("⏳ Waiting for you to complete authentication in your browser...")

	select {
	case result := <-callbackChan:
		fmt.Print("\r                                                                \r") // Clear line
		if result.Error != nil {
			if result.Error.Error() == "callback handled" {
				// If token was returned directly, use it
				if result.Token != "" {
					return s.handleDirectToken(result.Token)
				}
				return s.exchangeCodeForToken(ctx, result.Code, pkceChallenge.CodeVerifier)
			}
			return result.Error
		}
		// If token was returned directly, use it
		if result.Token != "" {
			return s.handleDirectToken(result.Token)
		}
		return s.exchangeCodeForToken(ctx, result.Code, pkceChallenge.CodeVerifier)
	case <-time.After(10 * time.Minute):
		fmt.Print("\r                                                                \r") // Clear line
		fmt.Println("⏰ Authentication timed out after 10 minutes")
		fmt.Println("   No problem! Just run 'pipeops auth login' again when ready")
		return fmt.Errorf("authentication timeout")
	case <-ctx.Done():
		fmt.Print("\r                                                                \r") // Clear line
		return ctx.Err()
	}
}

// handleDirectToken handles a token that was returned directly in the callback
func (s *PKCEOAuthService) handleDirectToken(token string) error {
	s.config.OAuth.AccessToken = token
	s.config.OAuth.RefreshToken = "" // Direct tokens don't have refresh tokens
	s.config.OAuth.ExpiresAt = time.Now().Add(30 * 24 * time.Hour) // Default 30 day expiry

	fmt.Println("🎉 Authentication successful!")
	fmt.Println("✅ You're now logged in to PipeOps")
	return nil
}

// findAvailablePort finds an available port for the callback server
func (s *PKCEOAuthService) findAvailablePort() (int, error) {
	// Try preferred ports first
	preferredPorts := []int{8085, 8086, 8087, 8088, 8089}
	for _, port := range preferredPorts {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			listener.Close()
			return port, nil
		}
	}

	// If all preferred ports are taken, let the OS assign a port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("failed to find available port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port, nil
}

// startCallbackServer starts HTTP server for OAuth callback
func (s *PKCEOAuthService) startCallbackServer(resultChan chan<- OAuthCallbackResult, expectedState string) (*http.Server, error) {
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
    <title>PipeOps - Authentication Failed</title>
    <style>
        :root { --bg: #0F172A; --card: #1E293B; --text: #F8FAFC; --subtext: #94A3B8; --error: #EF4444; }
        body { font-family: -apple-system, system-ui, sans-serif; background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: var(--card); padding: 2.5rem; border-radius: 1rem; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); width: 90%; max-width: 400px; text-align: center; border: 1px solid rgba(255,255,255,0.05); }
        .icon { width: 64px; height: 64px; background: rgba(239,68,68,0.1); color: var(--error); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 32px; margin: 0 auto 1.5rem; }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; font-weight: 600; }
        p { color: var(--subtext); line-height: 1.5; margin-bottom: 2rem; }
        .btn { display: block; width: 100%; padding: 0.75rem; background: transparent; border: 1px solid #334155; color: var(--text); border-radius: 0.5rem; font-weight: 500; cursor: pointer; transition: all 0.2s; font-size: 1rem; }
        .btn:hover { background: rgba(255,255,255,0.05); border-color: #475569; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon">✕</div>
        <h1>Authentication Failed</h1>
        <p>` + errDesc + `</p>
        <button class="btn" onclick="window.close()">Close Window</button>
    </div>
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
    <title>PipeOps - Security Check Failed</title>
    <style>
        :root { --bg: #0F172A; --card: #1E293B; --text: #F8FAFC; --subtext: #94A3B8; --warning: #F59E0B; }
        body { font-family: -apple-system, system-ui, sans-serif; background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: var(--card); padding: 2.5rem; border-radius: 1rem; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); width: 90%; max-width: 400px; text-align: center; border: 1px solid rgba(255,255,255,0.05); }
        .icon { width: 64px; height: 64px; background: rgba(245,158,11,0.1); color: var(--warning); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 32px; margin: 0 auto 1.5rem; }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; font-weight: 600; }
        p { color: var(--subtext); line-height: 1.5; margin-bottom: 2rem; }
        .btn { display: block; width: 100%; padding: 0.75rem; background: transparent; border: 1px solid #334155; color: var(--text); border-radius: 0.5rem; font-weight: 500; cursor: pointer; transition: all 0.2s; font-size: 1rem; }
        .btn:hover { background: rgba(255,255,255,0.05); border-color: #475569; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon">🛡️</div>
        <h1>Security Check Failed</h1>
        <p>Invalid security state parameter. Please try authenticating again.</p>
        <button class="btn" onclick="window.close()">Close Window</button>
    </div>
</body>
</html>`
			w.Write([]byte(statePage))
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("invalid state parameter")}
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		
		// Also check for token parameter (some servers return token directly)
		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.URL.Query().Get("access_token")
		}
		
		// Debug: log all query parameters
		if s.config.Settings != nil && s.config.Settings.Debug {
			fmt.Printf("🔍 Debug: Callback query params: %v\n", r.URL.Query())
		}
		
		if code == "" && token == "" {
			w.WriteHeader(400)
			noCodePage := `
<!DOCTYPE html>
<html>
<head>
    <title>PipeOps - Authorization Failed</title>
    <style>
        :root { --bg: #0F172A; --card: #1E293B; --text: #F8FAFC; --subtext: #94A3B8; --warning: #F59E0B; }
        body { font-family: -apple-system, system-ui, sans-serif; background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: var(--card); padding: 2.5rem; border-radius: 1rem; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); width: 90%; max-width: 400px; text-align: center; border: 1px solid rgba(255,255,255,0.05); }
        .icon { width: 64px; height: 64px; background: rgba(245,158,11,0.1); color: var(--warning); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 32px; margin: 0 auto 1.5rem; }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; font-weight: 600; }
        p { color: var(--subtext); line-height: 1.5; margin-bottom: 2rem; }
        .btn { display: block; width: 100%; padding: 0.75rem; background: transparent; border: 1px solid #334155; color: var(--text); border-radius: 0.5rem; font-weight: 500; cursor: pointer; transition: all 0.2s; font-size: 1rem; }
        .btn:hover { background: rgba(255,255,255,0.05); border-color: #475569; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon">🔍</div>
        <h1>Authorization Failed</h1>
        <p>The authorization code was not received. Please try authenticating again.</p>
        <button class="btn" onclick="window.close()">Close Window</button>
    </div>
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
    <title>PipeOps - Authenticated</title>
    <style>
        :root { --bg: #0F172A; --card: #1E293B; --text: #F8FAFC; --subtext: #94A3B8; --accent: #6366F1; --success: #10B981; }
        body { font-family: -apple-system, system-ui, sans-serif; background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: var(--card); padding: 2.5rem; border-radius: 1rem; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); width: 90%; max-width: 400px; text-align: center; border: 1px solid rgba(255,255,255,0.05); }
        .icon { width: 64px; height: 64px; background: rgba(16,185,129,0.1); color: var(--success); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 32px; margin: 0 auto 1.5rem; }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; font-weight: 600; }
        p { color: var(--subtext); line-height: 1.5; margin-bottom: 2rem; }
        .btn { display: block; width: 100%; padding: 0.75rem; background: var(--accent); color: white; border: none; border-radius: 0.5rem; font-weight: 500; cursor: pointer; transition: opacity 0.2s; font-size: 1rem; }
        .btn:hover { opacity: 0.9; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon">✓</div>
        <h1>Authenticated</h1>
        <p>You have successfully signed in to PipeOps CLI. You can now close this window and return to your terminal.</p>
        <button class="btn" onclick="window.close()">Close Window</button>
    </div>
    <script>setTimeout(() => window.close(), 3000);</script>
</body>
</html>`
		w.Write([]byte(successPage))
		resultChan <- OAuthCallbackResult{Code: code, Token: token}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.callbackPort),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("callback server error: %w", err)}
		}
	}()
	return server, nil
}

// exchangeCodeForToken exchanges authorization code for access token using PKCE
func (s *PKCEOAuthService) exchangeCodeForToken(ctx context.Context, code, codeVerifier string) error {
	// Prepare token request with PKCE (no client secret needed for public clients)
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", s.callbackPort)
	tokenReq := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectURI,
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

	// Debug: log raw response
	if s.config.Settings != nil && s.config.Settings.Debug {
		fmt.Printf("🔍 Debug: Token response: %s\n", string(body))
	}

	// Parse token response - new format includes redirect_url
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		RedirectURL  string `json:"redirect_url,omitempty"` // New field for redirect handling
		// Also check for data wrapper format
		Data struct {
			Token string `json:"token"`
		} `json:"data,omitempty"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Check for token in data wrapper (alternative format)
	accessToken := tokenResp.AccessToken
	if accessToken == "" && tokenResp.Data.Token != "" {
		accessToken = tokenResp.Data.Token
	}

	// Debug: show token format
	if s.config.Settings != nil && s.config.Settings.Debug {
		segments := strings.Count(accessToken, ".")
		if segments == 2 {
			fmt.Printf("🔍 Debug: Token format: JWT (3 segments)\n")
		} else {
			fmt.Printf("🔍 Debug: Token format: Opaque (%d segments)\n", segments+1)
		}
	}

	// Save tokens
	s.config.OAuth.AccessToken = accessToken
	s.config.OAuth.RefreshToken = tokenResp.RefreshToken
	s.config.OAuth.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Handle redirect URL if provided by the API
	if tokenResp.RedirectURL != "" {
		fmt.Printf("🔗 API provided redirect URL: %s\n", tokenResp.RedirectURL)
		// Optionally open the redirect URL in browser for any post-auth steps
		if err := OpenBrowser(tokenResp.RedirectURL); err != nil {
			fmt.Printf("⚠️  Could not open redirect URL automatically: %v\n", err)
			fmt.Printf("   You can manually visit: %s\n", tokenResp.RedirectURL)
		}
	}

	fmt.Println("🎉 Authentication successful!")
	fmt.Println("✅ You're now logged in to PipeOps")
	return nil
}

// Refresh uses the refresh token to obtain a new access token
func (s *PKCEOAuthService) Refresh(ctx context.Context) error {
	if s.config.OAuth.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// Prepare refresh request
	refreshReq := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": s.config.OAuth.RefreshToken,
		"client_id":     s.config.OAuth.ClientID,
	}

	jsonData, err := json.Marshal(refreshReq)
	if err != nil {
		return fmt.Errorf("failed to marshal refresh request: %w", err)
	}

	// Make refresh request
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.OAuth.BaseURL+"/oauth/token", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse refresh response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse refresh response: %w", err)
	}

	// Update config
	s.config.OAuth.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		s.config.OAuth.RefreshToken = tokenResp.RefreshToken
	}
	s.config.OAuth.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Save updated config
	if err := config.Save(s.config); err != nil {
		return fmt.Errorf("failed to save refreshed config: %w", err)
	}

	return nil
}

// IsAuthenticated checks if user is authenticated and attempts refresh if expired
func (s *PKCEOAuthService) IsAuthenticated() bool {
	if s.config.OAuth.AccessToken != "" && time.Now().Before(s.config.OAuth.ExpiresAt.Add(-5*time.Minute)) {
		return true
	}

	if s.config.OAuth.RefreshToken != "" {
		ctx := context.Background() // Use background context for simplicity; consider passing ctx if needed
		if err := s.Refresh(ctx); err == nil {
			return true
		} else {
			// Optional: log error or clear tokens if refresh fails
			fmt.Printf("Token refresh failed: %v\n", err)
		}
	}

	return false
}

// GetAccessToken returns the current access token
func (s *PKCEOAuthService) GetAccessToken() string {
	return s.config.OAuth.AccessToken
}
