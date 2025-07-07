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
	Short: "Logout from your PipeOps account",
	Long: `Logout from your PipeOps account and clear stored authentication tokens.

Examples:
  pipeops auth logout
  pipeops auth logout --json
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
				fmt.Println("Already logged out")
			}
			return
		}

		// Confirm logout unless force flag is used
		force, _ := cmd.Flags().GetBool("force")
		if !force && opts.Format != utils.OutputFormatJSON {
			if !utils.ConfirmAction("Are you sure you want to logout?") {
				fmt.Println("Logout cancelled")
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
			fmt.Println("Logged out successfully")
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) logout() {
	k.rootCmd.AddCommand(logoutCmd)

	// Add flags
	logoutCmd.Flags().BoolP("force", "f", false, "Force logout without confirmation")
}
