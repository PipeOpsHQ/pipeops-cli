package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
)

// UserInfo represents the user information returned by the OAuth userinfo endpoint
type UserInfo struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Avatar      string    `json:"avatar"`
	Verified    bool      `json:"verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastLoginAt time.Time `json:"last_login_at"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
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
	// Create request to userinfo endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", s.config.OAuth.BaseURL+"/oauth/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PipeOps-CLI/1.0")

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication expired - please run 'pipeops auth login'")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	// Parse response
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	return &userInfo, nil
}

// FormatUserInfo formats user information for display
func (ui *UserInfo) FormatUserInfo() string {
	output := fmt.Sprintf("👤 %s", ui.Name)
	if ui.Username != "" {
		output += fmt.Sprintf(" (@%s)", ui.Username)
	}
	output += "\n"

	if ui.Email != "" {
		verified := ""
		if ui.Verified {
			verified = " ✅"
		}
		output += fmt.Sprintf("📧 %s%s\n", ui.Email, verified)
	}

	if ui.ID != "" {
		output += fmt.Sprintf("🆔 %s\n", ui.ID)
	}

	if !ui.CreatedAt.IsZero() {
		output += fmt.Sprintf("📅 Member since: %s\n", ui.CreatedAt.Format("January 2, 2006"))
	}

	if !ui.LastLoginAt.IsZero() {
		output += fmt.Sprintf("🔄 Last login: %s\n", ui.LastLoginAt.Format("January 2, 2006 15:04 MST"))
	}

	if len(ui.Roles) > 0 {
		output += fmt.Sprintf("🏷️  Roles: %v\n", ui.Roles)
	}

	if len(ui.Permissions) > 0 {
		output += fmt.Sprintf("🔑 Permissions: %v\n", ui.Permissions)
	}

	return output
}
