package cmd

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
  pipeops logout
  pipeops logout --json
  pipeops logout --force`,
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
				fmt.Println("✅ You're already logged out")
				fmt.Println("🔑 When ready to return: pipeops auth login")
			}
			return
		}

		// Confirm logout unless force flag is used
		force, _ := cmd.Flags().GetBool("force")
		if !force && opts.Format != utils.OutputFormatJSON {
			if !utils.ConfirmAction("🚪 Are you sure you want to log out?") {
				fmt.Println("✅ Staying logged in")
				fmt.Println("🚀 Continue using PipeOps: pipeops project list")
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
			fmt.Println("👋 Successfully logged out!")
			fmt.Println("🔑 To log back in: pipeops auth login")
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
	logoutCmd.Flags().BoolP("force", "f", false, "Force logout without confirmation")
}
