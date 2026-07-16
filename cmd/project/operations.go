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
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(envVars)
		}
		rows := make([][]string, 0, len(envVars))
		for _, envVar := range envVars {
			rows = append(rows, []string{envVar.Key, envVar.Value})
		}
		utils.PrintTable([]string{"KEY", "VALUE"}, rows, opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var envSetCmd = &cobra.Command{
	Use:   "set <project-id> KEY=value [KEY=value...]",
	Short: "Set project environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		envVars, _, err := parseEnvPairs(args[1:])
		if err != nil {
			return err
		}
		updated, err := client.UpdateProjectEnvVariables(args[0], envVars)
		if err != nil {
			return fmt.Errorf("set project environment variables: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(updated)
		}
		utils.PrintSuccess("Project environment variables updated", opts)
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

	envCmd.AddCommand(envGetCmd, envSetCmd)
	root.AddCommand(getCmd, updateCmd, deleteCmd, deployCmd, restartCmd, stopCmd, envCmd, deploymentsCmd, deploymentHistoryCmd)
}
