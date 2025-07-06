package cmd

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "📜 List all projects",
	Long: `📜 List all projects in your PipeOps account.

Examples:
  - List all projects:
    pipeops list

  - List projects in JSON format:
    pipeops list --json`,
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
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
