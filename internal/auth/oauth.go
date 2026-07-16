package auth

import (
	"bytes"
	"context"
	"crypto/subtle"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

// PKCEOAuthService handles OAuth2 authentication with PKCE
type PKCEOAuthService struct {
	config       *config.Config
	callbackPort int
}

// NewPKCEOAuthService creates a new PKCE OAuth service
func NewPKCEOAuthService(cfg *config.Config) *PKCEOAuthService {
	return &PKCEOAuthService{
		config: cfg,
	}
}

// OAuthCallbackResult represents the result of OAuth callback
type OAuthCallbackResult struct {
	Code  string
	Token string // Direct token if server returns it
	Error error
}

var oauthCallbackPageTemplate = template.Must(template.New("oauth-callback-page").Parse(`<!DOCTYPE html>
<html>
<head>
    <title>PipeOps - {{.Title}}</title>
    <style>
        :root { --bg: #0F172A; --card: #1E293B; --text: #F8FAFC; --subtext: #94A3B8; --accent: #6366F1; --error: #EF4444; --warning: #F59E0B; --success: #10B981; }
        body { font-family: -apple-system, system-ui, sans-serif; background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
        .card { background: var(--card); padding: 2.5rem; border-radius: 1rem; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5); width: 90%; max-width: 400px; text-align: center; border: 1px solid rgba(255,255,255,0.05); }
        .icon { width: 64px; height: 64px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 32px; margin: 0 auto 1.5rem; }
        .icon.error { background: rgba(239,68,68,0.1); color: var(--error); }
        .icon.warning { background: rgba(245,158,11,0.1); color: var(--warning); }
        .icon.success { background: rgba(16,185,129,0.1); color: var(--success); }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; font-weight: 600; }
        p { color: var(--subtext); line-height: 1.5; margin-bottom: 2rem; }
        .btn { display: block; width: 100%; padding: 0.75rem; background: transparent; color: var(--text); border: 1px solid #334155; border-radius: 0.5rem; font-weight: 500; cursor: pointer; transition: all 0.2s; font-size: 1rem; }
        .btn.success { background: var(--accent); color: white; border: none; }
        .btn:hover { opacity: 0.9; background: rgba(255,255,255,0.05); border-color: #475569; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon {{.Tone}}">{{.Icon}}</div>
        <h1>{{.Heading}}</h1>
        <p>{{.Message}}</p>
        <button class="btn {{.Tone}}" onclick="window.close()">Close Window</button>
    </div>
    {{if .AutoClose}}<script>setTimeout(() => window.close(), 3000);</script>{{end}}
</body>
</html>`))

type oauthCallbackPageData struct {
	Title     string
	Heading   string
	Message   string
	Icon      string
	Tone      string
	AutoClose bool
}

func writeOAuthCallbackPage(w http.ResponseWriter, status int, data oauthCallbackPageData) {
	var body bytes.Buffer
	if err := oauthCallbackPageTemplate.Execute(&body, data); err != nil {
		http.Error(w, "authentication callback failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(body.Bytes())
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
	s.config.OAuth.RefreshToken = ""                               // Direct tokens don't have refresh tokens
	s.config.OAuth.ExpiresAt = time.Now().Add(30 * 24 * time.Hour) // Default 30 day expiry

	fmt.Println("🎉 Authentication successful!")
	fmt.Println("✅ You're now logged in to PipeOps")
	return nil
}

// findAvailablePort finds an available port for the callback server
func (s *PKCEOAuthService) findAvailablePort() (int, error) {
	// Try preferred ports first
	preferredPorts := []int{8080, 8085, 8086, 8087, 8088, 8089}
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
			writeOAuthCallbackPage(w, http.StatusBadRequest, oauthCallbackPageData{
				Title:   "Authentication Failed",
				Heading: "Authentication Failed",
				Message: errDesc,
				Icon:    "✕",
				Tone:    "error",
			})
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("authorization error: %s - %s", config.SanitizeLog(errParam), config.SanitizeLog(errDesc))}
			return
		}

		// Verify state
		state := r.URL.Query().Get("state")
		if !validOAuthState(state, expectedState) {
			writeOAuthCallbackPage(w, http.StatusBadRequest, oauthCallbackPageData{
				Title:   "Security Check Failed",
				Heading: "Security Check Failed",
				Message: "Invalid security state parameter. Please try authenticating again.",
				Icon:    "🛡️",
				Tone:    "warning",
			})
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
			fmt.Printf("🔍 Debug: Callback query params: %v\n", config.SanitizeLog(fmt.Sprintf("%v", r.URL.Query())))
		}

		if code == "" && token == "" {
			writeOAuthCallbackPage(w, http.StatusBadRequest, oauthCallbackPageData{
				Title:   "Authorization Failed",
				Heading: "Authorization Failed",
				Message: "The authorization code was not received. Please try authenticating again.",
				Icon:    "🔍",
				Tone:    "warning",
			})
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("no authorization code received")}
			return
		}

		// Success response
		writeOAuthCallbackPage(w, http.StatusOK, oauthCallbackPageData{
			Title:     "Authenticated",
			Heading:   "Authenticated",
			Message:   "You have successfully signed in to PipeOps CLI. You can now close this window and return to your terminal.",
			Icon:      "✓",
			Tone:      "success",
			AutoClose: true,
		})
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

func validOAuthState(received, expected string) bool {
	if expected == "" || received == "" {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(received), []byte(expected)) == 1 {
		return true
	}

	receivedWithoutPadding := strings.TrimRight(received, "=")
	expectedWithoutPadding := strings.TrimRight(expected, "=")
	return subtle.ConstantTimeCompare([]byte(receivedWithoutPadding), []byte(expectedWithoutPadding)) == 1
}

// exchangeCodeForToken exchanges authorization code for access token using PKCE
func (s *PKCEOAuthService) exchangeCodeForToken(ctx context.Context, code, codeVerifier string) error {
	client, err := s.newSDKClient()
	if err != nil {
		return err
	}

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

	req, err := client.NewRequest(http.MethodPost, "oauth/token", tokenReq)
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
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

	if _, err := client.Do(ctx, req, &tokenResp); err != nil {
		return fmt.Errorf("token exchange failed: %w", err)
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

	client, err := s.newSDKClient()
	if err != nil {
		return err
	}

	// Prepare refresh request
	refreshReq := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": s.config.OAuth.RefreshToken,
		"client_id":     s.config.OAuth.ClientID,
	}

	req, err := client.NewRequest(http.MethodPost, "oauth/token", refreshReq)
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	// Parse refresh response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if _, err := client.Do(ctx, req, &tokenResp); err != nil {
		return fmt.Errorf("refresh failed: %w", err)
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

func (s *PKCEOAuthService) newSDKClient() (*sdk.Client, error) {
	baseURL := config.GetAPIURL()
	if s.config != nil && s.config.OAuth != nil && strings.TrimSpace(s.config.OAuth.BaseURL) != "" {
		baseURL = s.config.OAuth.BaseURL
	}

	client, err := sdk.NewClient(baseURL, sdk.WithTimeout(30*time.Second), sdk.WithMaxRetries(3))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PipeOps SDK client: %w", err)
	}
	if s.config != nil && s.config.OAuth != nil && strings.TrimSpace(s.config.OAuth.AccessToken) != "" {
		client.SetToken(s.config.OAuth.AccessToken)
	}
	return client, nil
}

// IsAuthenticated checks if user is authenticated and attempts refresh if expired
func (s *PKCEOAuthService) IsAuthenticated() bool {
	if s.config == nil || s.config.OAuth == nil {
		return false
	}
	if s.config.IsAuthenticated() {
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
