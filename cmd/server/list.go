package server

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the command to list all servers
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all servers",
	Long: `List all servers in your PipeOps account.

Examples:
  pipeops server list
  pipeops server ls
  pipeops server list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		// Load configuration first
		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Allow overriding the workspace UUID per invocation.
		if workspaceUUID, _ := cmd.Flags().GetString("workspace"); workspaceUUID != "" {
			cfg.Settings.DefaultWorkspaceUUID = workspaceUUID
		}

		// Create client with the loaded configuration
		client := pipeops.NewClientWithConfig(cfg)

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Fetch servers from API
		utils.PrintInfo("Fetching all servers...", opts)

		serversResp, err := client.GetServers()
		if err != nil {
			// Handle authentication errors specifically
			if !utils.HandleAuthError(err, opts) {
				return
			}
			utils.HandleError(err, "Error fetching servers", opts)
			return
		}

		if len(serversResp.Servers) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON([]interface{}{})
			} else {
				fmt.Println("No servers found yet")
				fmt.Println()
				fmt.Println("Ready to create your first server?")
				fmt.Println("   Visit: https://app.pipeops.io")
			}
			return
		}

		// Format output
		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(serversResp.Servers)
		} else {
			// Prepare table data
			headers := []string{"SERVER ID", "NAME", "TYPE", "STATUS", "REGION", "IP", "CREATED"}
			var rows [][]string

			for _, server := range serversResp.Servers {
				name := utils.TruncateString(server.Name, 25)
				status := utils.GetStatusIcon(server.Status) + " " + server.Status
				ip := server.IP
				if ip == "" {
					ip = "N/A"
				}
				created := utils.FormatDateShort(server.CreatedAt)

				rows = append(rows, []string{
					server.ID,
					name,
					server.Type,
					status,
					server.Region,
					ip,
					created,
				})
			}

			utils.PrintTable(headers, rows, opts)
			utils.PrintSuccess(fmt.Sprintf("Found %d servers", len(serversResp.Servers)), opts)

			// Show helpful tips
			if !opts.Quiet {
				fmt.Printf("\n💡 TIPS\n")
				fmt.Printf("├─ View server details: pipeops server status <server-id>\n")
				fmt.Printf("└─ Install agent: pipeops agent install\n")
			}
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	listCmd.Flags().String("workspace", "", "Workspace UUID to scope server listing (or set PIPEOPS_WORKSPACE_UUID)")
}

// GetListCmd returns the list command for registration
func GetListCmd() *cobra.Command {
	return listCmd
}
