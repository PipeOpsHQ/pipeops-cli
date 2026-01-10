package models

import "time"

// Addon represents an addon service that can be deployed
type Addon struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Version     string            `json:"version"`
	Status      string            `json:"status"`
	Image       string            `json:"image"`
	Icon        string            `json:"icon,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Ports       []int             `json:"ports,omitempty"`
	EnvVars     map[string]string `json:"env_vars,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Service represents a service in a project
type Service struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	ProjectID string `json:"project_id"`
}

// AddonListResponse represents the response when listing addons
type AddonListResponse struct {
	Addons []Addon `json:"addons"`
	Total  int     `json:"total"`
}

// AddonDeployRequest represents a request to deploy an addon
type AddonDeployRequest struct {
	ID        string            `json:"id"`
	Server    string            `json:"Server"`
	Workspace string            `json:"Workspace"`
	ProjectID string            `json:"project_id,omitempty"`
	Name      string            `json:"name,omitempty"`
	EnvVars   map[string]string `json:"env_vars,omitempty"`
	Config    map[string]string `json:"config,omitempty"`
}

// AddonDeployResponse represents the response when deploying an addon
type AddonDeployResponse struct {
	DeploymentID string `json:"deployment_id"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

// AddonDeployment represents a deployed addon instance
type AddonDeployment struct {
	ID        string            `json:"id"`
	AddonID   string            `json:"addon_id"`
	ProjectID string            `json:"project_id"`
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	URL       string            `json:"url,omitempty"`
	EnvVars   map[string]string `json:"env_vars,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// AddonDeploymentsResponse represents the response when listing addon deployments
type AddonDeploymentsResponse struct {
	Deployments []AddonDeployment `json:"deployments"`
	Total       int               `json:"total"`
}
