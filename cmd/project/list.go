package project

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
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
		// Load configuration first
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Create client with the loaded configuration
		client := pipeops.NewClientWithConfigFunc(cfg)

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Fetch projects from API
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
				fmt.Println("No projects found yet")
				fmt.Println()
				fmt.Println("Ready to create your first project?")
				fmt.Println("   Visit: https://app.pipeops.io")
				fmt.Println("   Or check our docs for CLI project creation")
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

			// Encouraging summary with next steps
			fmt.Printf("\nFound %d project(s)\n", len(projectsResp.Projects))

			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Println("   pipeops project deploy <project-id>  # Deploy a project")
			fmt.Println("   pipeops project logs <project-id>    # View project logs")
		}
	},
	Args: cobra.NoArgs,
}

// NewList initializes and returns the list command
func (p *projectModel) listProjects() *cobra.Command {
	p.rootCmd.AddCommand(listCmd)
	return listCmd
}
