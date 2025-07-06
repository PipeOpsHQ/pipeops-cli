package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link [project-id]",
	Short: "🔗 Link a project to the current directory",
	Long: `🔗 Associate an existing PipeOps project with the current directory.
This allows you to run commands without specifying the project ID each time.

Examples:
  - Link a project to current directory:
    pipeops link proj-123

  - Link with interactive selection:
    pipeops link

After linking, you can use simplified commands:
  - pipeops logs (instead of pipeops project logs proj-123)
  - pipeops shell web-service (instead of pipeops shell proj-123 web-service)`,
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

		var projectID string

		if len(args) == 1 {
			// Project ID provided as argument
			projectID = args[0]
			if err := validation.ValidateProjectID(projectID); err != nil {
				fmt.Printf("❌ Invalid project ID: %v\n", err)
				return
			}
		} else {
			// Interactive project selection
			fmt.Println("🔍 Fetching your projects...")
			projectsResp, err := client.GetProjects()
			if err != nil {
				fmt.Printf("❌ Error fetching projects: %v\n", err)
				return
			}

			if len(projectsResp.Projects) == 0 {
				fmt.Println("📭 No projects found. Create a project first with 'pipeops project create'.")
				return
			}

			// Display projects
			fmt.Println("\n📂 Available projects:")
			for i, project := range projectsResp.Projects {
				fmt.Printf("  %d. %s (%s) - %s\n", i+1, project.Name, project.ID, project.Status)
			}

			// Get user selection
			fmt.Print("\n🎯 Select a project (1-", len(projectsResp.Projects), "): ")
			var selection int
			if _, err := fmt.Scanln(&selection); err != nil || selection < 1 || selection > len(projectsResp.Projects) {
				fmt.Println("❌ Invalid selection.")
				return
			}

			projectID = projectsResp.Projects[selection-1].ID
		}

		// Verify project exists
		fmt.Printf("🔍 Verifying project %s...\n", projectID)
		project, err := client.GetProject(projectID)
		if err != nil {
			fmt.Printf("❌ Error accessing project: %v\n", err)
			return
		}

		// Create .pipeops file in current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error getting current directory: %v\n", err)
			return
		}

		pipeopsFile := filepath.Join(currentDir, ".pipeops")

		// Write project ID to .pipeops file
		content := fmt.Sprintf("project_id=%s\n", projectID)
		if err := os.WriteFile(pipeopsFile, []byte(content), 0644); err != nil {
			fmt.Printf("❌ Error creating .pipeops file: %v\n", err)
			return
		}

		fmt.Printf("✅ Successfully linked project '%s' (%s) to current directory!\n", project.Name, projectID)
		fmt.Printf("📁 Created %s\n", pipeopsFile)
		fmt.Println("\n🎉 You can now use simplified commands:")
		fmt.Println("  - pipeops logs")
		fmt.Println("  - pipeops shell web-service")
		fmt.Println("  - pipeops proxy start web-service")
		fmt.Println("\n💡 Tip: Add .pipeops to your .gitignore to keep it local.")
	},
	Args: cobra.MaximumNArgs(1),
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
