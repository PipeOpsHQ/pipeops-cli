package auth

import (
	"fmt"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "me",
	Short: "👤 Get details of the currently logged-in user",
	Long: `👤 The "me" command retrieves and displays details about the currently logged-in user 
in the PipeOps platform. Use this to confirm your active user identity and associated account.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Logic to fetch and display user details
		fmt.Println("Fetching details of the currently logged-in user...")
		// Example: Replace this with actual API call or logic
		fmt.Println("Username: johndoe")
		fmt.Println("Email: johndoe@example.com")
		fmt.Println("Role: Admin")
	},
	Args: cobra.NoArgs, // This command does not accept arguments
}

func (k *authModel) me() {
	k.rootCmd.AddCommand(authCmd)
}
