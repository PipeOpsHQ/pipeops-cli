package pipeops

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

// Client represents the PipeOps client wrapping the Go SDK
type Client struct {
	sdkClient *sdk.Client
	config    *config.Config
}

// NewClient creates a new PipeOps client
func NewClient() *Client {
	cfg := config.DefaultConfig()
	baseURL := config.GetAPIURL()

	sdkClient, err := sdk.NewClient(baseURL)
	if err != nil {
		// Fallback to default if URL parsing fails
		sdkClient, _ = sdk.NewClient("")
	}

	return &Client{
		sdkClient: sdkClient,
		config:    cfg,
	}
}

// NewClientWithConfig creates a new PipeOps client with the provided configuration
func NewClientWithConfig(cfg *config.Config) *Client {
	baseURL := cfg.OAuth.BaseURL
	if baseURL == "" {
		baseURL = config.GetAPIURL()
	}

	sdkClient, err := sdk.NewClient(baseURL,
		sdk.WithTimeout(30*time.Second),
		sdk.WithMaxRetries(3),
	)
	if err != nil {
		// Fallback to default if URL parsing fails
		sdkClient, _ = sdk.NewClient("")
	}

	// Set the access token if available
	if cfg.OAuth != nil && cfg.OAuth.AccessToken != "" {
		sdkClient.SetToken(cfg.OAuth.AccessToken)
	}

	return &Client{
		sdkClient: sdkClient,
		config:    cfg,
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
	// Token verification is implicit in SDK through API calls
	// We'll use user settings endpoint as a verification method
	ctx := context.Background()
	_, _, err := c.sdkClient.Users.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return &models.PipeOpsTokenVerificationResponse{
		Valid: true,
	}, nil
}

// GetProjects retrieves all projects
func (c *Client) GetProjects() (*models.ProjectsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Projects.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to CLI models
	projects := make([]models.Project, len(resp.Data.Projects))
	for i, p := range resp.Data.Projects {
		projects[i] = models.Project{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt.Time,
			UpdatedAt:   p.UpdatedAt.Time,
		}
	}

	return &models.ProjectsResponse{
		Projects: projects,
		Total:    len(projects),
		Page:     1,
		PerPage:  len(projects),
	}, nil
}

// GetProject retrieves a specific project
func (c *Client) GetProject(projectID string) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Projects.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &models.Project{
		ID:          resp.Data.Project.ID,
		Name:        resp.Data.Project.Name,
		Description: resp.Data.Project.Description,
		Status:      resp.Data.Project.Status,
		CreatedAt:   resp.Data.Project.CreatedAt.Time,
		UpdatedAt:   resp.Data.Project.UpdatedAt.Time,
	}, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(req *models.ProjectCreateRequest) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	createReq := &sdk.CreateProjectRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	resp, _, err := c.sdkClient.Projects.Create(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &models.Project{
		ID:          resp.Data.Project.ID,
		Name:        resp.Data.Project.Name,
		Description: resp.Data.Project.Description,
		Status:      resp.Data.Project.Status,
		CreatedAt:   resp.Data.Project.CreatedAt.Time,
		UpdatedAt:   resp.Data.Project.UpdatedAt.Time,
	}, nil
}

// UpdateProject updates a project
func (c *Client) UpdateProject(projectID string, req *models.ProjectUpdateRequest) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	updateReq := &sdk.UpdateProjectRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	resp, _, err := c.sdkClient.Projects.Update(ctx, projectID, updateReq)
	if err != nil {
		return nil, err
	}

	return &models.Project{
		ID:          resp.Data.Project.ID,
		Name:        resp.Data.Project.Name,
		Description: resp.Data.Project.Description,
		Status:      resp.Data.Project.Status,
		CreatedAt:   resp.Data.Project.CreatedAt.Time,
		UpdatedAt:   resp.Data.Project.UpdatedAt.Time,
	}, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()
	_, err := c.sdkClient.Projects.Delete(ctx, projectID)
	return err
}

// GetLogs retrieves project logs
func (c *Client) GetLogs(req *models.LogsRequest) (*models.LogsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	
	// Build SDK request options
	opts := &sdk.LogsOptions{
		Limit: req.Limit,
	}

	resp, _, err := c.sdkClient.Projects.GetLogs(ctx, req.ProjectID, opts)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to CLI models
	logs := make([]models.LogEntry, 0)
	for _, logMap := range resp.Data.Logs {
		// Convert map to log entry
		entry := models.LogEntry{
			Message: fmt.Sprintf("%v", logMap),
		}
		logs = append(logs, entry)
	}

	return &models.LogsResponse{
		Logs:       logs,
		TotalCount: len(logs),
		HasMore:    false,
	}, nil
}

// StreamLogs streams project logs
func (c *Client) StreamLogs(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()
	
	// Build SDK request options
	opts := &sdk.LogsOptions{
		Limit: req.Limit,
	}

	// For now, just fetch logs (SDK may not have streaming support yet)
	resp, _, err := c.sdkClient.Projects.TailLogs(ctx, req.ProjectID, opts)
	if err != nil {
		return err
	}

	// Convert and callback with each log entry
	for _, logMap := range resp.Data.Logs {
		streamEntry := &models.StreamLogEntry{
			LogEntry: models.LogEntry{
				Message: fmt.Sprintf("%v", logMap),
			},
		}
		if err := callback(streamEntry); err != nil {
			return err
		}
	}

	return nil
}

// GetServices retrieves services for a project
func (c *Client) GetServices(projectID string, addonID string) (*models.ListServicesResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// Services may not be directly available in SDK yet
	// Return empty list for now
	return &models.ListServicesResponse{
		Services: []models.ServiceInfo{},
	}, nil
}

// StartProxy starts a proxy session
func (c *Client) StartProxy(req *models.ProxyRequest) (*models.ProxyResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// Proxy functionality may need direct HTTP implementation
	// or specific SDK support - keeping as-is for now
	return nil, errors.New("proxy not yet implemented with SDK")
}

// GetContainers retrieves containers for a project
func (c *Client) GetContainers(projectID string, addonID string) (*models.ListContainersResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// Containers may be part of Services or a separate endpoint
	// This may need specific SDK implementation
	return nil, errors.New("containers not yet implemented with SDK")
}

// StartExec starts an exec session
func (c *Client) StartExec(req *models.ExecRequest) (*models.ExecResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// Exec sessions may need WebSocket/terminal support
	// This may need specific SDK implementation
	return nil, errors.New("exec not yet implemented with SDK")
}

// StartShell starts a shell session
func (c *Client) StartShell(req *models.ShellRequest) (*models.ShellResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// Shell sessions may need WebSocket/terminal support
	// This may need specific SDK implementation
	return nil, errors.New("shell not yet implemented with SDK")
}

// GetAddons retrieves a list of addons
func (c *Client) GetAddons() (*models.AddonListResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.AddOns.List(ctx)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to CLI models
	addons := make([]models.Addon, len(resp.Data.AddOns))
	for i, a := range resp.Data.AddOns {
		addons[i] = models.Addon{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			Category:    a.Category,
			Icon:        a.Icon,
		}
	}

	return &models.AddonListResponse{
		Addons: addons,
	}, nil
}

// GetAddon retrieves a specific addon by ID
func (c *Client) GetAddon(addonID string) (*models.Addon, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.AddOns.Get(ctx, addonID)
	if err != nil {
		return nil, err
	}

	return &models.Addon{
		ID:          resp.Data.AddOn.ID,
		Name:        resp.Data.AddOn.Name,
		Description: resp.Data.AddOn.Description,
		Category:    resp.Data.AddOn.Category,
		Icon:        resp.Data.AddOn.Icon,
	}, nil
}

// DeployAddon deploys an addon
func (c *Client) DeployAddon(req *models.AddonDeployRequest) (*models.AddonDeployResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	// Convert map[string]string to map[string]interface{}
	config := make(map[string]interface{})
	for k, v := range req.Config {
		config[k] = v
	}

	sdkReq := &sdk.DeployAddOnRequest{
		AddOnUUID: req.AddonID,
		ProjectID: req.ProjectID,
		Config:    config,
	}

	resp, _, err := c.sdkClient.AddOns.Deploy(ctx, sdkReq)
	if err != nil {
		return nil, err
	}

	return &models.AddonDeployResponse{
		DeploymentID: resp.Data.Deployment.ID,
		Status:       resp.Data.Deployment.Status,
		Message:      resp.Message,
	}, nil
}

// GetAddonDeployments retrieves a list of addon deployments
func (c *Client) GetAddonDeployments(projectID string) ([]models.AddonDeployment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.AddOns.ListDeployments(ctx)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to CLI models and filter by projectID
	deployments := make([]models.AddonDeployment, 0)
	for _, d := range resp.Data.Deployments {
		// Filter by project ID if specified
		if projectID == "" || d.ProjectID == projectID {
			deployments = append(deployments, models.AddonDeployment{
				ID:        d.ID,
				ProjectID: d.ProjectID,
				AddonID:   d.AddOnID,
				Name:      d.AddOnName,
				Status:    d.Status,
				CreatedAt: d.CreatedAt.Time,
			})
		}
	}

	return deployments, nil
}

// DeleteAddonDeployment deletes an addon deployment
func (c *Client) DeleteAddonDeployment(deploymentID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()
	_, err := c.sdkClient.AddOns.DeleteDeployment(ctx, deploymentID)
	return err
}

// GetServers retrieves all servers
func (c *Client) GetServers() (*models.ServersResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Servers.List(ctx)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to CLI models
	servers := make([]models.Server, len(resp.Data.Servers))
	for i, s := range resp.Data.Servers {
		servers[i] = models.Server{
			ID:        s.ID,
			Name:      s.Name,
			Status:    s.Status,
			Type:      s.Provider, // Provider maps to Type in CLI
			Region:    s.Region,
			CreatedAt: s.CreatedAt.Time,
			UpdatedAt: s.UpdatedAt.Time,
		}
	}

	return &models.ServersResponse{
		Servers: servers,
		Total:   len(servers),
		Page:    1,
		PerPage: len(servers),
	}, nil
}

// GetServer retrieves a specific server by ID
func (c *Client) GetServer(serverID string) (*models.Server, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Servers.Get(ctx, serverID)
	if err != nil {
		return nil, err
	}

	return &models.Server{
		ID:        resp.Data.Server.ID,
		Name:      resp.Data.Server.Name,
		Status:    resp.Data.Server.Status,
		Type:      resp.Data.Server.Provider,
		Region:    resp.Data.Server.Region,
		CreatedAt: resp.Data.Server.CreatedAt.Time,
		UpdatedAt: resp.Data.Server.UpdatedAt.Time,
	}, nil
}

// CreateServer creates a new server
func (c *Client) CreateServer(req *models.ServerCreateRequest) (*models.Server, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	sdkReq := &sdk.CreateServerRequest{
		Name:     req.Name,
		Provider: req.Type, // Type maps to Provider in SDK
		Region:   req.Region,
	}

	resp, _, err := c.sdkClient.Servers.Create(ctx, sdkReq)
	if err != nil {
		return nil, err
	}

	return &models.Server{
		ID:        resp.Data.Server.ID,
		Name:      resp.Data.Server.Name,
		Status:    resp.Data.Server.Status,
		Type:      resp.Data.Server.Provider,
		Region:    resp.Data.Server.Region,
		CreatedAt: resp.Data.Server.CreatedAt.Time,
		UpdatedAt: resp.Data.Server.UpdatedAt.Time,
	}, nil
}

// UpdateServer updates an existing server
func (c *Client) UpdateServer(serverID string, req *models.ServerUpdateRequest) (*models.Server, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// SDK may not have Update method yet, return error for now
	return nil, errors.New("server update not yet supported in SDK")
}

// DeleteServer deletes a server
func (c *Client) DeleteServer(serverID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()
	_, err := c.sdkClient.Servers.Delete(ctx, serverID)
	return err
}
