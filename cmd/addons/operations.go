package addons

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

func addonsClient(cmd *cobra.Command, opts utils.OutputOptions) (pipeops.ClientAPI, error) {
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

var deployCmd = &cobra.Command{
	Use:   "deploy <addon-id>",
	Short: "Deploy an addon",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		server, _ := cmd.Flags().GetString("server")
		workspace, _ := cmd.Flags().GetString("workspace")
		projectID, _ := cmd.Flags().GetString("project")
		configPairs, _ := cmd.Flags().GetStringArray("config")
		configMap, err := parseConfigPairs(configPairs)
		if err != nil {
			return err
		}
		deployment, err := client.DeployAddon(&sdk.DeployAddOnRequest{
			ID:        args[0],
			Server:    server,
			Workspace: workspace,
			ProjectID: projectID,
			Config:    configMap,
		})
		if err != nil {
			return fmt.Errorf("deploy addon: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Addon deployment started", opts)
		}
		return printDeployment(deployment, opts)
	},
	Args: cobra.ExactArgs(1),
}

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List addon categories",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		categories, err := client.ListAddonCategories()
		if err != nil {
			return fmt.Errorf("list addon categories: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(categories)
		}
		rows := make([][]string, 0, len(categories))
		for _, category := range categories {
			id := category.UUID
			if id == "" {
				id = category.ID
			}
			rows = append(rows, []string{id, category.Name, category.Description})
		}
		utils.PrintTable([]string{"ID", "NAME", "DESCRIPTION"}, rows, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage addon deployments",
}

var deploymentGetCmd = &cobra.Command{
	Use:   "get <deployment-id>",
	Short: "Get addon deployment details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		deployment, err := client.GetAddonDeployment(args[0])
		if err != nil {
			return fmt.Errorf("get addon deployment: %w", err)
		}
		return printDeployment(deployment, opts)
	},
	Args: cobra.ExactArgs(1),
}

var deploymentDeleteCmd = &cobra.Command{
	Use:   "delete <deployment-id>",
	Short: "Delete an addon deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			return fmt.Errorf("--force is required to delete an addon deployment")
		}
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteAddonDeployment(args[0]); err != nil {
			return fmt.Errorf("delete addon deployment: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "deployment_id": args[0]})
		}
		utils.PrintSuccess("Addon deployment deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var deploymentSessionCmd = &cobra.Command{
	Use:   "session <session-id>",
	Short: "Get addon deployment session",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		session, err := client.GetAddonDeploymentSession(args[0])
		if err != nil {
			return fmt.Errorf("get addon deployment session: %w", err)
		}
		return utils.PrintJSON(session)
	},
	Args: cobra.ExactArgs(1),
}

var deploymentConfigsCmd = &cobra.Command{
	Use:   "configs <deployment-id>",
	Short: "View addon deployment configs",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		configs, err := client.ViewAddonDeploymentConfigs(args[0])
		if err != nil {
			return fmt.Errorf("view addon deployment configs: %w", err)
		}
		return utils.PrintJSON(configs)
	},
	Args: cobra.ExactArgs(1),
}

func parseConfigPairs(pairs []string) (map[string]interface{}, error) {
	configMap := make(map[string]interface{}, len(pairs))
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid config %q; expected KEY=value", pair)
		}
		configMap[key] = value
	}
	return configMap, nil
}

func printDeployment(deployment *models.AddonDeployment, opts utils.OutputOptions) error {
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(deployment)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"ID", deployment.ID},
		{"Name", deployment.Name},
		{"Status", deployment.Status},
		{"Category", deployment.Category},
		{"Environment", deployment.Environment},
		{"Version", deployment.Version},
		{"URL", deployment.DeploymentURL},
	}, opts)
	return nil
}

func init() {
	deployCmd.Flags().String("server", "", "Server/cluster UUID")
	deployCmd.Flags().String("workspace", "", "Workspace UUID")
	deployCmd.Flags().String("project", "", "Project ID")
	deployCmd.Flags().StringArray("config", nil, "Addon config in KEY=value form; repeatable")

	deploymentDeleteCmd.Flags().Bool("force", false, "Confirm addon deployment deletion")
	deploymentCmd.AddCommand(deploymentGetCmd, deploymentDeleteCmd, deploymentSessionCmd, deploymentConfigsCmd)
	AddonsCmd.AddCommand(deployCmd, categoriesCmd, deploymentCmd)
}
