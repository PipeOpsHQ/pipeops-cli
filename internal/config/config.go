package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	ConfigFileName = ".pipeops.json"
	ConfigDirName  = ".pipeops"
)

// Config represents the CLI configuration
type Config struct {
	OAuth    *OAuthConfig `json:"oauth,omitempty"`
	Settings *Settings    `json:"settings,omitempty"`
}

// OAuthConfig holds OAuth-related configuration
type OAuthConfig struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	BaseURL      string    `json:"base_url"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	Scopes       []string  `json:"scopes,omitempty"`
}

// Settings holds general CLI settings
type Settings struct {
	DefaultRegion string `json:"default_region,omitempty"`
	OutputFormat  string `json:"output_format,omitempty"`
	Debug         bool   `json:"debug,omitempty"`
}

// DefaultConfig returns a new config with default values
func DefaultConfig() *Config {
	return &Config{
		OAuth: &OAuthConfig{
			ClientID: "pipeops_default_client", // Original client ID
			BaseURL:  "https://api.pipeops.sh",
			Scopes:   []string{"read:user", "read:projects", "write:projects"},
		},
		Settings: &Settings{
			OutputFormat: "table",
			Debug:        false,
		},
	}
}

// Load reads configuration from disk
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Return default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Ensure defaults for missing fields
	if cfg.OAuth == nil {
		cfg.OAuth = DefaultConfig().OAuth
	}
	if cfg.Settings == nil {
		cfg.Settings = DefaultConfig().Settings
	}

	// Override with environment variables if available
	// if apiURL := os.Getenv("PIPEOPS_API_URL"); apiURL != "" {
	// 	cfg.OAuth.BaseURL = apiURL
	// }
	// if clientID := os.Getenv("PIPEOPS_CLIENT_ID"); clientID != "" {
	// 	cfg.OAuth.ClientID = clientID
	// }
	// if debug := os.Getenv("PIPEOPS_DEBUG"); debug == "true" {
	// 	cfg.Settings.Debug = true
	// }

	return &cfg, nil
}

// Save writes configuration to disk with secure permissions
func Save(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with secure permissions (read/write for owner only)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the full path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ConfigFileName), nil
}

// IsAuthenticated checks if the user has valid authentication
func (c *Config) IsAuthenticated() bool {
	if c.OAuth == nil || c.OAuth.AccessToken == "" {
		return false
	}

	// Check if token is expired (with 5 minute buffer)
	return time.Now().Before(c.OAuth.ExpiresAt.Add(-5 * time.Minute))
}

// ClearAuth removes authentication information
func (c *Config) ClearAuth() {
	if c.OAuth != nil {
		c.OAuth.AccessToken = ""
		c.OAuth.RefreshToken = ""
		c.OAuth.ExpiresAt = time.Time{}
	}
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}
