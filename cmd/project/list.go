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
	Short: "📜 List all projects",
	Long: `📜 List all projects in your PipeOps account.

Examples:
  - List all projects:
    pipeops project list

  - List projects in JSON format:
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
		utils.PrintInfo("Fetching all projects...", opts)

		projectsResp, err := client.GetProjects()
		if err != nil {
			utils.HandleError(err, "Error fetching projects", opts)
			return
		}

		if len(projectsResp.Projects) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON([]interface{}{})
			} else {
				utils.PrintWarning("No projects found. Create your first project to get started!", opts)
			}
			return
		}

		// Format output
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(projectsResp.Projects)
		} else {
			// Prepare table data
			headers := []string{"PROJECT ID", "PROJECT NAME", "STATUS", "CREATED"}
			var rows [][]string

			for _, project := range projectsResp.Projects {
				name := utils.TruncateString(project.Name, 30)
				status := utils.GetStatusIcon(project.Status) + " " + project.Status
				created := utils.FormatDateShort(project.CreatedAt)

				rows = append(rows, []string{
					project.ID,
					name,
					status,
					created,
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d projects", len(projectsResp.Projects)), opts)

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\n💡 TIPS\n")
				fmt.Printf("├─ Link a project: pipeops link <project-id>\n")
				fmt.Printf("├─ Create project: pipeops create <project-name>\n")
				fmt.Printf("└─ View project: pipeops status <project-id>\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

// NewList initializes and returns the list command
func (p *projectModel) listProjects() *cobra.Command {
	p.rootCmd.AddCommand(listCmd)
	return listCmd
}
