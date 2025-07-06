package project

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/spf13/cobra"
)

// listCmd represents the command to list all projects
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "📜 List all projects",
	Long: `📜 The "list" command displays all the projects associated with your PipeOps account.
You can use this command to quickly view project details like name, ID, and status.

Examples:
  - List all projects:
    pipeops project list`,
	Run: func(cmd *cobra.Command, args []string) {
		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if user is authenticated
		if !client.IsAuthenticated() {
			fmt.Println("❌ You are not logged in. Please run 'pipeops auth login' first.")
			return
		}

		// Fetch projects from API
		fmt.Printf("🔍 Fetching all projects...\n")

		projectsResp, err := client.GetProjects()
		if err != nil {
			fmt.Printf("❌ Error fetching projects: %v\n", err)
			return
		}

		if len(projectsResp.Projects) == 0 {
			fmt.Println("📭 No projects found. Create your first project to get started!")
			return
		}

		// Display header
		fmt.Printf("%-15s | %-30s | %-10s | %-12s\n", "PROJECT ID", "PROJECT NAME", "STATUS", "CREATED")
		fmt.Println(strings.Repeat("-", 80))

		// Display project details
		for _, project := range projectsResp.Projects {
			// Truncate long names to fit in table
			name := project.Name
			if len(name) > 30 {
				name = name[:27] + "..."
			}

			// Format created date
			createdDate := project.CreatedAt.Format("2006-01-02")

			fmt.Printf("%-15s | %-30s | %-10s | %-12s\n",
				project.ID, name, project.Status, createdDate)
		}

		fmt.Printf("\n✅ Found %d projects.\n", len(projectsResp.Projects))
	},
	Args: cobra.NoArgs, // This command does not accept arguments
}

// NewList initializes and returns the list command
func (p *projectModel) listProjects() *cobra.Command {
	p.rootCmd.AddCommand(listCmd)
	return listCmd
}
