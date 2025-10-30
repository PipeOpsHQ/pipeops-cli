/*
Copyright © 2024 9trocode

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/updater"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// Version is set at build time
var Version = "dev"

type Config struct {
	Version VersionInfo
	Updates UpdateSettings
}

type VersionInfo struct {
	Version string
}

type UpdateSettings struct {
	LastUpdateCheck time.Time `json:"last_update_check"`
	SkipUpdateCheck bool      `json:"skip_update_check"`
}

var Conf Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pipeops",
	Short:   "🚀 PipeOps CLI - Manage cloud-native development and deployment workflows",
	Long:    `🚀 PipeOps CLI is a command-line interface for managing cloud-native development and deployment workflows. Securely authenticate, manage projects and servers, deploy CI/CD pipelines, and monitor infrastructure—all from your terminal.`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set global JSON output flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			// Set a global flag that other commands can check
			cmd.Root().SetContext(context.WithValue(cmd.Root().Context(), "json", true))
		}

		// Check for updates periodically (but don't block the command)
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		go func() {
			if err := checkForUpdatesBackground(ctx, cmd); err != nil {
				// Log errors to stderr for debugging if verbose mode is enabled
				if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
					fmt.Fprintf(os.Stderr, "Update check warning: %v\n", err)
				}
			}
		}()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Changed {
			fmt.Println("🚀 PipeOps CLI Version:", Version)
			return
		}

		// Show help by default
		cmd.Help()
	},
}

// checkForUpdatesBackground runs a background update check
func checkForUpdatesBackground(ctx context.Context, cmd *cobra.Command) error {
	// Skip update check if specifically disabled
	if shouldSkipUpdateCheck(cmd) {
		return nil
	}

	// Check if it's been more than 24 hours since last check
	if !shouldCheckForUpdates() {
		return nil
	}

	// Get current version
	currentVersion := Version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	// Create update service
	updateService := updater.NewUpdateService(currentVersion)

	// Check for updates with a short timeout
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	release, hasUpdate, err := updateService.CheckForUpdates(checkCtx)
	if err != nil {
		return fmt.Errorf("update check failed: %w", err)
	}

	// Update last check time
	if err := updateLastCheckTime(); err != nil {
		return fmt.Errorf("failed to update check time: %w", err)
	}

	// If update available, show notification
	if hasUpdate {
		fmt.Printf("\n💡 A new version of PipeOps CLI is available: %s (current: %s)\n", release.TagName, currentVersion)
		fmt.Printf("   Run 'pipeops update' to install the latest version\n")
		fmt.Printf("   Run 'pipeops update check' to see what's new\n\n")
	}

	return nil
}

// shouldSkipUpdateCheck determines if update checking should be skipped
func shouldSkipUpdateCheck(cmd *cobra.Command) bool {
	// Skip for certain commands
	if cmd.Name() == "update" || cmd.Name() == "version" || cmd.Name() == "help" {
		return true
	}

	// Skip if running in CI/automated environment
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		return true
	}

	// Skip if explicitly disabled
	if os.Getenv("PIPEOPS_SKIP_UPDATE_CHECK") == "true" {
		return true
	}

	// Skip if JSON output is requested (likely automated)
	if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
		return true
	}

	return false
}

// shouldCheckForUpdates determines if it's time to check for updates
func shouldCheckForUpdates() bool {
	// Try to load existing config
	config := loadConfigSafely()

	// Check if enough time has passed since last check
	if time.Since(config.Updates.LastUpdateCheck) < 24*time.Hour {
		return false
	}

	return true
}

// loadConfigSafely loads config without exiting on errors
func loadConfigSafely() Config {
	var config Config

	home, err := os.UserHomeDir()
	if err != nil {
		if os.Getenv("PIPEOPS_DEBUG") == "true" {
			fmt.Fprintf(os.Stderr, "Warning: failed to get home directory: %v\n", err)
		}
		return config
	}

	filename := fmt.Sprintf("%s/.pipeops.json", home)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return config
	}

	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		if os.Getenv("PIPEOPS_DEBUG") == "true" {
			fmt.Fprintf(os.Stderr, "Warning: failed to read config: %v\n", err)
		}
		return config
	}

	if err := json.Unmarshal(dataBytes, &config); err != nil {
		if os.Getenv("PIPEOPS_DEBUG") == "true" {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse config: %v\n", err)
		}
	}
	return config
}

// updateLastCheckTime updates the last update check time
func updateLastCheckTime() error {
	config := loadConfigSafely()
	config.Updates.LastUpdateCheck = time.Now()

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	filename := fmt.Sprintf("%s/.pipeops.json", home)

	dataBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, dataBytes, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "Suppress non-essential output")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pipeops.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Prints out the current version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pipeops-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pipeops")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func GetConfig() (Config, error) {
	var filename string

	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get user home directory: %w", err)
	}

	filename = fmt.Sprintf("%s/%s", home, ".pipeops.json")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file does not exist: %s", filename)
	}

	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(dataBytes, &Conf)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return Conf, nil
}

func SaveConfig() error {
	var filename string

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	filename = fmt.Sprintf("%s/%s", home, ".pipeops.json")

	Conf.Version.Version = Version

	dataBytes, err := json.Marshal(Conf)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(filename, dataBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	if err := os.Chmod(filename, 0600); err != nil {
		return fmt.Errorf("failed to set config file permissions: %w", err)
	}

	return nil
}
