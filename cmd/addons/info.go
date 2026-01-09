package addons

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [addon-id]",
	Short: "Show addon details",
	Long: `Show detailed information about a specific addon.

If no addon ID is provided, an interactive selection will be shown.

Examples:
  - View addon details:
    pipeops addons info redis

  - Interactive selection:
    pipeops addons info`,
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

		var addonID string

		if len(args) > 0 {
			addonID = args[0]
		} else {
			// Interactive selection
			addonsResp, err := client.GetAddons()
			if err != nil {
				utils.HandleError(err, "Error fetching addons", opts)
				return
			}

			if len(addonsResp.Addons) == 0 {
				utils.PrintWarning("No addons available", opts)
				return
			}

			var options []string
			for _, addon := range addonsResp.Addons {
				options = append(options, fmt.Sprintf("%s (%s) - %s", addon.Name, addon.ID, addon.Category))
			}

			idx, _, err := utils.SelectOption("Select an addon", options)
			if err != nil {
				utils.HandleError(err, "Selection cancelled", opts)
				return
			}

			addonID = addonsResp.Addons[idx].ID
		}

		utils.PrintInfo(fmt.Sprintf("Getting addon '%s' information...", addonID), opts)

		addon, err := client.GetAddon(addonID)
		if err != nil {
			utils.HandleError(err, "Error fetching addon information", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(addon)
		} else {
			fmt.Printf("\nADDON DETAILS\n")
			fmt.Printf("├─ ID: %s\n", addon.ID)
			fmt.Printf("├─ Name: %s\n", addon.Name)
			fmt.Printf("├─ Category: %s\n", addon.Category)
			fmt.Printf("├─ Version: %s\n", addon.Version)
			fmt.Printf("├─ Status: %s %s\n", utils.GetStatusIcon(addon.Status), addon.Status)
			fmt.Printf("└─ Image: %s\n", addon.Image)

			if addon.Description != "" {
				fmt.Printf("\nDESCRIPTION\n")
				fmt.Printf("%s\n", addon.Description)
			}

			if len(addon.Tags) > 0 {
				fmt.Printf("\nTAGS\n")
				for i, tag := range addon.Tags {
					if i == len(addon.Tags)-1 {
						fmt.Printf("└─ %s\n", tag)
					} else {
						fmt.Printf("├─ %s\n", tag)
					}
				}
			}

			if !opts.Quiet {
				fmt.Printf("\nACTIONS\n")
				fmt.Printf("├─ Deploy: pipeops deploy --addon %s --project <project-id>\n", addon.ID)
				fmt.Printf("└─ List all addons: pipeops addons ls\n")
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	AddonsCmd.AddCommand(infoCmd)
}
