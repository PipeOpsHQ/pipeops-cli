package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/agent"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "⚙️ Manage agent-related commands and tasks",
	Long: `⚙️ The "agent" command provides subcommands to manage various aspects
of the PipeOps agent. These include setup, configuration, and interactions 
with supported environments like EKS, GKE, AKS, DKS, and K3s. 🌐

Examples:
  - Install the PipeOps agent:
    pipeops agent install`,
}

func init() {
	// Add the agent command as a subcommand of the root command
	rootCmd.AddCommand(agentCmd)

	// Initialize and register subcommands under the agent command
	registerAgentSubcommands()
}

// registerAgentSubcommands initializes and registers subcommands for the agent command
func registerAgentSubcommands() {
	// Create and register K3s-related commands under the agent command
	agentSub := agent.NewAgent(agentCmd)
	agentSub.Register()
}
