package project

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the command to list all projects
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long: `List all projects in your PipeOps account.

Examples:
  pipeops project list
  pipeops project list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Fetch projects from API
		projectsResp, err := client.GetProjects()
		if err != nil {
			utils.HandleError(err, "Error fetching projects", opts)
			return
		}

		if len(projectsResp.Projects) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON([]interface{}{})
			} else {
				fmt.Println("No projects found")
			}
			return
		}

		// Format output
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(projectsResp.Projects)
		} else {
			// Prepare table data
			headers := []string{"ID", "NAME", "STATUS", "CREATED"}
			var rows [][]string

			for _, project := range projectsResp.Projects {
				name := utils.TruncateString(project.Name, 30)
				status := project.Status
				created := utils.FormatDateShort(project.CreatedAt)

				rows = append(rows, []string{
					project.ID,
					name,
					status,
					created,
				})
			}

			utils.PrintTable(headers, rows, opts)
		}
	},
	Args: cobra.NoArgs,
}

// NewList initializes and returns the list command
func (p *projectModel) listProjects() *cobra.Command {
	p.rootCmd.AddCommand(listCmd)
	return listCmd
}
