package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/addons"
	"github.com/spf13/cobra"
)

var addonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "Manage addons",
	Long: `Manage addons in your PipeOps account.

Addons are pre-built services like databases, caches, and message queues
that can be deployed alongside your projects.

Examples:
  - List all deployable addons:
    pipeops addons
    pipeops addons available

  - View addon details:
    pipeops addons info <addon-id>

  - List deployed addons:
    pipeops addons list`,
	Aliases: []string{"addon"},
	Run:     addons.RunAvailableAddons,
	Args:    cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(addonsCmd)
	registerAddonsSubcommands()
}

func registerAddonsSubcommands() {
	// Add all subcommands from the addons package
	for _, cmd := range addons.AddonsCmd.Commands() {
		addonsCmd.AddCommand(cmd)
	}
}
