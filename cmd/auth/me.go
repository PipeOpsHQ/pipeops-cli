package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:     "me",
	Aliases: []string{"whoami"},
	Short:   "👤 Show current user information",
	Long: `👤 Display information about the currently authenticated user.

This command shows:
- User ID and username
- Email address and name
- Account creation date
- Authentication status

Examples:
  - Show user info:
    pipeops auth me

  - Show user info in JSON format:
    pipeops auth me --json

  - Alternative command:
    pipeops auth whoami`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Failed to load configuration", opts)
			return
		}

		// Check authentication
		if !cfg.IsAuthenticated() {
			utils.PrintError("Not authenticated. Please run 'pipeops auth login' to authenticate.", opts)
			return
		}

		// Create OAuth service and get user info
		authService := auth.NewOAuthService(cfg.OAuth)
		ctx := context.Background()

		utils.PrintInfo("Fetching user information...", opts)

		userInfo, err := authService.GetUserInfo(ctx)
		if err != nil {
			utils.HandleError(err, "Failed to get user information", opts)
			return
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(userInfo)
		} else {
			utils.PrintSuccess("User information retrieved successfully", opts)

			fmt.Printf("\n👤 USER INFORMATION\n")
			fmt.Printf("├─ ID: %d\n", userInfo.ID)
			fmt.Printf("├─ Username: %s\n", userInfo.Username)
			fmt.Printf("├─ Email: %s\n", userInfo.Email)
			fmt.Printf("├─ Name: %s %s\n", userInfo.FirstName, userInfo.LastName)
			fmt.Printf("├─ Created: %s\n", formatTime(userInfo.CreatedAt))
			fmt.Printf("├─ Updated: %s\n", formatTime(userInfo.UpdatedAt))
			fmt.Printf("├─ API Endpoint: %s\n", cfg.OAuth.BaseURL)
			fmt.Printf("└─ Token Status: %s Valid\n", utils.GetStatusIcon("success"))

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\n💡 TIPS\n")
				fmt.Printf("├─ List projects: pipeops list\n")
				fmt.Printf("├─ Create project: pipeops create <project-name>\n")
				fmt.Printf("├─ Check auth status: pipeops auth status\n")
				fmt.Printf("└─ Logout: pipeops auth logout\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) me() {
	k.rootCmd.AddCommand(meCmd)
}

// formatTime formats a time string for display
func formatTime(timeStr string) string {
	if timeStr == "" {
		return "N/A"
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}

	return t.Format("2006-01-02 15:04:05 MST")
}
