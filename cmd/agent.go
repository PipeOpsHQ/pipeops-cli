package cmd

import (
	cmd "github.com/PipeOpsHQ/pipeops-cli/cmd/agent"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agent-related commands and tasks.",
	Long: `The agent command provides subcommands to manage various aspects
of the PipeOps agent, such as setup, configuration, and interactions
with supported environments like EKS,GKE,AKS,DKS & K3s.`,
}

func init() {
	// Add the agent command as a subcommand of the root command
	rootCmd.AddCommand(agentCmd)

	// Initialize and register subcommands under the agent command
	registerAgentSubcommands()
}

// registerAgentSubcommands initializes and registers subcommands for the agent command
func registerAgentSubcommands() {
	// Create a new K3s-related command under the agent command
	k3sCmd := cmd.NewAgent(agentCmd)
	k3sCmd.Register()
}
