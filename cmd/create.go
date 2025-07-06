package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "🏗️ Create a new project",
	Long: `🏗️ Create a new project in your PipeOps account.

Examples:
  - Create a project interactively:
    pipeops create

  - Create a project with a name:
    pipeops create my-awesome-project

  - Create a project with description:
    pipeops create my-project --description "My awesome project"

  - Create and immediately link to current directory:
    pipeops create my-project --link`,
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

		// Get project name
		var projectName string
		if len(args) == 1 {
			projectName = args[0]
		} else {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintError("Project name is required for JSON output", opts)
				return
			}
			// Interactive mode
			name, err := utils.PromptUser("🎯 Enter project name: ")
			if err != nil {
				utils.HandleError(err, "Error reading input", opts)
				return
			}
			projectName = name
		}

		// Validate project name
		if err := validation.ValidateProjectName(projectName); err != nil {
			utils.PrintError(fmt.Sprintf("Invalid project name: %v", err), opts)
			return
		}

		// Get description from flag or prompt
		description, _ := cmd.Flags().GetString("description")
		if description == "" && opts.Format != utils.OutputFormatJSON {
			description = utils.PromptUserWithDefault("📝 Enter project description (optional)", "")
		}

		// Validate description if provided
		if description != "" {
			if err := validation.ValidateProjectDescription(description); err != nil {
				utils.PrintError(fmt.Sprintf("Invalid project description: %v", err), opts)
				return
			}
		}

		// Create project
		utils.PrintInfo("Creating project...", opts)

		project, err := client.CreateProject(projectName, description)
		if err != nil {
			utils.HandleError(err, "Error creating project", opts)
			return
		}

		// Output result
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(project)
		} else {
			utils.PrintSuccess(fmt.Sprintf("Created project '%s' with ID: %s", project.Name, project.ID), opts)

			// Show project details
			fmt.Printf("\n📂 PROJECT DETAILS\n")
			fmt.Printf("├─ Name: %s\n", project.Name)
			fmt.Printf("├─ ID: %s\n", project.ID)
			fmt.Printf("├─ Status: %s%s\n", utils.GetStatusIcon(project.Status), project.Status)
			if project.Description != "" {
				fmt.Printf("└─ Description: %s\n", project.Description)
			} else {
				fmt.Printf("└─ Description: (none)\n")
			}
		}

		// Handle linking
		shouldLink, _ := cmd.Flags().GetBool("link")
		if !shouldLink && opts.Format != utils.OutputFormatJSON {
			shouldLink = utils.ConfirmAction("🔗 Link this project to the current directory?")
		}

		if shouldLink {
			if err := linkProject(project.ID, opts); err != nil {
				utils.PrintWarning(fmt.Sprintf("Failed to link project: %v", err), opts)
			} else {
				utils.PrintSuccess("Project linked to current directory!", opts)
				if opts.Format != utils.OutputFormatJSON {
					fmt.Printf("\n🎉 You can now use simplified commands:\n")
					fmt.Printf("  - pipeops logs\n")
					fmt.Printf("  - pipeops shell <service-name>\n")
					fmt.Printf("  - pipeops status\n")
				}
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
}

// linkProject creates a .pipeops file in the current directory
func linkProject(projectID string, opts utils.OutputOptions) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	pipeopsFile := filepath.Join(currentDir, ".pipeops")
	content := fmt.Sprintf("project_id=%s\n", projectID)

	return os.WriteFile(pipeopsFile, []byte(content), 0644)
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Add flags
	createCmd.Flags().StringP("description", "d", "", "Project description")
	createCmd.Flags().BoolP("link", "l", false, "Link the project to the current directory after creation")
}
