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
	Short: "Manage your PipeOps authentication",
	Long: `Manage your PipeOps authentication and account access.

Common commands:
  pipeops auth login     Log in to your PipeOps account
  pipeops auth me        View your profile information
  pipeops auth status    Check your authentication status
  pipeops auth logout    Log out of your account

Get started by running: pipeops auth login`,
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
