package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "🔒 Manage authentication and access control.",
	Long: `🔒 The auth command provides a set of subcommands for managing 
authentication-related operations, such as login, logout, and credential management 
for interacting with projects on PipeOps. 🌐

Examples:
  - Logout from PipeOps:
    pipeops auth logout

  - Manage credentials:
    pipeops auth me`,
	Aliases: []string{"a"},
}

func init() {
	// Add the auth command as a subcommand of the root command
	rootCmd.AddCommand(authCmd)

	// Register subcommands under the auth command
	registerAuthSubcommands()
}

// registerAuthSubcommands initializes and registers subcommands for the auth command
func registerAuthSubcommands() {
	// Initialize and register authentication-related commands under the auth command
	authSub := auth.NewAuth(authCmd)
	authSub.Register()
}
