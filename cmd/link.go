package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link [project-id]",
	Short: "🔗 Link current directory to a PipeOps project",
	Long: `🔗 Link the current directory to a PipeOps project.

This command creates a local context file that associates your current directory
with a specific PipeOps project, enabling project-aware commands like deploy, logs, and status.

Examples:
  - Link to a specific project:
    pipeops link my-project-id

  - Interactive project selection:
    pipeops link

  - Link and set custom name:
    pipeops link my-project-id --name "My Local App"`,
	Args: cobra.MaximumNArgs(1),
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

		var projectID string
		var selectedProject *models.Project

		if len(args) > 0 {
			// Project ID provided as argument
			projectID = args[0]

			// Verify project exists
			utils.PrintInfo(fmt.Sprintf("Verifying project %s...", projectID), opts)
			project, err := client.GetProject(projectID)
			if err != nil {
				utils.HandleError(err, "Error fetching project", opts)
				return
			}
			selectedProject = project
		} else {
			// Interactive project selection
			utils.PrintInfo("Fetching your projects...", opts)
			projectsResp, err := client.GetProjects()
			if err != nil {
				utils.HandleError(err, "Error fetching projects", opts)
				return
			}

			if len(projectsResp.Projects) == 0 {
				utils.PrintWarning("No projects found. Create a project first at https://app.pipeops.io", opts)
				return
			}

			// Show projects and let user select
			fmt.Printf("\n📋 Available Projects:\n")
			for i, project := range projectsResp.Projects {
				status := utils.GetStatusIcon(project.Status)
				fmt.Printf("  %d. %s %s (%s)\n", i+1, status, project.Name, project.ID)
			}

			// Get user selection
			var selection int
			fmt.Printf("\nSelect a project (1-%d): ", len(projectsResp.Projects))
			_, err = fmt.Scanf("%d", &selection)
			if err != nil || selection < 1 || selection > len(projectsResp.Projects) {
				utils.HandleError(fmt.Errorf("invalid selection"), "Invalid project selection", opts)
				return
			}

			selectedProject = &projectsResp.Projects[selection-1]
			projectID = selectedProject.ID
		}

		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			utils.HandleError(err, "Error getting current directory", opts)
			return
		}

		// Create project context
		context := &utils.ProjectContext{
			ProjectID:   projectID,
			ProjectName: selectedProject.Name,
			Directory:   currentDir,
		}

		// Save context to .pipeops directory
		if err := utils.SaveProjectContext(context); err != nil {
			utils.HandleError(err, "Error saving project context", opts)
			return
		}

		// Success message
		utils.PrintSuccess(fmt.Sprintf("Successfully linked directory to project '%s' (%s)", selectedProject.Name, projectID), opts)

		if !opts.Quiet {
			fmt.Printf("\n📁 PROJECT CONTEXT\n")
			fmt.Printf("├─ Project: %s (%s)\n", selectedProject.Name, projectID)
			fmt.Printf("├─ Directory: %s\n", currentDir)
			fmt.Printf("└─ Context file: %s\n", filepath.Join(currentDir, ".pipeops", "project.json"))

			fmt.Printf("\n💡 NEXT STEPS\n")
			fmt.Printf("├─ Deploy: pipeops deploy\n")
			fmt.Printf("├─ View logs: pipeops logs\n")
			fmt.Printf("├─ Check status: pipeops status\n")
			fmt.Printf("└─ Manage env vars: pipeops env\n")
		}
	},
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "🔓 Unlink project from current directory",
	Long: `🔓 Remove the project association from the current directory.
After unlinking, you'll need to specify project IDs in commands again.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error getting current directory: %v\n", err)
			return
		}

		pipeopsFile := filepath.Join(currentDir, ".pipeops")

		// Check if .pipeops file exists
		if _, err := os.Stat(pipeopsFile); os.IsNotExist(err) {
			fmt.Println("📭 No project is linked to this directory.")
			return
		}

		// Remove .pipeops file
		if err := os.Remove(pipeopsFile); err != nil {
			fmt.Printf("❌ Error removing .pipeops file: %v\n", err)
			return
		}

		fmt.Println("✅ Successfully unlinked project from current directory!")
		fmt.Printf("🗑️  Removed %s\n", pipeopsFile)
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(unlinkCmd)
}
