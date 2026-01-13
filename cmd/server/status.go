package server

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

func GetStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "status [server-id]",
		Aliases: []string{"show"},
		Short:   "Get the status of a server",
		Long: `Get the status of a specific server in your PipeOps account.

Examples:
  pipeops server status <server-id>
  pipeops server show <server-id>`,
		Run: func(cmd *cobra.Command, args []string) {
			opts := utils.GetOutputOptions(cmd)

			if len(args) == 0 {
				utils.HandleError(fmt.Errorf("server ID is required"), "Usage: pipeops server status <server-id>", opts)
				return
			}
			serverID := args[0]

			cfg, err := config.Load()
			if err != nil {
				utils.HandleError(err, "Error loading configuration", opts)
				return
			}

			client := pipeops.NewClientWithConfig(cfg)

			if !utils.RequireAuth(client, opts) {
				return
			}

			utils.PrintInfo(fmt.Sprintf("Fetching status for server %s...", serverID), opts)

			server, err := client.GetServer(serverID)
			if err != nil {
				if !utils.HandleAuthError(err, opts) {
					return
				}
				utils.HandleError(err, "Error fetching server status", opts)
				return
			}

			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON(server)
			} else {
				headers := []string{"ATTRIBUTE", "VALUE"}
				var rows [][]string

				rows = append(rows, []string{"ID", server.ID})
				rows = append(rows, []string{"Name", server.Name})
				rows = append(rows, []string{"Type", server.Type})
				rows = append(rows, []string{"Status", server.Status})
				rows = append(rows, []string{"Region", server.Region})
				rows = append(rows, []string{"IP Address", server.IP})
				rows = append(rows, []string{"Created At", utils.FormatDate(server.CreatedAt)})

				utils.PrintTable(headers, rows, opts)
			}
		},
		Args: cobra.ExactArgs(1),
	}
}
