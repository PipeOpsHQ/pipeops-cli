package cmd

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List projects or addons",
	Long: `List all projects or addons in your PipeOps account.

Examples:
  - List all projects:
    pipeops list
    pipeops ls

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
					url := deployment.DeploymentURL
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
					fmt.Printf("\n[ TIPS ]\n")
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
					utils.PrintWarning("No projects found", opts)
					fmt.Printf("\n[ GET STARTED ]\n")
					fmt.Printf("├─ Create a project at: https://app.pipeops.io\n")
					fmt.Printf("├─ Import from GitHub: pipeops create --from-github\n")
					fmt.Printf("└─ Check documentation: https://docs.pipeops.io\n")
				}
				return
			}

			// Check if current directory is linked to a project
			linkedProjectID := ""
			if context, err := utils.LoadProjectContext(); err == nil {
				linkedProjectID = context.ProjectID
			}

			// Format output
			if opts.Format == utils.OutputFormatJSON {
				// Add linked status to JSON output
				type ProjectWithLink struct {
					*models.Project
					IsLinked bool `json:"is_linked"`
				}

				var projectsWithLink []ProjectWithLink
				for _, project := range projectsResp.Projects {
					p := project // Create a copy to avoid pointer issues
					projectsWithLink = append(projectsWithLink, ProjectWithLink{
						Project:  &p,
						IsLinked: p.ID == linkedProjectID,
					})
				}
				utils.PrintJSON(projectsWithLink)
			} else {
				// Enhanced table display
				fmt.Printf("\n[ PROJECTS OVERVIEW ]\n")
				fmt.Printf("├─ Total: %d projects\n", len(projectsResp.Projects))

				// Count projects by status
				statusCounts := make(map[string]int)
				for _, project := range projectsResp.Projects {
					statusCounts[project.Status]++
				}

				if len(statusCounts) > 0 {
					fmt.Printf("└─ Status: ")
					i := 0
					for status, count := range statusCounts {
						if i > 0 {
							fmt.Printf(", ")
						}
						fmt.Printf("%s %s (%d)", utils.GetStatusIcon(status), status, count)
						i++
					}
					fmt.Printf("\n")
				}

				fmt.Printf("\n")

				// Prepare enhanced table data
				headers := []string{"", "PROJECT ID", "PROJECT NAME", "STATUS", "ENVIRONMENT", "LAST UPDATED", "CREATED"}
				var rows [][]string

				for _, project := range projectsResp.Projects {
					// Check if this is the linked project
					linkedIndicator := "  "
					if project.ID == linkedProjectID {
						linkedIndicator = "→ "
					}

					name := utils.TruncateString(project.Name, 25)
					status := utils.GetStatusIcon(project.Status) + " " + project.Status

					// Get environment - using status as a proxy for environment
					environment := "production" // Default
					if project.Status == "development" || project.Status == "dev" {
						environment = "development"
					} else if project.Status == "staging" {
						environment = "staging"
					}

					// Format dates
					created := utils.FormatDateShort(project.CreatedAt)
					updated := utils.FormatDateShort(project.UpdatedAt)

					rows = append(rows, []string{
						linkedIndicator,
						project.ID,
						name,
						status,
						environment,
						updated,
						created,
					})
				}

				utils.PrintTable(headers, rows, opts)

				// Show linked project info
				if linkedProjectID != "" {
					fmt.Printf("\n→ Linked to current directory\n")
				}

				// Enhanced tips section
				if !opts.Quiet {
					fmt.Printf("\n[ QUICK ACTIONS ]\n")

					if linkedProjectID != "" {
						fmt.Printf("├─ Deploy linked project: pipeops deploy\n")
						fmt.Printf("├─ View linked status: pipeops status\n")
						fmt.Printf("├─ Unlink project: pipeops unlink\n")
					} else {
						fmt.Printf("├─ Link a project: pipeops link <project-id>\n")
						fmt.Printf("├─ Interactive link: pipeops link\n")
					}

					fmt.Printf("├─ View details: pipeops status <project-id>\n")
					fmt.Printf("├─ View logs: pipeops logs --project <project-id>\n")
					fmt.Printf("├─ List addons: pipeops list --addons\n")
					fmt.Printf("└─ Deploy addon: pipeops deploy --addon <addon-id>\n")

					// Add filtering hint
					fmt.Printf("\n[ COMING SOON ]\n")
					fmt.Printf("├─ Filter by status: pipeops list --status active\n")
					fmt.Printf("├─ Search projects: pipeops list --search <term>\n")
					fmt.Printf("└─ Sort options: pipeops list --sort name|created|updated\n")
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
