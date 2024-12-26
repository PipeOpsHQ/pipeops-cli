// config/config.go
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the PipeOps CLI configuration.
type Config struct {
	ServiceAccountToken string `json:"service_account_token"`
}

// getConfigPath returns the path to the configuration file.
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine user config directory: %w", err)
	}
	pipeopsConfigDir := filepath.Join(configDir, "pipeops-cli")
	if _, err := os.Stat(pipeopsConfigDir); os.IsNotExist(err) {
		if err := os.MkdirAll(pipeopsConfigDir, 0700); err != nil {
			return "", fmt.Errorf("unable to create config directory: %w", err)
		}
	}
	return filepath.Join(pipeopsConfigDir, "config.json"), nil
}

// SaveConfig saves the configuration to a file.
func SaveConfig(cfg Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// LoadConfig loads the configuration from a file.
func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	fileData, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("config file does not exist, please run 'pipeops install' first")
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(fileData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if cfg.ServiceAccountToken == "" {
		return nil, errors.New("service account token is missing in config")
	}
	return &cfg, nil
}
