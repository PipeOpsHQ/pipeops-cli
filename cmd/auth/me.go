package auth

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "👤 Show current user information",
	Long: `👤 Display information about the currently authenticated user.

Examples:
  - Show user info:
    pipeops auth me

  - Show user info in JSON format:
    pipeops auth me --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Get user information
		utils.PrintInfo("Fetching user information...", opts)

		// Verify token to ensure it's valid
		resp, err := client.VerifyToken()
		if err != nil {
			utils.HandleError(err, "Error verifying token", opts)
			return
		}

		if !resp.Valid {
			utils.PrintError("Invalid token. Please run 'pipeops auth login' to re-authenticate.", opts)
			return
		}

		// Get config for user details
		config := client.GetConfig()

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			userInfo := map[string]interface{}{
				"token_valid":  resp.Valid,
				"api_base_url": config.APIBaseURL,
			}
			if config.UserID != "" {
				userInfo["user_id"] = config.UserID
			}
			if config.Username != "" {
				userInfo["username"] = config.Username
			}
			if config.Email != "" {
				userInfo["email"] = config.Email
			}
			utils.PrintJSON(userInfo)
		} else {
			utils.PrintSuccess("User information retrieved successfully", opts)

			fmt.Printf("\n👤 USER INFORMATION\n")
			if config.UserID != "" {
				fmt.Printf("├─ User ID: %s\n", config.UserID)
			}
			if config.Username != "" {
				fmt.Printf("├─ Username: %s\n", config.Username)
			}
			if config.Email != "" {
				fmt.Printf("├─ Email: %s\n", config.Email)
			}
			fmt.Printf("├─ API Base URL: %s\n", config.APIBaseURL)
			fmt.Printf("└─ Token Status: %s Valid\n", utils.GetStatusIcon("success"))

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\n💡 TIPS\n")
				fmt.Printf("├─ List projects: pipeops list\n")
				fmt.Printf("├─ Create project: pipeops create <project-name>\n")
				fmt.Printf("└─ Logout: pipeops auth logout\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) me() {
	k.rootCmd.AddCommand(meCmd)
}
