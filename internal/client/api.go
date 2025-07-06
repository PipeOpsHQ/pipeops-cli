package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
)

// AuthenticatedClient wraps http.Client with automatic OAuth authentication
type AuthenticatedClient struct {
	baseURL     string
	httpClient  *http.Client
	authService *auth.OAuthService
	config      *config.Config
}

// NewAuthenticatedClient creates a new authenticated HTTP client
func NewAuthenticatedClient(cfg *config.Config) (*AuthenticatedClient, error) {
	if cfg.OAuth == nil {
		return nil, fmt.Errorf("OAuth configuration not found")
	}

	if !cfg.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated - please run 'pipeops auth login'")
	}

	authService := auth.NewOAuthService(cfg.OAuth)

	return &AuthenticatedClient{
		baseURL:     cfg.OAuth.BaseURL,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		authService: authService,
		config:      cfg,
	}, nil
}

// Do performs an HTTP request with automatic authentication
func (c *AuthenticatedClient) Do(req *http.Request) (*http.Response, error) {
	// Check if token is still valid
	if !c.authService.IsAuthenticated() {
		return nil, fmt.Errorf("authentication expired - please run 'pipeops auth login'")
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+c.authService.GetAccessToken())
	req.Header.Set("User-Agent", "PipeOps-CLI/1.0")

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Handle token refresh on 401 Unauthorized
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()

		// Try to refresh token if available
		if c.config.OAuth.RefreshToken != "" {
			if err := c.authService.RefreshToken(req.Context()); err == nil {
				// Save updated config
				if saveErr := config.Save(c.config); saveErr != nil {
					fmt.Printf("Warning: Failed to save refreshed token: %v\n", saveErr)
				}

				// Retry request with new token
				req.Header.Set("Authorization", "Bearer "+c.authService.GetAccessToken())
				return c.httpClient.Do(req)
			}
		}

		return nil, fmt.Errorf("authentication expired - please run 'pipeops auth login'")
	}

	return resp, nil
}

// Get performs a GET request to the specified path
func (c *AuthenticatedClient) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	return c.Do(req)
}

// Post performs a POST request to the specified path with the given body
func (c *AuthenticatedClient) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

// Put performs a PUT request to the specified path with the given body
func (c *AuthenticatedClient) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

// Delete performs a DELETE request to the specified path
func (c *AuthenticatedClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	return c.Do(req)
}
