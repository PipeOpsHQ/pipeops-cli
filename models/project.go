package models

import "time"

// Project represents a PipeOps project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id"`
}

// ProjectsResponse represents the response from the projects API
type ProjectsResponse struct {
	Projects []Project `json:"projects"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
}

// ProjectCreateRequest represents the request to create a project
type ProjectCreateRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	ServerID      string                 `json:"server_id,omitempty"`
	EnvironmentID string                 `json:"environment_id,omitempty"`
	Repository    string                 `json:"repository,omitempty"`
	Branch        string                 `json:"branch,omitempty"`
	BuildCommand  string                 `json:"build_command,omitempty"`
	StartCommand  string                 `json:"start_command,omitempty"`
	Port          int                    `json:"port,omitempty"`
	Framework     string                 `json:"framework,omitempty"`
	EnvVars       map[string]interface{} `json:"env_vars,omitempty"`
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
