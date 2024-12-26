package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/k3s"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Server related commands",
}


func init() {
	// Register the server command
	rootCmd.AddCommand(serverCmd)

	// Use the serverCmd to group k3s related commands under server
	k3sCmd := cmd.NewK3s(serverCmd)
	k3sCmd.Register()
}
