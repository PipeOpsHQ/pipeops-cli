package cmd

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy addons to projects",
	Long: `Deploy addons to projects.

Project code deployment is temporarily disabled. You can deploy addons to existing projects.

Examples:
  - Deploy an addon to a project:
    pipeops deploy --addon postgres --project proj-123

  - Deploy addon with environment variables:
    pipeops deploy --addon redis --project proj-123 --env REDIS_PASSWORD=secret

  - Deploy addon to linked project:
    pipeops deploy --addon postgres --env POSTGRES_DB=myapp`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Parse flags
		addonID, _ := cmd.Flags().GetString("addon")
		projectID, _ := cmd.Flags().GetString("project")
		envVars, _ := cmd.Flags().GetStringToString("env")

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

		if addonID != "" {
			// Deploy addon
			deployAddon(client, addonID, projectID, envVars, opts)
		} else {
			// Project deployment is disabled
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON(map[string]string{
					"status":  "disabled",
					"message": "Project deployment is temporarily disabled",
				})
			} else {
				utils.PrintWarning("Project deployment is temporarily disabled. You can still deploy addons to existing projects.", opts)
				utils.PrintInfo("\nALTERNATIVES:", opts)
				utils.PrintInfo("  - Deploy addons: `pipeops deploy --addon <addon-id> --project <project-id>`", opts)
				utils.PrintInfo("  - List addons: `pipeops list --addons`", opts)
				utils.PrintInfo("  - Use PipeOps web console for project deployment: https://app.pipeops.io", opts)
			}
		}
	},
	Args: cobra.NoArgs,
}

func deployAddon(client pipeops.ClientAPI, addonID, projectID string, envVars map[string]string, opts utils.OutputOptions) {
	// Get project ID if not provided
	if projectID == "" {
		projectContext, err := utils.LoadProjectContext()
		if err != nil || projectContext.ProjectID == "" {
			utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use --project flag or link a project with 'pipeops link'", opts)
			return
		}
		projectID = projectContext.ProjectID
	}

	// Get addon information
	spinner := utils.StartSpinner(fmt.Sprintf("Getting addon '%s' information...", addonID), opts)
	addon, err := client.GetAddon(addonID)
	utils.StopSpinner(spinner)

	if err != nil {
		utils.HandleError(err, "Error fetching addon information", opts)
		return
	}

	// Create deployment request
	req := &models.AddonDeployRequest{
		AddonID:   addonID,
		ProjectID: projectID,
		Name:      addon.Name,
		EnvVars:   envVars,
	}

	// Deploy addon
	spinner = utils.StartSpinner(fmt.Sprintf("Deploying addon '%s' to project '%s'...", addon.Name, projectID), opts)
	deployResp, err := client.DeployAddon(req)
	utils.StopSpinner(spinner)

	if err != nil {
		utils.HandleError(err, "Error deploying addon", opts)
		return
	}

	if opts.Format == utils.OutputFormatJSON {
		utils.PrintJSON(deployResp)
	} else {
		utils.PrintSuccess(fmt.Sprintf("Addon '%s' deployed successfully!", addon.Name), opts)
		utils.PrintInfo(fmt.Sprintf("Deployment ID: %s", deployResp.DeploymentID), opts)
		utils.PrintInfo(fmt.Sprintf("Status: %s", deployResp.Status), opts)

		if deployResp.Message != "" {
			utils.PrintInfo(fmt.Sprintf("Message: %s", deployResp.Message), opts)
		}

		// Show helpful tips
		if !opts.Quiet {
			fmt.Printf("\nNEXT STEPS\n")
			fmt.Printf("├─ Check status: pipeops status --project %s\n", projectID)
			fmt.Printf("├─ View logs: pipeops logs --project %s\n", projectID)
			fmt.Printf("└─ List deployments: pipeops list --deployments --project %s\n", projectID)
		}
	}
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Add flags
	deployCmd.Flags().StringP("addon", "a", "", "Addon ID to deploy")
	deployCmd.Flags().StringP("project", "p", "", "Target project ID")
	deployCmd.Flags().StringToStringP("env", "e", nil, "Environment variables (KEY=VALUE)")
}
