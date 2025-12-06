package deploy

import (
	"fmt"
	"os"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Deploy current directory to PipeOps",
	Long: `Deploy the current directory to PipeOps using the linked project.

This command automatically detects your project type and deploys it to PipeOps.
Make sure you have linked a project first using 'pipeops link'.

Examples:
  - Deploy current directory:
    pipeops deploy pipeline

  - Deploy with custom source:
    pipeops deploy pipeline --source ./my-app

  - Deploy with custom name:
    pipeops deploy pipeline --name "My App v2.0"`,
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

		// Get project context
		context, err := utils.LoadProjectContext()
		if err != nil {
			utils.HandleError(fmt.Errorf("no linked project found. Run 'pipeops link' first"), "Project not linked", opts)
			return
		}

		// Get flags
		source, _ := cmd.Flags().GetString("source")
		name, _ := cmd.Flags().GetString("name")

		if source == "" {
			source = "."
		}

		// Get deployment name
		if name == "" {
			if dir, err := os.Getwd(); err == nil {
				name = fmt.Sprintf("%s-deployment", utils.GetBaseName(dir))
			} else {
				name = "cli-deployment"
			}
		}

		utils.PrintInfo(fmt.Sprintf("Deploying %s to project %s...", source, context.ProjectName), opts)

		// Create the project using the API
		req := &models.ProjectCreateRequest{
			Name:        name,
			Description: fmt.Sprintf("Project created from %s", source),
		}

		project, err := client.CreateProject(req)
		if err != nil {
			utils.HandleError(err, "Error creating deployment", opts)
			return
		}

		// Format output
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(project)
		} else {
			utils.PrintSuccess("Deployment initiated successfully!", opts)
			fmt.Printf("\nDEPLOYMENT DETAILS\n")
			fmt.Printf("├─ Project: %s (%s)\n", context.ProjectName, context.ProjectID)
			fmt.Printf("├─ Source: %s\n", source)
			fmt.Printf("├─ Deployment: %s\n", project.Name)
			fmt.Printf("└─ Status: %s\n", project.Status)

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\nNEXT STEPS\n")
				fmt.Printf("├─ View logs: pipeops logs\n")
				fmt.Printf("├─ Check status: pipeops status\n")
				fmt.Printf("└─ Open shell: pipeops shell\n")
			}
		}
	},
}

// NewPipeline initializes and returns the pipeline command
func (p *deployModel) newPipeline() *cobra.Command {
	// Add flags
	pipelineCmd.Flags().StringP("source", "s", "", "Source directory to deploy (default: current directory)")
	pipelineCmd.Flags().StringP("name", "n", "", "Custom name for deployment")

	// Add the pipeline command as a subcommand to the parent command
	p.rootCmd.AddCommand(pipelineCmd)
	return pipelineCmd
}

// RegisterPipelineSubcommands initializes and registers subcommands for the pipeline command
func (p *deployModel) RegisterPipelineSubcommands() {
	// Add subcommands related to pipelines
	pipelineCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all pipelines",
		Long: `The "list" subcommand displays all the deployment pipelines in your project.

Example:
  pipeops deploy pipeline list`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing all pipelines...")
		},
	})

	pipelineCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new deployment pipeline",
		Long: `The "create" subcommand creates a new deployment pipeline in PipeOps.

Example:
  pipeops deploy pipeline create --name my-pipeline`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating a new pipeline...")
		},
	})

	pipelineCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a deployment pipeline",
		Long: `The "delete" subcommand deletes an existing deployment pipeline in PipeOps.

Example:
  pipeops deploy pipeline delete --id pipeline-id`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Deleting a pipeline...")
		},
	})
}
