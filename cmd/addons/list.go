package addons

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available addons",
	Long: `List all available addons in the PipeOps catalog.

Examples:
  - List all addons:
    pipeops addons list
    pipeops addons ls

  - List addons in JSON format:
    pipeops addons ls --json`,
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

		utils.PrintInfo("Fetching available addons...", opts)

		addonsResp, err := client.GetAddons()
		if err != nil {
			utils.HandleError(err, "Error fetching addons", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(addonsResp.Addons)
		} else {
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
			utils.PrintSuccess(fmt.Sprintf("Found %d addons", len(addonsResp.Addons)), opts)

			if !opts.Quiet {
				fmt.Printf("\nACTIONS\n")
				fmt.Printf("├─ View details: pipeops addons info <addon-id>\n")
				fmt.Printf("├─ Deploy addon: pipeops deploy --addon <addon-id> --project <project-id>\n")
				fmt.Printf("└─ List deployments: pipeops addons deployments --project <project-id>\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	AddonsCmd.AddCommand(listCmd)
}
