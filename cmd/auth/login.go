package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "🔐 Login to your PipeOps account",
	Long: `🔐 The "login" command allows you to authenticate with your PipeOps account
using your service account token. This will enable you to use other PipeOps CLI commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := pipeops.NewClient()

		// Load existing configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if already authenticated
		if client.IsAuthenticated() {
			fmt.Println("🔍 Checking existing authentication...")

			resp, err := client.VerifyToken()
			if err == nil && resp.Valid {
				fmt.Println("✅ You are already logged in!")
				config := client.GetConfig()
				if config.Username != "" {
					fmt.Printf("👤 Username: %s\n", config.Username)
				}
				if config.Email != "" {
					fmt.Printf("📧 Email: %s\n", config.Email)
				}
				return
			}

			fmt.Println("⚠️  Your current token is invalid or expired.")
		}

		// Get token from user
		fmt.Println("🔐 Please enter your PipeOps Service Account Token:")
		fmt.Print("Token: ")

		reader := bufio.NewReader(os.Stdin)
		token, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ Error reading token: %v\n", err)
			return
		}

		token = strings.TrimSpace(token)
		if token == "" {
			fmt.Println("❌ Token cannot be empty")
			return
		}

		// Validate token format
		if err := validation.ValidateToken(token); err != nil {
			fmt.Printf("❌ Invalid token format: %v\n", err)
			return
		}

		// Verify the token
		fmt.Println("🔍 Verifying token...")
		client.SetToken(token)

		resp, err := client.VerifyToken()
		if err != nil {
			fmt.Printf("❌ Error verifying token: %v\n", err)
			return
		}

		if !resp.Valid {
			fmt.Println("❌ Invalid token. Please check your token and try again.")
			return
		}

		// Save the configuration
		if err := client.SaveConfig(); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			return
		}

		fmt.Println("✅ Login successful!")
		fmt.Println("🎉 You can now use other PipeOps CLI commands.")
	},
	Args: cobra.NoArgs,
}

func (a *authModel) login() {
	a.rootCmd.AddCommand(loginCmd)
}
