package auth

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "🚪 Logout from your PipeOps account",
	Long: `🚪 Logout from your PipeOps account and clear stored authentication tokens.

This command will:
1. Remove stored access and refresh tokens
2. Clear cached user information
3. Require re-authentication for future API calls

Examples:
  - Logout:
    pipeops auth logout

  - Logout with JSON output:
    pipeops auth logout --json

  - Force logout without confirmation:
    pipeops auth logout --force`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Failed to load configuration", opts)
			return
		}

		// Check if user is authenticated
		if !cfg.IsAuthenticated() {
			if opts.Format == utils.OutputFormatJSON {
				result := map[string]interface{}{
					"success": true,
					"message": "Already logged out",
				}
				utils.PrintJSON(result)
			} else {
				utils.PrintWarning("You are not currently logged in.", opts)
			}
			return
		}

		// Confirm logout unless force flag is used
		force, _ := cmd.Flags().GetBool("force")
		if !force && opts.Format != utils.OutputFormatJSON {
			if !utils.ConfirmAction("🔐 Are you sure you want to logout?") {
				utils.PrintInfo("Logout cancelled.", opts)
				return
			}
		}

		// Clear authentication
		cfg.ClearAuth()
		if err := config.Save(cfg); err != nil {
			utils.HandleError(err, "Failed to save configuration", opts)
			return
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			result := map[string]interface{}{
				"success": true,
				"message": "Successfully logged out",
			}
			utils.PrintJSON(result)
		} else {
			utils.PrintSuccess("Successfully logged out!", opts)

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\n💡 NEXT STEPS\n")
				fmt.Printf("├─ Login again: pipeops auth login\n")
				fmt.Printf("└─ Get help: pipeops --help\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) logout() {
	k.rootCmd.AddCommand(logoutCmd)

	// Add flags
	logoutCmd.Flags().BoolP("force", "f", false, "Force logout without confirmation")
}
