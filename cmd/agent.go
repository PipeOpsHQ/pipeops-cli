package cmd

import (
	cmd "github.com/PipeOpsHQ/pipeops-cli/cmd/agent"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Agent related commands",
}

func init() {
	// Register the server command
	rootCmd.AddCommand(agentCmd)

	// Use the serverCmd to group k3s related commands under server
	k3sCmd := cmd.NewAgent(agentCmd)
	k3sCmd.Register()
}
