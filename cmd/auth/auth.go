package auth

import (
	"github.com/spf13/cobra"
)

// authModel represents the auth command model
type authModel struct {
	rootCmd *cobra.Command
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  `Manage authentication and user account operations.`,
}

// New initializes and returns auth command
func New() *cobra.Command {
	authModel := &authModel{
		rootCmd: authCmd,
	}

	authModel.login()
	authModel.logout()
	authModel.status()
	authModel.me()
	authModel.debug()
	authModel.consent()

	return authCmd
}
