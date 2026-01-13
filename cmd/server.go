package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage server-related operations.",
	Long: `The server command provides a set of subcommands for managing
server-related operations on PipeOps, such as provisioning, configuration, and
interactions with servers.

Examples:
  - List all servers:
    pipeops server list

  - Create a new server:
    pipeops server create

  - Delete a server:
    pipeops server delete <server-id>`,
}

func init() {
	// Add the server command as a subcommand of the root command
	rootCmd.AddCommand(serverCmd)

	// Register subcommands under the server command
	registerServerSubcommands()
}

// registerServerSubcommands initializes and registers subcommands for the server command
func registerServerSubcommands() {
	// Add server commands
	serverCmd.AddCommand(server.GetListCmd())
	serverCmd.AddCommand(server.GetCreateCmd())
	serverCmd.AddCommand(server.GetDeleteCmd())
}
