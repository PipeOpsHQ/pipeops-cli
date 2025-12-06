package project

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the project create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create project (temporarily disabled)",
	Long: `Project creation is temporarily disabled.

This feature is under development and will be available in a future release.

Available alternatives:
  - Use the PipeOps web console to create projects
  - Link existing projects: pipeops link <project-id>
  - List existing projects: pipeops project list`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(map[string]string{
				"status":  "disabled",
				"message": "Project creation is temporarily disabled",
			})
		} else {
			utils.PrintWarning("Project creation is temporarily disabled", opts)
			fmt.Printf("\nALTERNATIVES\n")
			fmt.Printf("├─ Use PipeOps web console to create projects\n")
			fmt.Printf("├─ Link existing projects: pipeops link <project-id>\n")
			fmt.Printf("├─ List projects: pipeops project list\n")
			fmt.Printf("└─ Deploy addons: pipeops deploy --addon <addon-id>\n")
		}
	},
	Args: cobra.NoArgs,
}

// NewCreate initializes and returns the create command
func (p *projectModel) createProject() *cobra.Command {
	p.rootCmd.AddCommand(createCmd)
	return createCmd
}
