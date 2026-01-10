package addons

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var deploymentsCmd = &cobra.Command{
	Use:     "deployments",
	Aliases: []string{"deps"},
	Short:   "List addon deployments in your workspace",
	Long: `List all addon deployments in your workspace.

Examples:
  - List all addon deployments:
    pipeops addons deployments`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		client := pipeops.NewClient()

		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		if !utils.RequireAuth(client, opts) {
			return
		}

		utils.PrintInfo("Fetching addon deployments...", opts)

		deployments, err := client.GetAddonDeployments("")
		if err != nil {
			// Check if it's a 500 error (API not fully implemented)
			if strings.Contains(err.Error(), "500") {
				utils.PrintWarning("The addon deployments API is not yet available. Please check the PipeOps dashboard for addon deployments.", opts)
				return
			}
			utils.HandleError(err, "Error fetching addon deployments", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(deployments)
		} else {
			if len(deployments) == 0 {
				utils.PrintWarning("No addon deployments found in this workspace", opts)
				return
			}

			headers := []string{"ID", "NAME", "CATEGORY", "STATUS", "ENVIRONMENT", "URL"}
			var rows [][]string

			for _, deployment := range deployments {
				url := deployment.DeploymentURL
				if url == "" {
					url = "N/A"
				} else {
					url = utils.TruncateString(url, 40)
				}

				rows = append(rows, []string{
					utils.TruncateString(deployment.ID, 20),
					deployment.Name,
					deployment.Category,
					utils.GetStatusIcon(deployment.Status) + " " + deployment.Status,
					deployment.Environment,
					url,
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d addon deployments", len(deployments)), opts)
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	AddonsCmd.AddCommand(deploymentsCmd)
}
