package addons

import "github.com/spf13/cobra"

var deploymentsCmd = &cobra.Command{
	Use:     "deployments",
	Aliases: []string{"deps"},
	Short:   "List addon deployments in your workspace",
	Long: `List all addon deployments in your workspace.

Examples:
  - List all addon deployments:
    pipeops addons deployments`,
	Run:  runAddonDeployments,
	Args: cobra.NoArgs,
}

func init() {
	AddonsCmd.AddCommand(deploymentsCmd)
}
