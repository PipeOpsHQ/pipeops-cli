package workspace

import (
	"context"
	"fmt"
	"strconv"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

// listCmd represents the command to list all workspaces
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all workspaces",
	Long: `List all workspaces in your PipeOps account.

Examples:
  pipeops workspace list
  pipeops workspace ls
  pipeops workspace list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		cfg, err := config.Load()
		if err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		client := pipeops.NewClientWithConfig(cfg)

		if !utils.RequireAuth(client, opts) {
			return
		}

		utils.PrintInfo("Fetching workspaces...", opts)

		workspaces, err := client.GetWorkspaces(context.Background())
		if err != nil {
			if !utils.HandleAuthError(err, opts) {
				return
			}
			utils.HandleError(err, "Error fetching workspaces", opts)
			return
		}

		if len(workspaces) == 0 {
			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON([]interface{}{})
			} else {
				fmt.Println("No workspaces found")
			}
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(workspaces)
			return
		}

		// Fetch user info to distinguish owned vs shared
		userInfoService := auth.NewUserInfoService(cfg)
		// We can get the token from the client config or re-read it
		token := client.GetConfig().OAuth.AccessToken
		userInfo, err := userInfoService.GetUserInfo(context.Background(), token)
		
		var currentUserID string
		if err == nil && userInfo != nil {
			currentUserID = strconv.Itoa(userInfo.ID)
		} else {
			// If we can't get user info, we can't distinguish, so show all as generic list
			// Or log a warning
			if opts.Verbose {
				utils.PrintWarning(fmt.Sprintf("Could not fetch user info to separate owned/shared workspaces: %v", err), opts)
			}
		}

		var ownedWorkspaces []sdk.Workspace
		var sharedWorkspaces []sdk.Workspace

		for _, ws := range workspaces {
			if currentUserID != "" && ws.OwnerID == currentUserID {
				ownedWorkspaces = append(ownedWorkspaces, ws)
			} else {
				sharedWorkspaces = append(sharedWorkspaces, ws)
			}
		}

		// If we couldn't determine ownership (currentUserID is empty), treat all as shared/generic list 
		// but since we want to show *something*, let's just dump them if we failed.
		// However, typically `GetUserInfo` should succeed if `GetWorkspaces` succeeded.
		
		if currentUserID == "" {
			// Fallback to old behavior if user info fetch failed
			printWorkspaceTable(workspaces, cfg, "WORKSPACES", opts)
		} else {
			if len(ownedWorkspaces) > 0 {
				printWorkspaceTable(ownedWorkspaces, cfg, "👤 YOUR WORKSPACES", opts)
			}
			
			if len(sharedWorkspaces) > 0 {
				if len(ownedWorkspaces) > 0 {
					fmt.Println() // Add spacing between tables
				}
				printWorkspaceTable(sharedWorkspaces, cfg, "👥 SHARED WORKSPACES", opts)
			}
		}

		if cfg.Settings != nil && cfg.Settings.DefaultWorkspaceUUID != "" {
			fmt.Printf("\n💡 Current workspace: %s\n", cfg.Settings.DefaultWorkspaceUUID)
		} else {
			fmt.Printf("\n💡 TIP: Select a default workspace with: pipeops workspace select\n")
		}
	},
	Args: cobra.NoArgs,
}

func printWorkspaceTable(workspaces []sdk.Workspace, cfg *config.Config, title string, opts utils.OutputOptions) {
	fmt.Printf("%s\n", title)
	headers := []string{"UUID", "NAME", "DESCRIPTION", "CREATED"}
	var rows [][]string

	for _, ws := range workspaces {
		name := utils.TruncateString(ws.Name, 30)
		desc := utils.TruncateString(ws.Description, 40)
		created := ""
		if ws.CreatedAt != nil {
			created = utils.FormatDateShort(ws.CreatedAt.Time)
		}

		// Mark current workspace
		if cfg.Settings != nil && cfg.Settings.DefaultWorkspaceUUID == ws.UUID {
			name = "✓ " + name
		}

		rows = append(rows, []string{
			ws.UUID,
			name,
			desc,
			created,
		})
	}

	utils.PrintTable(headers, rows, opts)
}

func (w *workspaceModel) list() {
	listCmd.Flags().Bool("json", false, "Output in JSON format")
	w.rootCmd.AddCommand(listCmd)
}
