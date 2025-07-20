package cmd

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "📜 List projects or addons",
	Long: `📜 List all projects or addons in your PipeOps account.

Examples:
  - List all projects:
    pipeops list

  - List all addons:
    pipeops list --addons

  - List projects in JSON format:
    pipeops list --json

  - List addon deployments for a project:
    pipeops list --deployments --project proj-123`,
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

		// Parse flags
		showAddons, _ := cmd.Flags().GetBool("addons")
		showDeployments, _ := cmd.Flags().GetBool("deployments")
		projectID, _ := cmd.Flags().GetString("project")

		if showDeployments {
			// List addon deployments for a project
			if projectID == "" {
				// Try to get from linked project
				projectContext, err := utils.LoadProjectContext()
				if err != nil || projectContext.ProjectID == "" {
					utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use --project flag or link a project with 'pipeops link'", opts)
					return
				}
				projectID = projectContext.ProjectID
			}

			utils.PrintInfo(fmt.Sprintf("Fetching addon deployments for project '%s'...", projectID), opts)

			deployments, err := client.GetAddonDeployments(projectID)
			if err != nil {
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

		} else if showAddons {
			// List available addons
			utils.PrintInfo("Fetching available addons...", opts)

			addonsResp, err := client.GetAddons()
			if err != nil {
				utils.HandleError(err, "Error fetching addons", opts)
				return
			}

			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON(addonsResp.Addons)
			} else {
				if len(addonsResp.Addons) == 0 {
					utils.PrintWarning("No addons found", opts)
					return
				}

				headers := []string{"ADDON ID", "NAME", "CATEGORY", "VERSION", "STATUS"}
				var rows [][]string

				for _, addon := range addonsResp.Addons {
					name := utils.TruncateString(addon.Name, 30)
					status := utils.GetStatusIcon(addon.Status) + " " + addon.Status

					rows = append(rows, []string{
						addon.ID,
						name,
						addon.Category,
						addon.Version,
						status,
					})
				}

				utils.PrintTable(headers, rows, opts)
				utils.PrintSuccess(fmt.Sprintf("Found %d addons", len(addonsResp.Addons)), opts)

				// Show helpful tips
				if !opts.Quiet {
					fmt.Printf("\n💡 TIPS\n")
					fmt.Printf("├─ Deploy addon: pipeops deploy --addon <addon-id> --project <project-id>\n")
					fmt.Printf("├─ View deployments: pipeops list --deployments --project <project-id>\n")
					fmt.Printf("└─ Get addon info: pipeops status --addon <addon-id>\n")
				}
			}

		} else {
			// List projects (default behavior)
			utils.PrintInfo("Fetching all projects...", opts)

			projectsResp, err := client.GetProjects()
			if err != nil {
				// Handle authentication errors specifically
				if !utils.HandleAuthError(err, opts) {
					return
				}
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
					fmt.Printf("├─ List addons: pipeops list --addons\n")
					fmt.Printf("└─ View project: pipeops status <project-id>\n")
				}
			}
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().Bool("addons", false, "List available addons instead of projects")
	listCmd.Flags().Bool("deployments", false, "List addon deployments for a project")
	listCmd.Flags().StringP("project", "p", "", "Project ID (for listing deployments)")
}
