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
	fmt.Println("🚀 Starting PipeOps CLI authentication...")

	// Generate PKCE challenge
	pkceChallenge, err := GeneratePKCEChallenge()
	if err != nil {
		return fmt.Errorf("failed to generate PKCE challenge: %w", err)
	}

	// Generate state parameter for CSRF protection
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

	fmt.Println("📱 Opening browser for authentication...")
	fmt.Printf("🌐 Please visit: %s\n", authURL)

	// Open browser
	if err := OpenBrowser(authURL); err != nil {
		fmt.Printf("⚠️  Could not open browser: %v\n", err)
		fmt.Println("Please copy the URL above and paste it in your browser.")
	}

	// Start callback server
	callbackChan := make(chan OAuthCallbackResult, 1)
	server := s.startCallbackServer(callbackChan, state)
	defer server.Close()

	// Wait for callback
	select {
	case result := <-callbackChan:
		if result.Error != nil {
			return result.Error
		}
		return s.exchangeCodeForToken(ctx, result.Code, pkceChallenge.CodeVerifier)
	case <-time.After(10 * time.Minute):
		return fmt.Errorf("authentication timeout")
	case <-ctx.Done():
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
				time.Sleep(2 * time.Second)
				resultChan <- OAuthCallbackResult{Error: fmt.Errorf("callback handled")}
			}()
		}()

		// Check for errors
		if errParam := r.URL.Query().Get("error"); errParam != "" {
			errDesc := r.URL.Query().Get("error_description")
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("❌ Authentication failed: %s", errDesc)))
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("authorization error: %s - %s", errParam, errDesc)}
			return
		}

		// Verify state
		state := r.URL.Query().Get("state")
		if state != expectedState {
			w.WriteHeader(400)
			w.Write([]byte("❌ Invalid state parameter"))
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("invalid state parameter")}
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			w.WriteHeader(400)
			w.Write([]byte("❌ No authorization code received"))
			resultChan <- OAuthCallbackResult{Error: fmt.Errorf("no authorization code received")}
			return
		}

		// Success response
		w.WriteHeader(200)
		w.Write([]byte("✅ Authentication successful! You can close this window."))
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
	fmt.Println("🔄 Exchanging authorization code for access token...")

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
		return fmt.Errorf("failed to marshal token request: %w", err)
	}

	fmt.Printf("🔍 Token request (PKCE): %s\n", string(jsonData))

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
		fmt.Printf("❌ Token exchange failed (status %d): %s\n", resp.StatusCode, string(body))
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

	fmt.Println("✅ Authentication successful!")
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
