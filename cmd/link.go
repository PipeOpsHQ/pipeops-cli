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
	Short: "Link current directory to a PipeOps project",
	Long: `Link the current directory to a PipeOps project.

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
			fmt.Printf("\nAvailable Projects:\n")
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
			fmt.Printf("\nPROJECT CONTEXT\n")
			fmt.Printf("├─ Project: %s (%s)\n", selectedProject.Name, projectID)
			fmt.Printf("├─ Directory: %s\n", currentDir)
			fmt.Printf("└─ Context file: %s\n", filepath.Join(currentDir, ".pipeops", "project.json"))

			fmt.Printf("\nNEXT STEPS\n")
			fmt.Printf("├─ Deploy: pipeops deploy\n")
			fmt.Printf("├─ View logs: pipeops logs\n")
			fmt.Printf("├─ Check status: pipeops status\n")
			fmt.Printf("└─ Manage env vars: pipeops env\n")
		}
	},
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink project from current directory",
	Long: `Remove the project association from the current directory.

This command removes both the new context format (.pipeops/project.json) and 
the legacy format (.pipeops file) to ensure complete unlinking.

After unlinking, you'll need to specify project IDs in commands again.

Examples:
  - Unlink current directory:
    pipeops unlink
    
  - Force unlink (no confirmation):
    pipeops unlink --force`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		force, _ := cmd.Flags().GetBool("force")

		currentDir, err := os.Getwd()
		if err != nil {
			utils.HandleError(err, "Error getting current directory", opts)
			return
		}

		// Check if any project is linked
		context, contextErr := utils.LoadProjectContext()
		hasContext := contextErr == nil

		// Check for legacy .pipeops file
		legacyFile := filepath.Join(currentDir, ".pipeops")
		hasLegacy := false
		if _, err := os.Stat(legacyFile); err == nil {
			hasLegacy = true
		}

		// Check for new context directory
		contextDir := filepath.Join(currentDir, ".pipeops")
		contextFile := filepath.Join(contextDir, "project.json")
		hasContextDir := false
		if _, err := os.Stat(contextFile); err == nil {
			hasContextDir = true
		}

		if !hasContext && !hasLegacy && !hasContextDir {
			utils.PrintWarning("No project is linked to this directory", opts)
			return
		}

		// Show what will be unlinked
		if hasContext && !opts.Quiet {
			fmt.Printf("\nCURRENT LINK\n")
			fmt.Printf("├─ Project: %s (%s)\n", context.ProjectName, context.ProjectID)
			fmt.Printf("├─ Directory: %s\n", context.Directory)
			fmt.Printf("└─ Linked at: %s\n", utils.FormatDate(context.LinkedAt))
		}

		// Confirm unlinking unless force flag is set
		if !force && !opts.Quiet {
			if !utils.ConfirmAction("\nAre you sure you want to unlink this project?") {
				utils.PrintInfo("Unlink cancelled", opts)
				return
			}
		}

		// Track what was removed
		var removedItems []string

		// Remove new context directory and all its contents
		if hasContextDir {
			if err := os.RemoveAll(contextDir); err != nil {
				utils.PrintWarning(fmt.Sprintf("Could not remove .pipeops directory: %v", err), opts)
			} else {
				removedItems = append(removedItems, ".pipeops/")
			}
		}

		// Remove legacy .pipeops file if it exists and wasn't removed as part of directory
		if hasLegacy && !hasContextDir {
			if err := os.Remove(legacyFile); err != nil {
				utils.PrintWarning(fmt.Sprintf("Could not remove legacy .pipeops file: %v", err), opts)
			} else {
				removedItems = append(removedItems, ".pipeops")
			}
		}

		// Success message
		if len(removedItems) > 0 {
			utils.PrintSuccess("Successfully unlinked project from current directory", opts)

			if !opts.Quiet {
				fmt.Printf("\n🗑️  REMOVED\n")
				for i, item := range removedItems {
					if i == len(removedItems)-1 {
						fmt.Printf("└─ %s\n", item)
					} else {
						fmt.Printf("├─ %s\n", item)
					}
				}

				fmt.Printf("\n💡 NEXT STEPS\n")
				fmt.Printf("├─ Link another project: pipeops link\n")
				fmt.Printf("├─ List projects: pipeops list\n")
				fmt.Printf("└─ Specify project ID directly in commands\n")
			}
		} else {
			utils.PrintError("Failed to unlink project", opts)
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(unlinkCmd)

	// Add flags for unlink command
	unlinkCmd.Flags().BoolP("force", "f", false, "Force unlink without confirmation")
}
