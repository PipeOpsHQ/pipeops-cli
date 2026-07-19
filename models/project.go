package models

import "time"

// Project represents a PipeOps project
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	// URL is the public app URL (agent-managed, PKS LB, or NonPks domain).
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
}

// ProjectsResponse represents the response from the projects API
type ProjectsResponse struct {
	Projects []Project `json:"projects"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
}

// ProjectEnvVar is a key/value environment variable for project create.
type ProjectEnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ProjectCreateRequest is the CLI-level request for creating a project.
// CreateProject maps this onto the control-plane CreateProjectRequest shape
// (clusterUUID, environment_uuid, buildSettings, envVariables, …).
type ProjectCreateRequest struct {
	Name string `json:"name"`
	// Description is accepted by the CLI for local display/compat but is not
	// part of POST /project/create.
	Description string `json:"description,omitempty"`

	// ClusterUUID is the server/cluster UUID (CLI flags: --server, --cluster).
	ClusterUUID string `json:"cluster_uuid,omitempty"`
	// EnvironmentUUID is the environment UUID (CLI flag: --environment).
	EnvironmentUUID string `json:"environment_uuid,omitempty"`
	// Environment is the environment name/slug (CLI flag: --environment-name).
	// Defaults to "development" when empty.
	Environment string `json:"environment,omitempty"`

	Repository string `json:"repository,omitempty"`
	Branch     string `json:"branch,omitempty"`
	// Source is the VCS source: github | gitlab | bitbucket | image.
	// Defaults to "github" when empty.
	Source   string `json:"source,omitempty"`
	Username string `json:"username,omitempty"`

	BuildCommand string `json:"build_command,omitempty"`
	// StartCommand maps to buildSettings.runCommand.
	StartCommand string `json:"start_command,omitempty"`
	// BuildMethod maps to buildSettings.buildMethod (e.g. nodejs, dockerfile).
	// Defaults to "nodejs" when empty (or "dockerfile" when language suggests it).
	BuildMethod string `json:"build_method,omitempty"`
	// Port is the HTTP service port. Defaults to 3000 for non-worker projects.
	Port int `json:"port,omitempty"`

	Framework          string `json:"framework,omitempty"`
	RepositoryLanguage string `json:"repository_language,omitempty"`

	// EnvVariables are sent as envVariables[]. Prefer this over EnvVars.
	EnvVariables []ProjectEnvVar `json:"env_variables,omitempty"`
	// EnvVars is a legacy map form kept for callers that still set KEY→value.
	// CreateProject merges it into EnvVariables when EnvVariables is empty.
	EnvVars map[string]interface{} `json:"env_vars,omitempty"`

	WorkspaceUUID string `json:"workspace_uuid,omitempty"`
	CommitURL     string `json:"commit_url,omitempty"`
	CommitSha     string `json:"commit_sha,omitempty"`

	// Worker sets buildSettings.worker=true and omits networkSettings.
	Worker bool `json:"worker,omitempty"`
}

// ProjectUpdateRequest represents the request to update a project
type ProjectUpdateRequest struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Status       string `json:"status,omitempty"`
	BuildCommand string `json:"build_command,omitempty"`
	StartCommand string `json:"start_command,omitempty"`
	Port         int    `json:"port,omitempty"`
}
