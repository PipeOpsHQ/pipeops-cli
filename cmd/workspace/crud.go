package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

func workspaceClient(cmd *cobra.Command, opts utils.OutputOptions) (pipeops.ClientAPI, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}
	client := pipeops.NewClientWithConfig(cfg)
	if !utils.RequireAuth(client, opts) {
		return nil, nil
	}
	return client, nil
}

var getCmd = &cobra.Command{
	Use:   "get <workspace-id>",
	Short: "Get workspace details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := workspaceClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		workspace, err := client.GetWorkspace(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("get workspace: %w", err)
		}
		return printWorkspace(workspace, opts)
	},
	Args: cobra.ExactArgs(1),
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := workspaceClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		teamID, _ := cmd.Flags().GetString("team")
		workspace, err := client.CreateWorkspace(context.Background(), &sdk.CreateWorkspaceRequest{
			Name:        name,
			Description: description,
			TeamID:      teamID,
		})
		if err != nil {
			return fmt.Errorf("create workspace: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Workspace created", opts)
		}
		return printWorkspace(workspace, opts)
	},
	Args: cobra.NoArgs,
}

var updateCmd = &cobra.Command{
	Use:   "update <workspace-id>",
	Short: "Update a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := workspaceClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		workspace, err := client.UpdateWorkspace(context.Background(), args[0], &sdk.UpdateWorkspaceRequest{
			Name:        name,
			Description: description,
		})
		if err != nil {
			return fmt.Errorf("update workspace: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Workspace updated", opts)
		}
		return printWorkspace(workspace, opts)
	},
	Args: cobra.ExactArgs(1),
}

var deleteCmd = &cobra.Command{
	Use:   "delete <workspace-id>",
	Short: "Delete a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			return fmt.Errorf("--force is required to delete a workspace")
		}
		client, err := workspaceClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteWorkspace(context.Background(), args[0]); err != nil {
			return fmt.Errorf("delete workspace: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "workspace_id": args[0]})
		}
		utils.PrintSuccess("Workspace deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func printWorkspace(workspace *sdk.Workspace, opts utils.OutputOptions) error {
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(workspace)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"ID", workspace.ID},
		{"UUID", workspace.UUID},
		{"Name", workspace.Name},
		{"Description", workspace.Description},
		{"Owner ID", workspace.OwnerID},
		{"Created", utils.FormatDate(timestampPtr(workspace.CreatedAt))},
	}, opts)
	return nil
}

func timestampPtr(ts *sdk.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.Time
}

func (w *workspaceModel) crud() {
	createCmd.Flags().String("name", "", "Workspace name")
	createCmd.Flags().String("description", "", "Workspace description")
	createCmd.Flags().String("team", "", "Team ID")
	_ = createCmd.MarkFlagRequired("name")

	updateCmd.Flags().String("name", "", "Workspace name")
	updateCmd.Flags().String("description", "", "Workspace description")
	deleteCmd.Flags().Bool("force", false, "Confirm workspace deletion")

	w.rootCmd.AddCommand(getCmd, createCmd, updateCmd, deleteCmd)
}
