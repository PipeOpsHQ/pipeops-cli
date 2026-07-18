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

func parseEnvPairs(pairs []string) ([]models.ProjectEnvVar, error) {
	envVars := make([]models.ProjectEnvVar, 0, len(pairs))
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid env var %q; expected KEY=value", pair)
		}
		envVars = append(envVars, models.ProjectEnvVar{Key: key, Value: value})
	}
	return envVars, nil
}

// ProjectCreateRequestFromFlags builds a ProjectCreateRequest from create flags.
// Exported so the top-level `pipeops create` alias can reuse the same mapping.
func ProjectCreateRequestFromFlags(cmd *cobra.Command) (*models.ProjectCreateRequest, error) {
	return projectCreateRequestFromFlags(cmd)
}

func projectCreateRequestFromFlags(cmd *cobra.Command) (*models.ProjectCreateRequest, error) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	serverID, _ := cmd.Flags().GetString("server")
	clusterID, _ := cmd.Flags().GetString("cluster")
	environmentID, _ := cmd.Flags().GetString("environment")
	environmentName, _ := cmd.Flags().GetString("environment-name")
	repository, _ := cmd.Flags().GetString("repository")
	branch, _ := cmd.Flags().GetString("branch")
	source, _ := cmd.Flags().GetString("source")
	username, _ := cmd.Flags().GetString("username")
	buildCommand, _ := cmd.Flags().GetString("build-command")
	startCommand, _ := cmd.Flags().GetString("start-command")
	buildMethod, _ := cmd.Flags().GetString("build-method")
	framework, _ := cmd.Flags().GetString("framework")
	language, _ := cmd.Flags().GetString("language")
	port, _ := cmd.Flags().GetInt("port")
	envPairs, _ := cmd.Flags().GetStringArray("env")
	workspace, _ := cmd.Flags().GetString("workspace")
	commitURL, _ := cmd.Flags().GetString("commit-url")
	commitSha, _ := cmd.Flags().GetString("commit-sha")
	worker, _ := cmd.Flags().GetBool("worker")

	envVars, err := parseEnvPairs(envPairs)
	if err != nil {
		return nil, err
	}

	clusterUUID := strings.TrimSpace(serverID)
	if clusterUUID == "" {
		clusterUUID = strings.TrimSpace(clusterID)
	}

	return &models.ProjectCreateRequest{
		Name:               name,
		Description:        description,
		ClusterUUID:        clusterUUID,
		EnvironmentUUID:    environmentID,
		Environment:        environmentName,
		Repository:         repository,
		Branch:             branch,
		Source:             source,
		Username:           username,
		BuildCommand:       buildCommand,
		StartCommand:       startCommand,
		BuildMethod:        buildMethod,
		Port:               port,
		Framework:          framework,
		RepositoryLanguage: language,
		EnvVariables:       envVars,
		WorkspaceUUID:      workspace,
		CommitURL:          commitURL,
		CommitSha:          commitSha,
		Worker:             worker,
	}, nil
}

// AddProjectCreateFlags registers shared create flags on a command.
// Exported for the top-level `pipeops create` alias.
func AddProjectCreateFlags(cmd *cobra.Command) {
	addProjectCreateFlags(cmd)
}

func addProjectCreateFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Project name")
	cmd.Flags().String("description", "", "Project description (CLI-only; not sent on create)")
	cmd.Flags().String("server", "", "Server/cluster UUID (alias for --cluster)")
	cmd.Flags().String("cluster", "", "Cluster UUID (alias for --server)")
	cmd.Flags().String("environment", "", "Environment UUID")
	cmd.Flags().String("environment-name", "development", "Environment name/slug (e.g. development)")
	cmd.Flags().String("repository", "", "Repository URL or owner/repo")
	cmd.Flags().String("branch", "", "Repository branch")
	cmd.Flags().String("source", "github", "Source provider: github, gitlab, bitbucket, or image")
	cmd.Flags().String("username", "", "Repository owner/username (defaults from --repository)")
	cmd.Flags().String("build-command", "", "Build command")
	cmd.Flags().String("start-command", "", "Start/run command")
	cmd.Flags().String("build-method", "", "Build method (default: nodejs; use dockerfile for Docker builds)")
	cmd.Flags().Int("port", 0, "Application port (default 3000 for web projects)")
	cmd.Flags().String("framework", "", "Framework name")
	cmd.Flags().String("language", "", "Repository language (maps to repositoryLanguage)")
	cmd.Flags().StringArray("env", nil, "Environment variable in KEY=value form; repeatable")
	cmd.Flags().String("workspace", "", "Workspace UUID (defaults to configured workspace)")
	cmd.Flags().String("commit-url", "", "Commit URL")
	cmd.Flags().String("commit-sha", "", "Commit SHA")
	cmd.Flags().Bool("worker", false, "Create as a worker project (no network port)")
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
