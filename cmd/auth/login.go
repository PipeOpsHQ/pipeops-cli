package auth

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "🔐 Login to your PipeOps account",
	Long: `🔐 Login to your PipeOps account using your authentication token.

Examples:
  - Login interactively:
    pipeops auth login

  - Login with token:
    pipeops auth login --token <your-token>

  - Login with JSON output:
    pipeops auth login --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Get token from flag or prompt
		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintError("Token is required for JSON output", opts)
				return
			}

			// Interactive mode
			fmt.Print("🔑 Enter your PipeOps token: ")
			tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				utils.HandleError(err, "Error reading token", opts)
				return
			}
			token = strings.TrimSpace(string(tokenBytes))
			fmt.Println() // Add newline after password input
		}

		// Validate token
		if err := validation.ValidateToken(token); err != nil {
			utils.PrintError(fmt.Sprintf("Invalid token: %v", err), opts)
			return
		}

		// Verify token with API
		utils.PrintInfo("Verifying token...", opts)

		client.SetToken(token)
		resp, err := client.VerifyToken()
		if err != nil {
			utils.HandleError(err, "Error verifying token", opts)
			return
		}

		if !resp.Valid {
			utils.PrintError("Invalid token. Please check your token and try again.", opts)
			return
		}

		// Save configuration
		if err := client.SaveConfig(); err != nil {
			utils.HandleError(err, "Error saving configuration", opts)
			return
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			result := map[string]interface{}{
				"success":     true,
				"token_valid": resp.Valid,
			}
			utils.PrintJSON(result)
		} else {
			utils.PrintSuccess("Login successful!", opts)

			fmt.Printf("\n🎉 You're now logged in to PipeOps!\n")

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\n💡 NEXT STEPS\n")
				fmt.Printf("├─ View user info: pipeops auth me\n")
				fmt.Printf("├─ List projects: pipeops list\n")
				fmt.Printf("├─ Create project: pipeops create <project-name>\n")
				fmt.Printf("└─ Get help: pipeops --help\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

func (k *authModel) login() {
	k.rootCmd.AddCommand(loginCmd)

	// Add flags
	loginCmd.Flags().StringP("token", "t", "", "Authentication token")
}
