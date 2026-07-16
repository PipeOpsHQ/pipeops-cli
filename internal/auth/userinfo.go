package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
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
}

// NewUserInfoService creates a new userinfo service
func NewUserInfoService(cfg *config.Config) *UserInfoService {
	return &UserInfoService{
		config: cfg,
	}
}

// GetUserInfo fetches user information from the OAuth userinfo endpoint
func (s *UserInfoService) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := s.newSDKClient(accessToken)
	if err != nil {
		return nil, err
	}

	profileResp, _, err := client.Users.GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}

	return userFromSDK(profileResp.Data.User), nil
}

func (s *UserInfoService) newSDKClient(accessToken string) (*sdk.Client, error) {
	baseURL := config.GetAPIURL()
	if s.config != nil && s.config.OAuth != nil && s.config.OAuth.BaseURL != "" {
		baseURL = s.config.OAuth.BaseURL
	}
	client, err := sdk.NewClient(baseURL, sdk.WithTimeout(30*time.Second), sdk.WithMaxRetries(3))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PipeOps SDK client: %w", err)
	}
	client.SetToken(accessToken)
	return client, nil
}

func userFromSDK(user sdk.User) *UserInfo {
	id, _ := strconv.Atoi(user.ID)
	info := &UserInfo{
		ID:                 id,
		UUID:               config.SanitizeLog(user.UUID),
		Email:              config.SanitizeLog(user.Email),
		Name:               config.SanitizeLog(user.FullName),
		FirstName:          config.SanitizeLog(user.FirstName),
		LastName:           config.SanitizeLog(user.LastName),
		Avatar:             config.SanitizeLog(user.AvatarURL),
		Verified:           user.EmailVerified,
		SubscriptionActive: user.IsSubscriptionActive,
	}
	if user.CreatedAt != nil {
		info.CreatedAt = user.CreatedAt.Time
	}
	if user.UpdatedAt != nil {
		info.UpdatedAt = user.UpdatedAt.Time
	}
	return info
}

// FormatUserInfo formats user information for display
func (ui *UserInfo) FormatUserInfo() string {
	// Use the best available name for the header
	displayName := config.SanitizeLog(ui.GetDisplayName())
	output := fmt.Sprintf("👤 %s", displayName)

	// Add username if different from display name
	if ui.Username != "" && ui.Username != displayName {
		output += fmt.Sprintf(" (@%s)", config.SanitizeLog(ui.Username))
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
		output += fmt.Sprintf("📧 %s%s\n", config.SanitizeLog(ui.Email), verified)
	}

	// User ID and UUID
	if ui.ID != 0 {
		output += fmt.Sprintf("🆔 %s", ui.GetIDString())
		if ui.UUID != "" {
			output += fmt.Sprintf(" (UUID: %s)", config.SanitizeLog(ui.UUID))
		}
		output += "\n"
	}

	// Avatar
	if ui.Avatar != "" {
		output += fmt.Sprintf("🖼️  Avatar: %s\n", config.SanitizeLog(ui.Avatar))
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
		output += fmt.Sprintf("🏷️  Roles: %v\n", sanitizeStringSlice(ui.Roles))
	}

	if len(ui.Permissions) > 0 {
		output += fmt.Sprintf("🔑 Permissions: %v\n", sanitizeStringSlice(ui.Permissions))
	}

	return output
}

func sanitizeStringSlice(values []string) []string {
	sanitized := make([]string, len(values))
	for i, value := range values {
		sanitized[i] = config.SanitizeLog(value)
	}
	return sanitized
}
