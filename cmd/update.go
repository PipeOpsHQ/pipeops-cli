package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/updater"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install CLI updates",
	Long: `Check for and install CLI updates.

The update command allows you to:
- Check for newer versions of PipeOps CLI
- Install the latest version with your consent
- View release notes and changes

Examples:
  pipeops update              # Check for updates and prompt to install
  pipeops update check        # Just check for updates without installing
  pipeops update --yes        # Install updates without prompting
  pipeops update --json       # Get update information in JSON format`,
	Run: func(cmd *cobra.Command, args []string) {
		runUpdateCheck(cmd, args, true) // Default behavior includes installation
	},
}

// updateCheckCmd represents the update check command
var updateCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available updates",
	Long: `Check for available updates without installing.

This command will check if a newer version of PipeOps CLI is available
and display information about it, but won't install anything.

Examples:
  pipeops update check        # Check for updates
  pipeops update check --json # Get update info in JSON format`,
	Run: func(cmd *cobra.Command, args []string) {
		runUpdateCheck(cmd, args, false) // Check only, no installation
	},
}

// runUpdateCheck handles the update checking logic
func runUpdateCheck(cmd *cobra.Command, args []string, allowInstall bool) {
	opts := utils.GetOutputOptions(cmd)

	// Get current version
	currentVersion := Version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	// Create update service
	updateService := updater.NewUpdateService(currentVersion)

	// Check for updates
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	utils.PrintInfo("Checking for updates...", opts)

	release, hasUpdate, err := updateService.CheckForUpdates(ctx)
	if err != nil {
		utils.HandleError(err, "Failed to check for updates", opts)
		return
	}

	// Handle no updates available
	if !hasUpdate {
		if opts.Format == utils.OutputFormatJSON {
			result := map[string]interface{}{
				"current_version": currentVersion,
				"latest_version":  release.TagName,
				"up_to_date":      true,
				"message":         "You are using the latest version",
			}
			utils.PrintJSON(result)
		} else {
			utils.PrintSuccess(fmt.Sprintf("You are using the latest version (%s)", currentVersion), opts)
		}
		return
	}

	// Handle updates available
	if opts.Format == utils.OutputFormatJSON {
		result := map[string]interface{}{
			"current_version":  currentVersion,
			"latest_version":   release.TagName,
			"up_to_date":       false,
			"update_available": true,
			"release_name":     release.Name,
			"release_date":     release.PublishedAt.Format("2006-01-02"),
			"release_notes":    release.Body,
			"download_url":     fmt.Sprintf("https://github.com/%s/releases/tag/%s", updater.GetGitHubRepo(), release.TagName),
		}
		utils.PrintJSON(result)
		return
	}

	// Display update information in human-readable format
	fmt.Printf("\nA new version of PipeOps CLI is available!\n")
	fmt.Printf("   Current version: %s\n", currentVersion)
	fmt.Printf("   Latest version:  %s\n", release.TagName)
	fmt.Printf("   Release date:    %s\n", release.PublishedAt.Format("January 2, 2006"))

	if release.Name != "" && release.Name != release.TagName {
		fmt.Printf("   Release name:    %s\n", release.Name)
	}

	// Show release notes if available
	if release.Body != "" {
		fmt.Printf("\nRelease Notes:\n")

		// Format release notes nicely
		lines := strings.Split(release.Body, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("   %s\n", line)
			}
		}
	}

	fmt.Printf("\nView full release: https://github.com/%s/releases/tag/%s\n",
		updater.GetGitHubRepo(), release.TagName)

	// If not allowing installation, just show the information
	if !allowInstall {
		fmt.Printf("\nTo install the update, run: pipeops update\n")
		return
	}

	// Ask for user consent to update
	skipPrompt, _ := cmd.Flags().GetBool("yes")
	if !skipPrompt {
		fmt.Printf("\n")
		if !utils.ConfirmAction("Would you like to update now?") {
			fmt.Println("Update cancelled. You can update later by running 'pipeops update'")
			return
		}
	}

	// Perform the update
	fmt.Printf("\nStarting update process...\n")
	if err := updateService.UpdateCLI(ctx, release, opts); err != nil {
		utils.HandleError(err, "Failed to update CLI", opts)
		return
	}

	fmt.Printf("\n[OK] Update completed successfully!\n")
	fmt.Printf("   You are now running version %s\n", release.TagName)
	fmt.Printf("\n[INFO] You may need to restart your terminal or shell to use the updated version.\n")
}

func init() {
	// Add check subcommand
	updateCmd.AddCommand(updateCheckCmd)

	// Add flags
	updateCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt and install updates automatically")
	updateCmd.Flags().Bool("check-only", false, "Only check for updates without installing (same as 'update check')")

	// Add to root command
	rootCmd.AddCommand(updateCmd)
}
