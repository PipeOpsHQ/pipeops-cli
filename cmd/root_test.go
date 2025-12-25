package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestShouldSkipUpdateCheck(t *testing.T) {
	// Save original environment variables
	origCI := os.Getenv("CI")
	origGithubActions := os.Getenv("GITHUB_ACTIONS")
	origSkip := os.Getenv("PIPEOPS_SKIP_UPDATE_CHECK")

	// Restore them after test
	defer func() {
		if origCI != "" {
			os.Setenv("CI", origCI)
		} else {
			os.Unsetenv("CI")
		}
		if origGithubActions != "" {
			os.Setenv("GITHUB_ACTIONS", origGithubActions)
		} else {
			os.Unsetenv("GITHUB_ACTIONS")
		}
		if origSkip != "" {
			os.Setenv("PIPEOPS_SKIP_UPDATE_CHECK", origSkip)
		} else {
			os.Unsetenv("PIPEOPS_SKIP_UPDATE_CHECK")
		}
	}()

	// Clear environment variables for test isolation
	os.Unsetenv("CI")
	os.Unsetenv("GITHUB_ACTIONS")
	os.Unsetenv("PIPEOPS_SKIP_UPDATE_CHECK")

	tests := []struct {
		name     string
		cmdName  string
		envVars  map[string]string
		jsonFlag bool
		want     bool
	}{
		{
			name:    "skip for update command",
			cmdName: "update",
			want:    true,
		},
		{
			name:    "skip for version command",
			cmdName: "version",
			want:    true,
		},
		{
			name:    "skip for help command",
			cmdName: "help",
			want:    true,
		},
		{
			name:    "skip in CI environment",
			cmdName: "project",
			envVars: map[string]string{"CI": "true"},
			want:    true,
		},
		{
			name:    "skip in GitHub Actions",
			cmdName: "project",
			envVars: map[string]string{"GITHUB_ACTIONS": "true"},
			want:    true,
		},
		{
			name:    "skip when explicitly disabled",
			cmdName: "project",
			envVars: map[string]string{"PIPEOPS_SKIP_UPDATE_CHECK": "true"},
			want:    true,
		},

		{
			name:    "don't skip for normal command",
			cmdName: "project",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Create a test command
			testCmd := &cobra.Command{
				Use: tt.cmdName,
			}
			testCmd.PersistentFlags().Bool("json", tt.jsonFlag, "JSON output")
			if tt.jsonFlag {
				testCmd.Flags().Set("json", "true")
			}

			got := shouldSkipUpdateCheck(testCmd)
			if got != tt.want {
				t.Errorf("shouldSkipUpdateCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionInfo(t *testing.T) {
	if Version == "" {
		// Version should be set via ldflags, but may be empty in tests
		t.Log("Version is empty (expected in tests)")
	}

	// Test that version command is available
	cmd := rootCmd
	if cmd.Version == "" && Version != "" {
		t.Error("rootCmd.Version should be set when Version is set")
	}
}

func TestRootCommandFlags(t *testing.T) {
	// Check that persistent flags are registered
	flags := []string{"json", "verbose", "quiet", "config"}

	for _, flagName := range flags {
		flag := rootCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected persistent flag %q not found", flagName)
		}
	}
}

func TestConfigFunctions(t *testing.T) {
	// Test that GetConfig and SaveConfig return errors properly
	// Note: These tests require a valid config file setup

	// Create temporary directory
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	// Test GetConfig with non-existent file
	_, err := GetConfig()
	if err == nil {
		t.Error("GetConfig() should return error when config file doesn't exist")
	}

	// Test SaveConfig
	Conf = Config{
		Version: VersionInfo{Version: "test"},
	}
	err = SaveConfig()
	if err != nil {
		t.Errorf("SaveConfig() unexpected error: %v", err)
	}

	// Now GetConfig should work
	_, err = GetConfig()
	if err != nil {
		t.Errorf("GetConfig() unexpected error after save: %v", err)
	}
}
