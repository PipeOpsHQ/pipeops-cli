package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create project (temporarily disabled)",
	Long: `Project creation is temporarily disabled.

This feature is under development and will be available in a future release.

Available alternatives:
  - Use the PipeOps web console to create projects
  - Link existing projects: pipeops link <project-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(map[string]string{
				"status":  "disabled",
				"message": "Project creation is temporarily disabled",
			})
		} else {
			utils.PrintWarning("Project creation is temporarily disabled. This feature is under development and will be available in a future release.", opts)
			utils.PrintInfo("\nAvailable alternatives:", opts)
			utils.PrintInfo("  - Use the PipeOps web console to create projects: https://app.pipeops.io", opts)
			utils.PrintInfo("  - Link existing projects: `pipeops link <project-id>`", opts)
			utils.PrintInfo("  - List projects: `pipeops list`", opts)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
