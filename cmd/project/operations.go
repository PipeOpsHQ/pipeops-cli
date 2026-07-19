package project

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <project-id>",
	Short: "Get project details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		project, err := client.GetProject(args[0])
		if err != nil {
			return fmt.Errorf("get project: %w", err)
		}
		printProject(project, opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var updateCmd = &cobra.Command{
	Use:   "update <project-id>",
	Short: "Update project configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		project, err := client.UpdateProject(args[0], projectUpdateRequestFromFlags(cmd))
		if err != nil {
			return fmt.Errorf("update project: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Project updated", opts)
		}
		printProject(project, opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var deleteCmd = &cobra.Command{
	Use:   "delete <project-id>",
	Short: "Delete a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			return fmt.Errorf("--force is required to delete a project")
		}
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteProject(args[0]); err != nil {
			return fmt.Errorf("delete project: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "project_id": args[0]})
		}
		utils.PrintSuccess("Project deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var deployCmd = &cobra.Command{
	Use:   "deploy <project-id>",
	Short: "Trigger a project deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProjectAction(cmd, args[0], "deployed", func(client interface{ DeployProject(string) error }) error {
			return client.DeployProject(args[0])
		})
	},
	Args: cobra.ExactArgs(1),
}

var restartCmd = &cobra.Command{
	Use:   "restart <project-id>",
	Short: "Restart a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProjectAction(cmd, args[0], "restarted", func(client interface{ RestartProject(string) error }) error {
			return client.RestartProject(args[0])
		})
	},
	Args: cobra.ExactArgs(1),
}

var stopCmd = &cobra.Command{
	Use:   "stop <project-id>",
	Short: "Stop a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProjectAction(cmd, args[0], "stopped", func(client interface{ StopProject(string) error }) error {
			return client.StopProject(args[0])
		})
	},
	Args: cobra.ExactArgs(1),
}

func runProjectAction[T any](cmd *cobra.Command, projectID, status string, action func(T) error) error {
	opts := utils.GetOutputOptions(cmd)
	client, err := authenticatedClient(cmd, opts)
	if err != nil || client == nil {
		return err
	}
	typed, ok := any(client).(T)
	if !ok {
		return fmt.Errorf("client does not support project action")
	}
	if err := action(typed); err != nil {
		return fmt.Errorf("project %s: %w", status, err)
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(map[string]string{"status": status, "project_id": projectID})
	}
	utils.PrintSuccess(fmt.Sprintf("Project %s", status), opts)
	return nil
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage project environment variables",
}

var envGetCmd = &cobra.Command{
	Use:   "get <project-id>",
	Short: "Get project environment variables",
	Long: `Get project environment variables.

Values are masked by default so secrets are not printed to the terminal or
shell history. Pass --reveal to show plaintext values (use only on a trusted
terminal; prefer --json with --reveal for scripting).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		envVars, err := client.GetProjectEnvVariables(args[0])
		if err != nil {
			return fmt.Errorf("get project environment variables: %w", err)
		}
		reveal, _ := cmd.Flags().GetBool("reveal")
		display := maskEnvVariables(envVars, reveal)
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(display)
		}
		rows := make([][]string, 0, len(display))
		for _, envVar := range display {
			rows = append(rows, []string{envVar.Key, envVar.Value})
		}
		utils.PrintTable([]string{"KEY", "VALUE"}, rows, opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

// maskEnvVariables returns a copy of env vars with values redacted unless reveal is true.
func maskEnvVariables(envVars []sdk.EnvVariable, reveal bool) []sdk.EnvVariable {
	if reveal {
		return envVars
	}
	out := make([]sdk.EnvVariable, 0, len(envVars))
	for _, envVar := range envVars {
		out = append(out, sdk.EnvVariable{
			Key:   envVar.Key,
			Value: maskSecretValue(envVar.Value),
		})
	}
	return out
}

// maskSecretValue redacts a secret for display while preserving a hint of length.
func maskSecretValue(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 4 {
		return "****"
	}
	return "****" + value[len(value)-2:]
}

var envSetCmd = &cobra.Command{
	Use:   "set <project-id> KEY=value [KEY=value...]",
	Short: "Set project environment variables",
	Long: `Set project environment variables.

By default keys are merged into existing envs (prefer-client: client values win,
other keys kept). Pass --replace for a full replace of the entire env set
(dashboard-style). PORT is injected server-side from network settings when missing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		parsed, err := parseEnvPairs(args[1:])
		if err != nil {
			return err
		}
		envVars := make([]sdk.EnvVariable, 0, len(parsed))
		for _, ev := range parsed {
			envVars = append(envVars, sdk.EnvVariable{Key: ev.Key, Value: ev.Value})
		}
		// Prefer-client default: merge=true. --replace forces full replace.
		replace, _ := cmd.Flags().GetBool("replace")
		merge := !replace
		if flag := cmd.Flags().Lookup("merge"); flag != nil && flag.Changed {
			merge, _ = cmd.Flags().GetBool("merge")
		}
		updated, err := client.UpdateProjectEnvVariables(args[0], envVars, merge)
		if err != nil {
			return fmt.Errorf("set project environment variables: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(updated)
		}
		if merge {
			utils.PrintSuccess("Project environment variables merged", opts)
		} else {
			utils.PrintSuccess("Project environment variables replaced", opts)
		}
		return nil
	},
	Args: cobra.MinimumNArgs(2),
}

var deploymentsCmd = &cobra.Command{
	Use:   "deployments <project-id>",
	Short: "List project deployments",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.ListProjectDeployments(args[0], &sdk.ProjectDeploymentListOptions{
			FilterBy: cmd.Flag("filter").Value.String(),
			Page:     intFlag(cmd, "page", 1),
			Limit:    intFlag(cmd, "limit", 20),
		})
		if err != nil {
			return fmt.Errorf("list project deployments: %w", err)
		}
		return printDeploymentRecords(resp.Data, opts)
	},
	Args: cobra.ExactArgs(1),
}

var deploymentHistoryCmd = &cobra.Command{
	Use:   "deployment-history <project-id>",
	Short: "List project deployment history",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.ListProjectDeploymentHistory(args[0], &sdk.ProjectDeploymentHistoryOptions{
			Page:  intFlag(cmd, "page", 1),
			Limit: intFlag(cmd, "limit", 20),
		})
		if err != nil {
			return fmt.Errorf("list project deployment history: %w", err)
		}
		return printDeploymentRecords(resp.Data, opts)
	},
	Args: cobra.ExactArgs(1),
}

func printDeploymentRecords(records []sdk.ProjectDeploymentRecord, opts utils.OutputOptions) error {
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(records)
	}
	rows := make([][]string, 0, len(records))
	for _, record := range records {
		rows = append(rows, []string{
			deploymentRecordValue(record, "id", "ID", "uuid", "UUID"),
			deploymentRecordValue(record, "name", "Name", "deployment_name", "DeploymentName"),
			deploymentRecordValue(record, "status", "Status"),
			deploymentRecordValue(record, "created_at", "CreatedAt", "createdAt"),
		})
	}
	utils.PrintTable([]string{"ID", "NAME", "STATUS", "CREATED"}, rows, opts)
	return nil
}

func registerOperationCommands(root *cobra.Command) {
	addProjectUpdateFlags(updateCmd)
	deleteCmd.Flags().Bool("force", false, "Confirm project deletion")
	deploymentsCmd.Flags().String("filter", "", "Deployment filter")
	deploymentsCmd.Flags().Int("page", 1, "Page number")
	deploymentsCmd.Flags().Int("limit", 20, "Page size")
	deploymentHistoryCmd.Flags().Int("page", 1, "Page number")
	deploymentHistoryCmd.Flags().Int("limit", 20, "Page size")
	envGetCmd.Flags().Bool("reveal", false, "Show plaintext secret values (default: masked)")
	envSetCmd.Flags().Bool("merge", true, "Merge keys into existing envs (default true; prefer-client)")
	envSetCmd.Flags().Bool("replace", false, "Full-replace entire env set instead of merging")

	envCmd.AddCommand(envGetCmd, envSetCmd)
	root.AddCommand(getCmd, updateCmd, deleteCmd, deployCmd, restartCmd, stopCmd, envCmd, deploymentsCmd, deploymentHistoryCmd)
}
