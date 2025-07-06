package pipeops

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PipeOpsHQ/pipeops-cli/libs"
	"github.com/PipeOpsHQ/pipeops-cli/models"
)

// Config represents the PipeOps configuration
type Config struct {
	Token      string `json:"token,omitempty"`
	APIBaseURL string `json:"api_base_url,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
}

// Client represents a PipeOps client
type Client struct {
	config     *Config
	httpClient libs.HttpClients
}

// NewClient creates a new PipeOps client
func NewClient() *Client {
	return &Client{
		config:     &Config{},
		httpClient: libs.NewHttpClient(),
	}
}

// LoadConfig loads the configuration from the config file
func (c *Client) LoadConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, create default config
		c.config = &Config{
			APIBaseURL: "https://api.pipeops.io",
		}
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, c.config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// SaveConfig saves the configuration to the config file
func (c *Client) SaveConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IsAuthenticated checks if the user is authenticated
func (c *Client) IsAuthenticated() bool {
	return c.config.Token != ""
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.config.Token = token
}

// GetToken returns the current authentication token
func (c *Client) GetToken() string {
	return c.config.Token
}

// VerifyToken verifies the authentication token
func (c *Client) VerifyToken() (*models.PipeOpsTokenVerificationResponse, error) {
	if c.config.Token == "" {
		return nil, fmt.Errorf("no token set")
	}

	return c.httpClient.VerifyToken(c.config.Token, "")
}

// GetConfig returns the current configuration
func (c *Client) GetConfig() *Config {
	return c.config
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".pipeops", "config.json"), nil
}

// GetProjects retrieves all projects for the authenticated user
func (c *Client) GetProjects() (*models.ProjectsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.GetProjects(c.config.Token)
}

// GetProject retrieves a specific project by ID
func (c *Client) GetProject(projectID string) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.GetProject(c.config.Token, projectID)
}

// CreateProject creates a new project
func (c *Client) CreateProject(name, description string) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	req := &models.ProjectCreateRequest{
		Name:        name,
		Description: description,
	}

	return c.httpClient.CreateProject(c.config.Token, req)
}

// UpdateProject updates an existing project
func (c *Client) UpdateProject(projectID, name, description, status string) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	req := &models.ProjectUpdateRequest{}
	if name != "" {
		req.Name = name
	}
	if description != "" {
		req.Description = description
	}
	if status != "" {
		req.Status = status
	}

	return c.httpClient.UpdateProject(c.config.Token, projectID, req)
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID string) error {
	if !c.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	return c.httpClient.DeleteProject(c.config.Token, projectID)
}

// GetLogs retrieves logs for a project or addon
func (c *Client) GetLogs(req *models.LogsRequest) (*models.LogsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.GetLogs(c.config.Token, req)
}

// StreamLogs streams logs in real-time for a project or addon
func (c *Client) StreamLogs(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error {
	if !c.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	return c.httpClient.StreamLogs(c.config.Token, req, callback)
}

// StartProxy starts a proxy session for a project or addon service
func (c *Client) StartProxy(req *models.ProxyRequest) (*models.ProxyResponse, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.StartProxy(c.config.Token, req)
}

// GetServices retrieves available services for a project or addon
func (c *Client) GetServices(projectID string, addonID string) (*models.ListServicesResponse, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.GetServices(c.config.Token, projectID, addonID)
}

// GetContainers retrieves available containers for a project or addon
func (c *Client) GetContainers(projectID string, addonID string) (*models.ListContainersResponse, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.GetContainers(c.config.Token, projectID, addonID)
}

// StartExec starts an exec session for a project or addon container
func (c *Client) StartExec(req *models.ExecRequest) (*models.ExecResponse, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.StartExec(c.config.Token, req)
}

// StartShell starts a shell session for a project or addon container
func (c *Client) StartShell(req *models.ShellRequest) (*models.ShellResponse, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	return c.httpClient.StartShell(c.config.Token, req)
}
