package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

var environmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"env", "environments"},
	Short:   "Manage environments",
}

var environmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		envs, err := client.ListEnvironments(context.Background())
		if err != nil {
			return fmt.Errorf("list environments: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(envs)
		}
		rows := make([][]string, 0, len(envs))
		for _, env := range envs {
			rows = append(rows, []string{envID(env), env.Name, env.WorkspaceID, sdkTime(env.CreatedAt)})
		}
		utils.PrintTable([]string{"ID", "NAME", "WORKSPACE", "CREATED"}, rows, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

var environmentGetCmd = &cobra.Command{
	Use:   "get <environment-id>",
	Short: "Get environment details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		env, err := client.GetEnvironment(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("get environment: %w", err)
		}
		return printEnvironment(env, opts)
	},
	Args: cobra.ExactArgs(1),
}

var environmentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		workspace, _ := cmd.Flags().GetString("workspace")
		cluster, _ := cmd.Flags().GetString("cluster")
		envPairs, _ := cmd.Flags().GetStringArray("env")
		envVars, err := parseSDKEnvPairs(envPairs)
		if err != nil {
			return err
		}
		env, err := client.CreateEnvironment(context.Background(), &sdk.CreateEnvironmentRequest{
			Name:          name,
			WorkspaceUUID: workspace,
			ClusterUUID:   cluster,
			EnvVariables:  envVars,
		})
		if err != nil {
			return fmt.Errorf("create environment: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Environment created", opts)
		}
		return printEnvironment(env, opts)
	},
	Args: cobra.NoArgs,
}

var environmentUpdateCmd = &cobra.Command{
	Use:   "update <environment-id>",
	Short: "Update an environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		env, err := client.UpdateEnvironment(context.Background(), args[0], &sdk.UpdateEnvironmentRequest{Name: name})
		if err != nil {
			return fmt.Errorf("update environment: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Environment updated", opts)
		}
		return printEnvironment(env, opts)
	},
	Args: cobra.ExactArgs(1),
}

var environmentDeleteCmd = &cobra.Command{
	Use:   "delete <environment-id>",
	Short: "Delete an environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			return fmt.Errorf("--force is required to delete an environment")
		}
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteEnvironment(context.Background(), args[0]); err != nil {
			return fmt.Errorf("delete environment: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "environment_id": args[0]})
		}
		utils.PrintSuccess("Environment deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var environmentVarsCmd = &cobra.Command{
	Use:   "vars",
	Short: "Manage environment variables",
}

var environmentVarsSetCmd = &cobra.Command{
	Use:   "set <environment-id> KEY=value [KEY=value...]",
	Short: "Set environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		envVars, err := parseSDKEnvPairs(args[1:])
		if err != nil {
			return err
		}
		if err := client.SetEnvironmentVariables(context.Background(), args[0], envVars); err != nil {
			return fmt.Errorf("set environment variables: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]interface{}{"status": "updated", "environment_id": args[0], "env_variables": envVars})
		}
		utils.PrintSuccess("Environment variables updated", opts)
		return nil
	},
	Args: cobra.MinimumNArgs(2),
}

func rootClient(opts utils.OutputOptions) (pipeops.ClientAPI, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}
	client := pipeops.NewClientWithConfigFunc(cfg)
	if !utils.RequireAuth(client, opts) {
		return nil, nil
	}
	return client, nil
}

func parseSDKEnvPairs(pairs []string) ([]sdk.EnvVariable, error) {
	envVars := make([]sdk.EnvVariable, 0, len(pairs))
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid env var %q; expected KEY=value", pair)
		}
		envVars = append(envVars, sdk.EnvVariable{Key: key, Value: value})
	}
	return envVars, nil
}

func printEnvironment(env *sdk.Environment, opts utils.OutputOptions) error {
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(env)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"ID", envID(*env)},
		{"Name", env.Name},
		{"Workspace ID", env.WorkspaceID},
		{"Created", sdkTime(env.CreatedAt)},
		{"Updated", sdkTime(env.UpdatedAt)},
	}, opts)
	return nil
}

func envID(env sdk.Environment) string {
	if env.UUID != "" {
		return env.UUID
	}
	return env.ID
}

func sdkTime(ts *sdk.Timestamp) string {
	if ts == nil {
		return "N/A"
	}
	return utils.FormatDate(ts.Time)
}

func init() {
	environmentCreateCmd.Flags().String("name", "", "Environment name")
	environmentCreateCmd.Flags().String("workspace", "", "Workspace UUID")
	environmentCreateCmd.Flags().String("cluster", "", "Cluster UUID")
	environmentCreateCmd.Flags().StringArray("env", nil, "Environment variable in KEY=value form; repeatable")
	_ = environmentCreateCmd.MarkFlagRequired("name")

	environmentUpdateCmd.Flags().String("name", "", "Environment name")
	_ = environmentUpdateCmd.MarkFlagRequired("name")
	environmentDeleteCmd.Flags().Bool("force", false, "Confirm environment deletion")

	environmentVarsCmd.AddCommand(environmentVarsSetCmd)
	environmentCmd.AddCommand(
		environmentListCmd,
		environmentGetCmd,
		environmentCreateCmd,
		environmentUpdateCmd,
		environmentDeleteCmd,
		environmentVarsCmd,
	)
	rootCmd.AddCommand(environmentCmd)
}
