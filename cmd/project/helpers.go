package project

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

func authenticatedClient(cmd *cobra.Command, opts utils.OutputOptions) (pipeops.ClientAPI, error) {
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

func printProject(project *models.Project, opts utils.OutputOptions) {
	if opts.Format == utils.OutputFormatJSON {
		_ = utils.PrintJSON(project)
		return
	}
	rows := [][]string{
		{"ID", project.ID},
		{"Name", project.Name},
		{"Status", project.Status},
		{"URL", project.URL},
		{"Description", project.Description},
		{"Created", utils.FormatDate(project.CreatedAt)},
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, rows, opts)
}

func parseEnvPairs(pairs []string) ([]sdk.EnvVariable, map[string]interface{}, error) {
	envVars := make([]sdk.EnvVariable, 0, len(pairs))
	envMap := make(map[string]interface{}, len(pairs))
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, nil, fmt.Errorf("invalid env var %q; expected KEY=value", pair)
		}
		envVars = append(envVars, sdk.EnvVariable{Key: key, Value: value})
		envMap[key] = value
	}
	return envVars, envMap, nil
}

func projectCreateRequestFromFlags(cmd *cobra.Command) (*models.ProjectCreateRequest, error) {
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

	_, envMap, err := parseEnvPairs(envPairs)
	if err != nil {
		return nil, err
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
		EnvVars:       envMap,
	}, nil
}

func addProjectCreateFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Project name")
	cmd.Flags().String("description", "", "Project description")
	cmd.Flags().String("server", "", "Server/cluster UUID")
	cmd.Flags().String("environment", "", "Environment UUID")
	cmd.Flags().String("repository", "", "Repository URL or identifier")
	cmd.Flags().String("branch", "", "Repository branch")
	cmd.Flags().String("build-command", "", "Build command")
	cmd.Flags().String("start-command", "", "Start command")
	cmd.Flags().Int("port", 0, "Application port")
	cmd.Flags().String("framework", "", "Framework name")
	cmd.Flags().StringArray("env", nil, "Environment variable in KEY=value form; repeatable")
	_ = cmd.MarkFlagRequired("name")
}

func addProjectUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Project name")
	cmd.Flags().String("description", "", "Project description")
	cmd.Flags().String("build-command", "", "Build command")
	cmd.Flags().String("start-command", "", "Start command")
	cmd.Flags().Int("port", 0, "Application port")
}

func projectUpdateRequestFromFlags(cmd *cobra.Command) *models.ProjectUpdateRequest {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	buildCommand, _ := cmd.Flags().GetString("build-command")
	startCommand, _ := cmd.Flags().GetString("start-command")
	port, _ := cmd.Flags().GetInt("port")
	return &models.ProjectUpdateRequest{
		Name:         name,
		Description:  description,
		BuildCommand: buildCommand,
		StartCommand: startCommand,
		Port:         port,
	}
}

func deploymentRecordValue(record sdk.ProjectDeploymentRecord, keys ...string) string {
	for _, key := range keys {
		if value, ok := record[key]; ok && value != nil {
			return fmt.Sprintf("%v", value)
		}
	}
	return ""
}

func intFlag(cmd *cobra.Command, name string, fallback int) int {
	value, _ := cmd.Flags().GetInt(name)
	if value == 0 {
		return fallback
	}
	return value
}

func parseBoolPointer(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
