package addons

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List deployed addons in your workspace",
	Long: `List addon deployments in your current PipeOps workspace.

Examples:
  - List deployed addons:
    pipeops addons list
    pipeops addons ls

  - List deployed addons in JSON format:
    pipeops addons ls --json`,
	Run:  runAddonDeployments,
	Args: cobra.NoArgs,
}

var availableCmd = &cobra.Command{
	Use:     "available",
	Aliases: []string{"catalog"},
	Short:   "List deployable addons",
	Long: `List all addons available to deploy from the PipeOps catalog.

Examples:
  - List deployable addons:
    pipeops addons
    pipeops addons available

  - List deployable addons in JSON format:
    pipeops addons --json
    pipeops addons available --json`,
	Run:  RunAvailableAddons,
	Args: cobra.NoArgs,
}

func RunAvailableAddons(cmd *cobra.Command, args []string) {
	opts := utils.GetOutputOptions(cmd)
	client := pipeops.NewClient()

	if err := client.LoadConfig(); err != nil {
		utils.HandleError(err, "Error loading configuration", opts)
		return
	}

	if !utils.RequireAuth(client, opts) {
		return
	}

	utils.PrintInfo("Fetching deployable addons...", opts)

	addonsResp, err := client.GetAddons()
	if err != nil {
		utils.HandleError(err, "Error fetching addons", opts)
		return
	}

	if opts.Format == utils.OutputFormatJSON {
		utils.PrintJSON(addonsResp.Addons)
		return
	}

	if len(addonsResp.Addons) == 0 {
		utils.PrintWarning("No addons found", opts)
		return
	}

	headers := []string{"ID", "NAME", "CATEGORY", "STATUS"}
	var rows [][]string

	for _, addon := range addonsResp.Addons {
		name := utils.TruncateString(addon.Name, 30)
		status := utils.GetStatusIcon(addon.Status) + " " + addon.Status

		rows = append(rows, []string{
			addon.ID,
			name,
			addon.Category,
			status,
		})
	}

	utils.PrintTable(headers, rows, opts)
	utils.PrintSuccess(fmt.Sprintf("Found %d deployable addons", len(addonsResp.Addons)), opts)

	if !opts.Quiet {
		fmt.Printf("\nACTIONS\n")
		fmt.Printf("├─ View details: pipeops addons info <addon-id>\n")
		fmt.Printf("├─ Deploy addon: pipeops addons deploy <addon-id>\n")
		fmt.Printf("└─ List deployed addons: pipeops addons list\n")
	}
}

func runAddonDeployments(cmd *cobra.Command, args []string) {
	opts := utils.GetOutputOptions(cmd)
	client := pipeops.NewClient()

	if err := client.LoadConfig(); err != nil {
		utils.HandleError(err, "Error loading configuration", opts)
		return
	}

	if !utils.RequireAuth(client, opts) {
		return
	}

	utils.PrintInfo("Fetching deployed addons...", opts)

	deployments, err := client.GetAddonDeployments()
	if err != nil {
		if strings.Contains(err.Error(), "500") {
			utils.PrintWarning("The addon deployments API is not yet available. Please check the PipeOps dashboard for addon deployments.", opts)
			return
		}
		utils.HandleError(err, "Error fetching addon deployments", opts)
		return
	}

	if opts.Format == utils.OutputFormatJSON {
		utils.PrintJSON(deployments)
		return
	}

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

func init() {
	AddonsCmd.AddCommand(listCmd, availableCmd)
}
