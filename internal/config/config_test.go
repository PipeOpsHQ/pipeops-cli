package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if cfg.OAuth == nil {
		t.Error("OAuth config is nil")
	}

	if cfg.Settings == nil {
		t.Error("Settings config is nil")
	}

	if cfg.OAuth.ClientID == "" {
		t.Error("ClientID is empty")
	}

	if cfg.Settings.OutputFormat == "" {
		t.Error("OutputFormat is empty")
	}
}

func TestGetClientID(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "returns environment variable when set",
			envValue: "test_client_id",
			want:     "test_client_id",
		},
		{
			name:     "returns default when env not set",
			envValue: "",
			want:     DefaultClientID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("PIPEOPS_CLIENT_ID", tt.envValue)
				defer os.Unsetenv("PIPEOPS_CLIENT_ID")
			} else {
				os.Unsetenv("PIPEOPS_CLIENT_ID")
			}

			got := GetClientID()
			if got != tt.want {
				t.Errorf("GetClientID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAPIURL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "returns environment variable when set",
			envValue: "https://custom-api.example.com",
			want:     "https://custom-api.example.com",
		},
		{
			name:     "returns default when env not set",
			envValue: "",
			want:     DefaultAPIURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("PIPEOPS_API_URL", tt.envValue)
				defer os.Unsetenv("PIPEOPS_API_URL")
			} else {
				os.Unsetenv("PIPEOPS_API_URL")
			}

			got := GetAPIURL()
			if got != tt.want {
				t.Errorf("GetAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ConfigFileName)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create test config
	cfg := DefaultConfig()
	cfg.OAuth.AccessToken = "test_token"
	cfg.OAuth.RefreshToken = "refresh_token"
	cfg.OAuth.ExpiresAt = time.Now().Add(1 * time.Hour)

	// Save config
	err := Save(cfg)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Config file permissions = %v, want 0600", info.Mode().Perm())
	}

	// Load config
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify loaded config
	if loadedCfg.OAuth.AccessToken != cfg.OAuth.AccessToken {
		t.Errorf("AccessToken = %v, want %v", loadedCfg.OAuth.AccessToken, cfg.OAuth.AccessToken)
	}

	if loadedCfg.OAuth.RefreshToken != cfg.OAuth.RefreshToken {
		t.Errorf("RefreshToken = %v, want %v", loadedCfg.OAuth.RefreshToken, cfg.OAuth.RefreshToken)
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Override home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Load config (should return default)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Should have default values
	if cfg.OAuth == nil || cfg.Settings == nil {
		t.Error("Load() should return default config when file doesn't exist")
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name        string
		accessToken string
		expiresAt   time.Time
		want        bool
	}{
		{
			name:        "valid token",
			accessToken: "valid_token",
			expiresAt:   time.Now().Add(1 * time.Hour),
			want:        true,
		},
		{
			name:        "expired token",
			accessToken: "expired_token",
			expiresAt:   time.Now().Add(-1 * time.Hour),
			want:        false,
		},
		{
			name:        "token expiring soon (within 5 min buffer)",
			accessToken: "expiring_token",
			expiresAt:   time.Now().Add(3 * time.Minute),
			want:        false,
		},
		{
			name:        "no token",
			accessToken: "",
			expiresAt:   time.Now().Add(1 * time.Hour),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				OAuth: &OAuthConfig{
					AccessToken: tt.accessToken,
					ExpiresAt:   tt.expiresAt,
				},
			}

			got := cfg.IsAuthenticated()
			if got != tt.want {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearAuth(t *testing.T) {
	cfg := &Config{
		OAuth: &OAuthConfig{
			AccessToken:  "test_token",
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
	}

	cfg.ClearAuth()

	if cfg.OAuth.AccessToken != "" {
		t.Error("AccessToken should be cleared")
	}

	if cfg.OAuth.RefreshToken != "" {
		t.Error("RefreshToken should be cleared")
	}

	if !cfg.OAuth.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should be zero time")
	}
}

func TestEnvironmentVariableOverrides(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Override home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save a config with default values
	cfg := DefaultConfig()
	cfg.OAuth.BaseURL = "https://default.example.com"
	cfg.OAuth.ClientID = "default_client"
	cfg.Settings.Debug = false

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Set environment variables
	os.Setenv("PIPEOPS_API_URL", "https://override.example.com")
	os.Setenv("PIPEOPS_CLIENT_ID", "override_client")
	os.Setenv("PIPEOPS_DEBUG", "true")
	defer func() {
		os.Unsetenv("PIPEOPS_API_URL")
		os.Unsetenv("PIPEOPS_CLIENT_ID")
		os.Unsetenv("PIPEOPS_DEBUG")
	}()

	// Load config - should use environment overrides
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loadedCfg.OAuth.BaseURL != "https://override.example.com" {
		t.Errorf("BaseURL = %v, want https://override.example.com", loadedCfg.OAuth.BaseURL)
	}

	if loadedCfg.OAuth.ClientID != "override_client" {
		t.Errorf("ClientID = %v, want override_client", loadedCfg.OAuth.ClientID)
	}

	if !loadedCfg.Settings.Debug {
		t.Error("Debug should be true from environment variable")
	}
}
