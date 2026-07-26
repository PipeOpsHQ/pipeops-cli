package addons

import "github.com/spf13/cobra"

var deploymentsCmd = &cobra.Command{
	Use:     "deployments",
	Aliases: []string{"deps"},
	Short:   "List addon deployments in your workspace",
	Long: `List all addon deployments in your workspace.

Examples:
  - List all addon deployments:
    pipeops addons deployments
    pipeops addons deployments --workspace <workspace-uuid>`,
	Run:  runAddonDeployments,
	Args: cobra.NoArgs,
}

func init() {
	deploymentsCmd.Flags().String("workspace", "", "Workspace UUID (or set PIPEOPS_WORKSPACE_UUID / pipeops workspace select)")
	AddonsCmd.AddCommand(deploymentsCmd)
}
