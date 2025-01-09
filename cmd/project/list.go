package project

import (
	"strings"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
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
		// Mock data to simulate fetching project details
		log.Info("🔍 Fetching all projects...\n")

		// Example project data
		projects := []struct {
			ID     string
			Name   string
			Status string
		}{
			{"proj-001", "My First Project", "Active"},
			{"proj-002", "Demo Project", "Inactive"},
			{"proj-003", "Test Pipeline Project", "Active"},
		}

		// Display header
		log.Infof("%-15s | %-30s | %-10s\n", "PROJECT ID", "PROJECT NAME", "STATUS")
		log.Info(strings.Repeat("-", 60))

		// Display project details
		for _, project := range projects {
			log.Infof("%-15s | %-30s | %-10s\n", project.ID, project.Name, project.Status)
		}

		log.Info("\n✅ Projects listed successfully.")
	},
	Args: cobra.NoArgs, // This command does not accept arguments
}

// NewList initializes and returns the list command
func (p *projectModel) listProjects() *cobra.Command {
	p.rootCmd.AddCommand(listCmd)
	return listCmd
}
