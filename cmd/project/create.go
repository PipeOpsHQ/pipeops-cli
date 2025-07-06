package project

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "✨ Create a new project",
	Long: `✨ The "create" command creates a new project in your PipeOps account.
You can specify the project name and optionally provide a description.

Examples:
  - Create a project with just a name:
    pipeops project create --name "My New Project"

  - Create a project with name and description:
    pipeops project create --name "My New Project" --description "A sample project"`,
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

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			fmt.Println("❌ Project name is required. Use --name flag to specify the name.")
			return
		}

		// Validate project name
		if err := validation.ValidateProjectName(name); err != nil {
			fmt.Printf("❌ Invalid project name: %v\n", err)
			return
		}

		// Validate project description if provided
		if description != "" {
			if err := validation.ValidateProjectDescription(description); err != nil {
				fmt.Printf("❌ Invalid project description: %v\n", err)
				return
			}
		}

		// Create project
		fmt.Printf("🔍 Creating project '%s'...\n", name)

		project, err := client.CreateProject(name, description)
		if err != nil {
			fmt.Printf("❌ Error creating project: %v\n", err)
			return
		}

		fmt.Println("✅ Project created successfully!")
		fmt.Printf("🆔 Project ID: %s\n", project.ID)
		fmt.Printf("📝 Name: %s\n", project.Name)
		if project.Description != "" {
			fmt.Printf("📄 Description: %s\n", project.Description)
		}
		fmt.Printf("📊 Status: %s\n", project.Status)
		fmt.Printf("📅 Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
	},
	Args: cobra.NoArgs,
}

func init() {
	createCmd.Flags().StringP("name", "n", "", "Name of the project (required)")
	createCmd.Flags().StringP("description", "d", "", "Description of the project (optional)")
	createCmd.MarkFlagRequired("name")
}

func (p *projectModel) createProject() {
	p.rootCmd.AddCommand(createCmd)
}
