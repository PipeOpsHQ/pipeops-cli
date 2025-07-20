package pipeops

import (
	"errors"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/libs"
	"github.com/PipeOpsHQ/pipeops-cli/models"
)

// Client represents the PipeOps client
type Client struct {
	httpClient libs.HttpClients
	config     *config.Config
}

// NewClient creates a new PipeOps client
func NewClient() *Client {
	return &Client{
		httpClient: libs.NewHttpClient(),
		config:     config.DefaultConfig(),
	}
}

// NewClientWithConfig creates a new PipeOps client with the provided configuration
func NewClientWithConfig(cfg *config.Config) *Client {
	baseURL := cfg.OAuth.BaseURL
	if baseURL == "" {
		baseURL = config.GetAPIURL()
	}

	return &Client{
		httpClient: libs.NewHttpClientWithURL(baseURL),
		config:     cfg,
	}
}

// LoadConfig loads the configuration from the config file
func (c *Client) LoadConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c.config = cfg
	return nil
}

// SaveConfig saves the configuration to the config file
func (c *Client) SaveConfig() error {
	return config.Save(c.config)
}

// GetConfig returns the configuration
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// IsAuthenticated checks if the user is authenticated
func (c *Client) IsAuthenticated() bool {
	return c.config.IsAuthenticated()
}

// GetToken returns the authentication token
func (c *Client) GetToken() string {
	if c.config.OAuth == nil {
		return ""
	}
	return c.config.OAuth.AccessToken
}

// GetOperatorID returns the operator ID
func (c *Client) GetOperatorID() string {
	// This could be stored in the config or derived from the token
	return ""
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	if c.config.OAuth == nil {
		c.config.OAuth = &config.OAuthConfig{}
	}
	c.config.OAuth.AccessToken = token
}

// SetOperatorID sets the operator ID
func (c *Client) SetOperatorID(operatorID string) {
	// This could be stored in the config if needed
	// For now, just skip since it's not used in the current OAuth config
}

// VerifyToken verifies the authentication token
func (c *Client) VerifyToken() (*models.PipeOpsTokenVerificationResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.VerifyToken(c.GetToken(), c.GetOperatorID())
}

// GetProjects retrieves all projects
func (c *Client) GetProjects() (*models.ProjectsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetProjects(c.GetToken())
}

// GetProject retrieves a specific project
func (c *Client) GetProject(projectID string) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetProject(c.GetToken(), projectID)
}

// CreateProject creates a new project
func (c *Client) CreateProject(req *models.ProjectCreateRequest) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.CreateProject(c.GetToken(), req)
}

// UpdateProject updates a project
func (c *Client) UpdateProject(projectID string, req *models.ProjectUpdateRequest) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.UpdateProject(c.GetToken(), projectID, req)
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	return c.httpClient.DeleteProject(c.GetToken(), projectID)
}

// GetLogs retrieves project logs
func (c *Client) GetLogs(req *models.LogsRequest) (*models.LogsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetLogs(c.GetToken(), req)
}

// StreamLogs streams project logs
func (c *Client) StreamLogs(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	return c.httpClient.StreamLogs(c.GetToken(), req, callback)
}

// GetServices retrieves services for a project
func (c *Client) GetServices(projectID string, addonID string) (*models.ListServicesResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetServices(c.GetToken(), projectID, addonID)
}

// StartProxy starts a proxy session
func (c *Client) StartProxy(req *models.ProxyRequest) (*models.ProxyResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.StartProxy(c.GetToken(), req)
}

// GetContainers retrieves containers for a project
func (c *Client) GetContainers(projectID string, addonID string) (*models.ListContainersResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetContainers(c.GetToken(), projectID, addonID)
}

// StartExec starts an exec session
func (c *Client) StartExec(req *models.ExecRequest) (*models.ExecResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.StartExec(c.GetToken(), req)
}

// StartShell starts a shell session
func (c *Client) StartShell(req *models.ShellRequest) (*models.ShellResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.StartShell(c.GetToken(), req)
}

// GetAddons retrieves a list of addons
func (c *Client) GetAddons() (*models.AddonListResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetAddons(c.GetToken())
}

// GetAddon retrieves a specific addon by ID
func (c *Client) GetAddon(addonID string) (*models.Addon, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetAddon(c.GetToken(), addonID)
}

// DeployAddon deploys an addon
func (c *Client) DeployAddon(req *models.AddonDeployRequest) (*models.AddonDeployResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.DeployAddon(c.GetToken(), req)
}

// GetAddonDeployments retrieves a list of addon deployments
func (c *Client) GetAddonDeployments(projectID string) ([]models.AddonDeployment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetAddonDeployments(c.GetToken(), projectID)
}

// DeleteAddonDeployment deletes an addon deployment
func (c *Client) DeleteAddonDeployment(deploymentID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	return c.httpClient.DeleteAddonDeployment(c.GetToken(), deploymentID)
}

// GetServers retrieves all servers
func (c *Client) GetServers() (*models.ServersResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetServers(c.GetToken())
}

// GetServer retrieves a specific server by ID
func (c *Client) GetServer(serverID string) (*models.Server, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.GetServer(c.GetToken(), serverID)
}

// CreateServer creates a new server
func (c *Client) CreateServer(req *models.ServerCreateRequest) (*models.Server, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.CreateServer(c.GetToken(), req)
}

// UpdateServer updates an existing server
func (c *Client) UpdateServer(serverID string, req *models.ServerUpdateRequest) (*models.Server, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	return c.httpClient.UpdateServer(c.GetToken(), serverID, req)
}

// DeleteServer deletes a server
func (c *Client) DeleteServer(serverID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	return c.httpClient.DeleteServer(c.GetToken(), serverID)
}
