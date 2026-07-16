package cmd

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a PipeOps project",
	Long: `Create a new PipeOps project.

This is a convenience alias for project creation:
  pipeops create --name api --server <server-id> --environment <env-id>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load configuration: %w", err)
		}
		client := pipeops.NewClientWithConfigFunc(cfg)
		if !utils.RequireAuth(client, opts) {
			return nil
		}
		req, err := topLevelProjectCreateRequest(cmd)
		if err != nil {
			return err
		}
		project, err := client.CreateProject(req)
		if err != nil {
			return fmt.Errorf("create project: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(project)
		}
		utils.PrintSuccess("Project created", opts)
		utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
			{"ID", project.ID},
			{"Name", project.Name},
			{"Status", project.Status},
			{"Description", project.Description},
		}, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

func topLevelProjectCreateRequest(cmd *cobra.Command) (*models.ProjectCreateRequest, error) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	serverID, _ := cmd.Flags().GetString("server")
	environmentID, _ := cmd.Flags().GetString("environment")
	repository, _ := cmd.Flags().GetString("repository")
	branch, _ := cmd.Flags().GetString("branch")
	buildCommand, _ := cmd.Flags().GetString("build-command")
	startCommand, _ := cmd.Flags().GetString("start-command")
	framework, _ := cmd.Flags().GetString("framework")
	port, _ := cmd.Flags().GetInt("port")
	envPairs, _ := cmd.Flags().GetStringArray("env")

	envVars := make(map[string]interface{}, len(envPairs))
	for _, pair := range envPairs {
		key, value, ok := strings.Cut(pair, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid env var %q; expected KEY=value", pair)
		}
		envVars[key] = value
	}

	return &models.ProjectCreateRequest{
		Name:          name,
		Description:   description,
		ServerID:      serverID,
		EnvironmentID: environmentID,
		Repository:    repository,
		Branch:        branch,
		BuildCommand:  buildCommand,
		StartCommand:  startCommand,
		Port:          port,
		Framework:     framework,
		EnvVars:       envVars,
	}, nil
}

func init() {
	createCmd.Flags().String("name", "", "Project name")
	createCmd.Flags().String("description", "", "Project description")
	createCmd.Flags().String("server", "", "Server/cluster UUID")
	createCmd.Flags().String("environment", "", "Environment UUID")
	createCmd.Flags().String("repository", "", "Repository URL or identifier")
	createCmd.Flags().String("branch", "", "Repository branch")
	createCmd.Flags().String("build-command", "", "Build command")
	createCmd.Flags().String("start-command", "", "Start command")
	createCmd.Flags().Int("port", 0, "Application port")
	createCmd.Flags().String("framework", "", "Framework name")
	createCmd.Flags().StringArray("env", nil, "Environment variable in KEY=value form; repeatable")
	_ = createCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(createCmd)
}
