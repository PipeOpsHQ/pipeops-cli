package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
)

// UserInfo represents the user information returned by the OAuth userinfo endpoint
type UserInfo struct {
	ID                 int       `json:"id"`       // Server returns number
	UUID               string    `json:"uuid"`     // Server returns UUID
	Username           string    `json:"username"` // May not be in response
	Email              string    `json:"email"`
	Name               string    `json:"name"`                // May not be in response
	FirstName          string    `json:"first_name"`          // Server field
	LastName           string    `json:"last_name"`           // Server field
	Avatar             string    `json:"avatar"`              // Server field
	Verified           bool      `json:"email_verified"`      // Server field (email_verified)
	SubscriptionActive bool      `json:"subscription_active"` // Server field
	CreatedAt          time.Time `json:"created_at"`          // May not be in response
	UpdatedAt          time.Time `json:"updated_at"`          // May not be in response
	LastLoginAt        time.Time `json:"last_login_at"`       // May not be in response
	Roles              []string  `json:"roles"`               // May not be in response
	Permissions        []string  `json:"permissions"`         // May not be in response
}

// GetIDString returns the ID as a string for compatibility
func (ui *UserInfo) GetIDString() string {
	return fmt.Sprintf("%d", ui.ID)
}

// GetFullName returns the full name from first and last name
func (ui *UserInfo) GetFullName() string {
	if ui.FirstName != "" && ui.LastName != "" {
		return ui.FirstName + " " + ui.LastName
	}
	if ui.Name != "" {
		return ui.Name
	}
	if ui.FirstName != "" {
		return ui.FirstName
	}
	if ui.LastName != "" {
		return ui.LastName
	}
	return ui.Username
}

// GetDisplayName returns the best available name for display
func (ui *UserInfo) GetDisplayName() string {
	if fullName := ui.GetFullName(); fullName != "" && fullName != ui.Username {
		return fullName
	}
	if ui.Username != "" {
		return ui.Username
	}
	return ui.Email
}

// UserInfoService handles OAuth userinfo API calls
type UserInfoService struct {
	config *config.Config
	client *http.Client
}

// NewUserInfoService creates a new userinfo service
func NewUserInfoService(cfg *config.Config) *UserInfoService {
	return &UserInfoService{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetUserInfo fetches user information from the OAuth userinfo endpoint
func (s *UserInfoService) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	// First attempt with Bearer token (standard OAuth 2.0)
	userInfo, err := s.getUserInfoWithBearer(ctx, accessToken)
	if err == nil {
		return userInfo, nil
	}

	// If Bearer token fails due to JWT format issues, try alternative approaches
	if s.isJWTFormatError(err) {
		// Check if token is actually a JWT that needs different handling
		if s.isJWT(accessToken) {
			return s.getUserInfoWithJWT(ctx, accessToken)
		}

		// For opaque tokens, provide helpful error message
		return nil, fmt.Errorf("server expects JWT tokens but received opaque token - this is a server configuration issue. Server response: %w", err)
	}

	return nil, err
}

// getUserInfoWithBearer makes a standard OAuth 2.0 Bearer token request
func (s *UserInfoService) getUserInfoWithBearer(ctx context.Context, accessToken string) (*UserInfo, error) {
	// Create request to userinfo endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", s.config.OAuth.BaseURL+"/oauth/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PipeOps-CLI/1.0")

	// Debug information
	if s.config.Settings != nil && s.config.Settings.Debug {
		fmt.Printf("🔍 Debug: Making Bearer token request to %s\n", config.SanitizeLog(req.URL.String()))
		tokenPreview := accessToken
		if len(accessToken) > 20 {
			tokenPreview = accessToken[:20] + "..."
		}
		fmt.Printf("🔍 Debug: Token preview: %s\n", tokenPreview)
		fmt.Printf("🔍 Debug: Token format: %s\n", s.getTokenFormat(accessToken))
	}

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for better error messages
	bodyBytes, readErr := io.ReadAll(resp.Body)
	var responseBody string
	if readErr == nil {
		responseBody = string(bodyBytes)
	}

	// Debug response information
	if s.config.Settings != nil && s.config.Settings.Debug {
		fmt.Printf("🔍 Debug: Response status: %d\n", resp.StatusCode)
		fmt.Printf("🔍 Debug: Response headers: %v\n", resp.Header)
		fmt.Printf("🔍 Debug: Response body: %s\n", config.SanitizeLog(responseBody))
	}

	// Check response status with detailed error messages
	if resp.StatusCode == http.StatusUnauthorized {
		if responseBody != "" {
			return nil, fmt.Errorf("authentication failed - server response: %s", responseBody)
		}
		return nil, fmt.Errorf("authentication expired or invalid - please run 'pipeops auth login'")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("userinfo endpoint not found - the API might not support this endpoint yet")
	}

	if resp.StatusCode != http.StatusOK {
		if responseBody != "" {
			return nil, fmt.Errorf("userinfo request failed with status %d: %s", resp.StatusCode, responseBody)
		}
		return nil, fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	// Parse response
	var userInfo UserInfo
	if readErr != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", readErr)
	}

	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w (response: %s)", err, responseBody)
	}

	return &userInfo, nil
}

// getUserInfoWithJWT handles JWT-specific authentication
func (s *UserInfoService) getUserInfoWithJWT(ctx context.Context, jwtToken string) (*UserInfo, error) {
	// This method can be implemented if needed for JWT-specific handling
	// For now, it attempts the same Bearer token approach
	return s.getUserInfoWithBearer(ctx, jwtToken)
}

// isJWTFormatError checks if the error is related to JWT format issues
func (s *UserInfoService) isJWTFormatError(err error) bool {
	errorStr := err.Error()
	return strings.Contains(errorStr, "invalid number of segments") ||
		strings.Contains(errorStr, "token contains an invalid") ||
		strings.Contains(errorStr, "malformed token")
}

// isJWT checks if a token is in JWT format (has 3 segments)
func (s *UserInfoService) isJWT(token string) bool {
	return strings.Count(token, ".") == 2
}

// getTokenFormat returns a human-readable description of the token format
func (s *UserInfoService) getTokenFormat(token string) string {
	segments := strings.Count(token, ".")
	if segments == 2 {
		return "JWT (3 segments)"
	}
	return fmt.Sprintf("Opaque (%d segments)", segments+1)
}

// FormatUserInfo formats user information for display
func (ui *UserInfo) FormatUserInfo() string {
	// Use the best available name for the header
	displayName := ui.GetDisplayName()
	output := fmt.Sprintf("👤 %s", displayName)

	// Add username if different from display name
	if ui.Username != "" && ui.Username != displayName {
		output += fmt.Sprintf(" (@%s)", ui.Username)
	}
	output += "\n"

	// Email with verification status
	if ui.Email != "" {
		verified := ""
		if ui.Verified {
			verified = " ✅"
		} else {
			verified = " ⚠️"
		}
		output += fmt.Sprintf("📧 %s%s\n", ui.Email, verified)
	}

	// User ID and UUID
	if ui.ID != 0 {
		output += fmt.Sprintf("🆔 %s", ui.GetIDString())
		if ui.UUID != "" {
			output += fmt.Sprintf(" (UUID: %s)", ui.UUID)
		}
		output += "\n"
	}

	// Avatar
	if ui.Avatar != "" {
		output += fmt.Sprintf("🖼️  Avatar: %s\n", ui.Avatar)
	}

	// Subscription status
	if ui.SubscriptionActive {
		output += "💎 Subscription: Active\n"
	} else {
		output += "💡 Subscription: Inactive\n"
	}

	// Account dates (if available)
	if !ui.CreatedAt.IsZero() {
		output += fmt.Sprintf("📅 Member since: %s\n", ui.CreatedAt.Format("January 2, 2006"))
	}

	if !ui.LastLoginAt.IsZero() {
		output += fmt.Sprintf("🔄 Last login: %s\n", ui.LastLoginAt.Format("January 2, 2006 15:04 MST"))
	}

	// Roles and permissions (if available)
	if len(ui.Roles) > 0 {
		output += fmt.Sprintf("🏷️  Roles: %v\n", ui.Roles)
	}

	if len(ui.Permissions) > 0 {
		output += fmt.Sprintf("🔑 Permissions: %v\n", ui.Permissions)
	}

	return output
}
