package cmd

import (
	cmd "github.com/PipeOpsHQ/pipeops-cli/cmd/k3s"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage server-related operations.",
	Long: `The server command provides a set of subcommands for managing
server-related operations, such as provisioning, configurations and interactions
with servers on PipeOps`,
}

func init() {
	// Add the server command as a subcommand of the root command
	rootCmd.AddCommand(serverCmd)

	// Register subcommands under the server command
	registerServerSubcommands()
}

// registerServerSubcommands initializes and registers subcommands for the server command
func registerServerSubcommands() {
	// Initialize K3s-related commands under the server command
	k3sCmd := cmd.NewK3s(serverCmd)
	k3sCmd.Register()
}
