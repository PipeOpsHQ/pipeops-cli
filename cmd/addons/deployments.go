package addons

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var deploymentsCmd = &cobra.Command{
	Use:     "deployments",
	Aliases: []string{"deps"},
	Short:   "List addon deployments for a project",
	Long: `List all addon deployments for a specific project.

If no project ID is provided, an interactive selection will be shown.

Examples:
  - List deployments for a project:
    pipeops addons deployments --project proj-123

  - Interactive project selection:
    pipeops addons deployments`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		client := pipeops.NewClient()

		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		if !utils.RequireAuth(client, opts) {
			return
		}

		projectID, _ := cmd.Flags().GetString("project")

		if projectID == "" {
			// Try linked project first
			projectContext, err := utils.LoadProjectContext()
			if err == nil && projectContext.ProjectID != "" {
				projectID = projectContext.ProjectID
			} else {
				// Interactive project selection
				projectsResp, err := client.GetProjects()
				if err != nil {
					utils.HandleError(err, "Error fetching projects", opts)
					return
				}

				if len(projectsResp.Projects) == 0 {
					utils.PrintWarning("No projects found", opts)
					return
				}

				var options []string
				for _, p := range projectsResp.Projects {
					status := utils.GetStatusIcon(p.Status)
					options = append(options, fmt.Sprintf("%s %s (%s)", status, p.Name, p.ID))
				}

				idx, _, err := utils.SelectOption("Select a project", options)
				if err != nil {
					utils.HandleError(err, "Selection cancelled", opts)
					return
				}

				projectID = projectsResp.Projects[idx].ID
			}
		}

		utils.PrintInfo(fmt.Sprintf("Fetching addon deployments for project '%s'...", projectID), opts)

		deployments, err := client.GetAddonDeployments(projectID)
		if err != nil {
			// Check if it's a 500 error (API not fully implemented)
			if strings.Contains(err.Error(), "500") {
				utils.PrintWarning("The addon deployments API is not yet available. Please check the PipeOps dashboard for addon deployments.", opts)
				return
			}
			utils.HandleError(err, "Error fetching addon deployments", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(deployments)
		} else {
			if len(deployments) == 0 {
				utils.PrintWarning("No addon deployments found for this project", opts)
				return
			}

			headers := []string{"DEPLOYMENT ID", "ADDON NAME", "STATUS", "URL", "CREATED"}
			var rows [][]string

			for _, deployment := range deployments {
				url := deployment.URL
				if url == "" {
					url = "N/A"
				}

				rows = append(rows, []string{
					deployment.ID,
					deployment.Name,
					utils.GetStatusIcon(deployment.Status) + " " + deployment.Status,
					url,
					utils.FormatDateShort(deployment.CreatedAt),
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d addon deployments", len(deployments)), opts)
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	AddonsCmd.AddCommand(deploymentsCmd)
	deploymentsCmd.Flags().StringP("project", "p", "", "Project ID")
}
