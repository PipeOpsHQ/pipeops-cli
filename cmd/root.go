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
	"fmt"
	"os"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/updater"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cfgFile string

// Version is set at build time
var Version = "dev"

var Conf config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pipeops",
	Short: "PipeOps CLI - Manage cloud-native development and deployment workflows",
	Long:    `PipeOps CLI is a command-line interface for managing cloud-native development and deployment workflows. Securely authenticate, manage projects and servers, deploy CI/CD pipelines, and monitor infrastructure—all from your terminal.`,
	Version: Version,
	SilenceErrors: true, // We handle errors in main.go
	SilenceUsage:  true, // Don't show usage on error
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set global JSON output flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			// Set a global flag that other commands can check
			cmd.Root().SetContext(context.WithValue(cmd.Root().Context(), "json", true))
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Check for updates after command completes
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		checkForUpdatesAndPrompt(ctx, cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Changed {
			fmt.Println("PipeOps CLI Version:", Version)
			return
		}

		// Show help by default
		cmd.Help()
	},
}

// checkForUpdatesAndPrompt checks for updates and prompts user to update
func checkForUpdatesAndPrompt(ctx context.Context, cmd *cobra.Command) {
	// Skip update check if specifically disabled
	if shouldSkipUpdateCheck(cmd) {
		return
	}

	// Check if it's been more than 24 hours since last check
	if !shouldCheckForUpdates() {
		return
	}

	// Get current version
	currentVersion := Version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	// Create update service
	updateService := updater.NewUpdateService(currentVersion)

	// Check for updates with a short timeout
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	release, hasUpdate, err := updateService.CheckForUpdates(checkCtx)
	if err != nil {
		// Silently fail - don't interrupt user
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintf(os.Stderr, "Update check warning: %v\n", err)
		}
		return
	}

	// Update last check time
	_ = updateLastCheckTime()

	// If update available, show notification and prompt
	if hasUpdate {
		fmt.Printf("\n")
		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		fmt.Printf("%s A new version of PipeOps CLI is available: %s (current: %s)\n", yellow("📦"), release.TagName, currentVersion)
		fmt.Printf("\n")

		// Ask user if they want to update now
		fmt.Printf("Would you like to update now? [y/N]: ")
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" || response == "yes" || response == "Yes" {
			fmt.Printf("\n")
			opts := utils.OutputOptions{Format: utils.OutputFormatTable}
			if err := updateService.UpdateCLI(ctx, release, opts); err != nil {
				fmt.Printf("❌ Update failed: %v\n", err)
				fmt.Printf("   You can manually update by running: pipeops update\n")
			}
		} else {
			fmt.Printf("\n💡 Run 'pipeops update' when you're ready to update.\n")
		}
	}
}

// checkForUpdatesBackground runs a background update check (deprecated, kept for reference)
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
		fmt.Printf("\nA new version of PipeOps CLI is available: %s (current: %s)\n", release.TagName, currentVersion)
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
	cfg, err := config.Load()
	if err != nil {
		return true
	}

	if cfg.Updates == nil {
		return true
	}

	// Check if enough time has passed since last check
	if time.Since(cfg.Updates.LastUpdateCheck) < 24*time.Hour {
		return false
	}

	return true
}

// updateLastCheckTime updates the last update check time
func updateLastCheckTime() error {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	if cfg.Updates == nil {
		cfg.Updates = &config.UpdateSettings{}
	}
	cfg.Updates.LastUpdateCheck = time.Now()

	return config.Save(cfg)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "Suppress non-essential output")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pipeops.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Prints out the current version")

	// Custom Help Template
	bold := color.New(color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	rootCmd.SetHelpTemplate(fmt.Sprintf(`%s
  {{.Long}}

%s
  {{.UseLine}}

%s
{{if .HasAvailableSubCommands}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  %s {{.Short}}
{{end}}{{end}}{{end}}
%s
{{if .HasAvailableFlags}}{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}{{end}}

%s
  %s

`,
		bold("PipeOps CLI"),
		bold("USAGE"),
		bold("AVAILABLE COMMANDS"),
		green("{{rpad .Name .NamePadding }}"),
		bold("FLAGS"),
		bold("LEARN MORE"),
		cyan("Use 'pipeops [command] --help' for more information about a command."),
	))

	// Custom Usage Template
	rootCmd.SetUsageTemplate(fmt.Sprintf(`%s
  {{.UseLine}}

%s
{{if .HasAvailableSubCommands}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  %s {{.Short}}
{{end}}{{end}}{{end}}
%s
{{if .HasAvailableFlags}}{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}{{end}}
`,
		bold("USAGE"),
		bold("AVAILABLE COMMANDS"),
		yellow("{{rpad .Name .NamePadding }}"),
		bold("FLAGS"),
	))
}
