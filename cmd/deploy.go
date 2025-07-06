package cmd

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/cmd/deploy"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [project-id]",
	Short: "🚀 Deploy to a project",
	Long: `🚀 Deploy your application to a project. Uses the linked project if available.
This command provides access to deployment operations.

Examples:
  - Deploy to linked project:
    pipeops deploy

  - Deploy to specific project:
    pipeops deploy proj-123

  - Deploy pipeline:
    pipeops deploy pipeline --source ./my-app

Use 'pipeops deploy --help' to see all available deployment commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Get project ID
		var projectID string
		var err error

		if len(args) == 1 {
			projectID = args[0]
		} else {
			projectID, err = utils.GetProjectIDOrLinked("")
			if err != nil {
				utils.PrintError(err.Error(), opts)
				fmt.Println("\n💡 Available deployment commands:")
				fmt.Println("  - pipeops deploy pipeline    # Deploy a pipeline")
				fmt.Println("  - pipeops link <project-id>  # Link a project to this directory")
				return
			}
		}

		// Validate project ID
		if err := validation.ValidateProjectID(projectID); err != nil {
			utils.PrintError(fmt.Sprintf("Invalid project ID: %v", err), opts)
			return
		}

		// Show project context
		utils.PrintProjectContextWithOptions(projectID, opts)

		// Show deployment options
		if opts.Format == utils.OutputFormatJSON {
			deploymentInfo := map[string]interface{}{
				"project_id": projectID,
				"available_commands": []string{
					"pipeline",
				},
			}
			utils.PrintJSON(deploymentInfo)
		} else {
			fmt.Printf("\n🚀 DEPLOYMENT OPTIONS for project %s\n", projectID)
			fmt.Printf("├─ Pipeline deployment: pipeops deploy pipeline\n")
			fmt.Printf("└─ Use --help for more information\n")

			fmt.Printf("\n💡 QUICK START\n")
			fmt.Printf("Deploy a pipeline from current directory:\n")
			fmt.Printf("  pipeops deploy pipeline --source .\n")
		}
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Register subcommands under the deploy command
	registerDeploySubcommands()
}

// registerDeploySubcommands initializes and registers subcommands for the deploy command
func registerDeploySubcommands() {
	// Initialize and register deploy-related commands under the deploy command
	deployCmd := deploy.NewDeploy(deployCmd)
	deployCmd.Register()
}
