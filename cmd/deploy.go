package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/deploy"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "🚀 Manage deployment-related operations",
	Long: `🔧 The deploy command provides a set of subcommands for managing 
deployment-related operations, such as configuring deployments, monitoring, 
and interacting with deployments on PipeOps. 🌐`,
	Aliases: []string{"d"},
}

func init() {
	// Add the deploy command as a subcommand of the root command
	rootCmd.AddCommand(deployCmd)

	// Register subcommands under the deploy command
	registerDeploySubcommands()
}

// registerDeploySubcommands initializes and registers subcommands for the deploy command
func registerDeploySubcommands() {
	// Initialize and register deploy-related commands under the deploy command
	deployCmd := deploy.NewDeploy(deployCmd)
	deployCmd.Register()
}
