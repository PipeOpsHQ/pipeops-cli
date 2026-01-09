package pipeops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/libs"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/manifoldco/promptui"
)

// Client represents the PipeOps client wrapping the Go SDK
type Client struct {
	sdkClient    *sdk.Client
	config       *config.Config
	legacyClient libs.HttpClients // Required for exec/shell/containers not in SDK
}

// NewClient creates a new PipeOps client
func NewClient() ClientAPI {
	cfg := config.DefaultConfig()
	baseURL := config.GetAPIURL()

	sdkClient, err := sdk.NewClient(baseURL)
	if err != nil {
		// Fallback to default if URL parsing fails
		sdkClient, _ = sdk.NewClient("")
	}

	legacyClient := libs.NewHttpClientWithURL(baseURL)

	return &Client{
		sdkClient:    sdkClient,
		config:       cfg,
		legacyClient: legacyClient,
	}
}

// NewClientWithConfig creates a new PipeOps client with the provided configuration
func NewClientWithConfig(cfg *config.Config) ClientAPI {
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

	legacyClient := libs.NewHttpClientWithURL(baseURL)

	return &Client{
		sdkClient:    sdkClient,
		config:       cfg,
		legacyClient: legacyClient,
	}
}

// LoadConfig loads the configuration from the config file
func (c *Client) LoadConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c.config = cfg

	// Ensure the SDK client picks up the latest token and base URL.
	if cfg.OAuth != nil && strings.TrimSpace(cfg.OAuth.AccessToken) != "" {
		c.sdkClient.SetToken(cfg.OAuth.AccessToken)
	}

	baseURL := config.GetAPIURL()
	if cfg.OAuth != nil && strings.TrimSpace(cfg.OAuth.BaseURL) != "" {
		baseURL = cfg.OAuth.BaseURL
	}
	c.legacyClient = libs.NewHttpClientWithURL(baseURL)
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

func (c *Client) getClusterUUID() string {
	// Prefer env var so it can be overridden per-invocation without editing config.
	if v := strings.TrimSpace(os.Getenv("PIPEOPS_CLUSTER_UUID")); v != "" {
		return v
	}
	if c.config != nil && c.config.Settings != nil {
		return strings.TrimSpace(c.config.Settings.DefaultClusterUUID)
	}
	return ""
}

func (c *Client) getWorkspaceUUID() string {
	// Prefer env var so it can be overridden per-invocation without editing config.
	if v := strings.TrimSpace(os.Getenv("PIPEOPS_WORKSPACE_UUID")); v != "" {
		return v
	}
	if c.config != nil && c.config.Settings != nil {
		return strings.TrimSpace(c.config.Settings.DefaultWorkspaceUUID)
	}
	return ""
}

func (c *Client) resolveWorkspaceUUID(ctx context.Context) (string, error) {
	if workspaceUUID := c.getWorkspaceUUID(); workspaceUUID != "" {
		return workspaceUUID, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Workspaces.List(ctx)
	if err != nil {
		return "", err
	}

	workspaces := resp.Data.Workspaces
	if len(workspaces) == 0 {
		return "", errors.New("no workspaces found - create a workspace first")
	}

	// If only one workspace, auto-select it
	if len(workspaces) == 1 {
		uuid := strings.TrimSpace(workspaces[0].UUID)
		if uuid != "" {
			// Save selection for future use
			c.saveDefaultWorkspace(uuid)
			return uuid, nil
		}
	}

	// Multiple workspaces exist - prompt user to select
	fmt.Println("\nYou have multiple workspaces. Please select one:")
	options := make([]string, len(workspaces))
	for i, ws := range workspaces {
		options[i] = fmt.Sprintf("%s (%s)", ws.Name, ws.UUID)
	}

	idx, err := promptSelectWorkspace(options)
	if err != nil {
		return "", fmt.Errorf("workspace selection cancelled: %w", err)
	}

	selectedUUID := strings.TrimSpace(workspaces[idx].UUID)
	if selectedUUID == "" {
		return "", errors.New("selected workspace has no UUID")
	}

	// Save selection for future use
	c.saveDefaultWorkspace(selectedUUID)
	fmt.Printf("✓ Selected workspace: %s\n\n", workspaces[idx].Name)

	return selectedUUID, nil
}

// promptSelectWorkspace uses promptui to let user select a workspace
func promptSelectWorkspace(options []string) (int, error) {
	prompt := promptui.Select{
		Label: "Select a workspace",
		Items: options,
		Size:  10,
	}

	idx, _, err := prompt.Run()
	return idx, err
}

// saveDefaultWorkspace saves the selected workspace to config
func (c *Client) saveDefaultWorkspace(uuid string) {
	if c.config != nil && c.config.Settings != nil {
		c.config.Settings.DefaultWorkspaceUUID = uuid
		_ = config.Save(c.config) // Best effort save
	}
}

func sdkStatusCode(err error) (int, bool) {
	var apiErr *sdk.ErrorResponse
	if errors.As(err, &apiErr) && apiErr != nil && apiErr.Response != nil {
		return apiErr.Response.StatusCode, true
	}
	return 0, false
}

func timestampToTime(ts *sdk.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.Time
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

	// Resolve workspace UUID to scope project listing
	workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
	if err != nil {
		return nil, err
	}

	// Use SDK with workspace scoping
	opts := &sdk.ProjectListOptions{
		WorkspaceUUID: workspaceUUID,
	}
	resp, _, err := c.sdkClient.Projects.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to models
	var projects []models.Project
	for _, p := range resp.Data.Projects {
		id := strings.TrimSpace(p.UUID)
		if id == "" {
			id = p.ID
		}
		projects = append(projects, models.Project{
			ID:          id,
			Name:        p.Name,
			Description: p.Description,
			Status:      p.Status,
		})
	}

	return &models.ProjectsResponse{Projects: projects}, nil
}

// GetProject retrieves a specific project
func (c *Client) GetProject(projectID string) (*models.Project, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()

	// SDK v0.6.2+ handles workspace fallback automatically
	resp, _, err := c.sdkClient.Projects.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}

	id := strings.TrimSpace(resp.Data.Project.UUID)
	if id == "" {
		id = resp.Data.Project.ID
	}

	return &models.Project{
		ID:          id,
		Name:        resp.Data.Project.Name,
		Description: resp.Data.Project.Description,
		Status:      resp.Data.Project.Status,
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

	id := strings.TrimSpace(resp.Data.Project.UUID)
	if id == "" {
		id = resp.Data.Project.ID
	}

	return &models.Project{
		ID:          id,
		Name:        resp.Data.Project.Name,
		Description: resp.Data.Project.Description,
		Status:      resp.Data.Project.Status,
		CreatedAt:   timestampToTime(resp.Data.Project.CreatedAt),
		UpdatedAt:   timestampToTime(resp.Data.Project.UpdatedAt),
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

	id := strings.TrimSpace(resp.Data.Project.UUID)
	if id == "" {
		id = resp.Data.Project.ID
	}

	return &models.Project{
		ID:          id,
		Name:        resp.Data.Project.Name,
		Description: resp.Data.Project.Description,
		Status:      resp.Data.Project.Status,
		CreatedAt:   timestampToTime(resp.Data.Project.CreatedAt),
		UpdatedAt:   timestampToTime(resp.Data.Project.UpdatedAt),
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

// DeployProject triggers a deployment for a project
func (c *Client) DeployProject(projectID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()
	_, err := c.sdkClient.Projects.Deploy(ctx, projectID)
	return err
}

// GetLogs retrieves project logs
func (c *Client) GetLogs(req *models.LogsRequest) (*models.LogsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()

	// Resolve workspace UUID
	workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
	if err != nil {
		return nil, err
	}

	// Use SDK with workspace scoping
	opts := &sdk.LogsOptions{
		Limit:         req.Limit,
		WorkspaceUUID: workspaceUUID,
		App:           "project",
	}
	if req.Since != nil {
		opts.StartTime = req.Since.Format(time.RFC3339)
	}
	if req.Until != nil {
		opts.EndTime = req.Until.Format(time.RFC3339)
	}

	resp, _, err := c.sdkClient.Projects.GetLogs(ctx, req.ProjectID, opts)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to models
	var entries []models.LogEntry
	for _, logMap := range resp.Data.Logs {
		entry := models.LogEntry{}

		// Extract message
		if m, ok := logMap["message"]; ok {
			entry.Message = fmt.Sprintf("%v", m)
		} else if m, ok := logMap["log"]; ok {
			entry.Message = fmt.Sprintf("%v", m)
		} else {
			// Fallback: stringify the whole map
			entry.Message = fmt.Sprintf("%v", logMap)
		}

		// Extract timestamp
		if ts, ok := logMap["timestamp"]; ok {
			if tsStr, ok := ts.(string); ok {
				if t, err := time.Parse(time.RFC3339Nano, tsStr); err == nil {
					entry.Timestamp = t
				} else if t, err := time.Parse(time.RFC3339, tsStr); err == nil {
					entry.Timestamp = t
				}
			}
		} else if ts, ok := logMap["time"]; ok {
			if tsStr, ok := ts.(string); ok {
				if t, err := time.Parse(time.RFC3339Nano, tsStr); err == nil {
					entry.Timestamp = t
				} else if t, err := time.Parse(time.RFC3339, tsStr); err == nil {
					entry.Timestamp = t
				}
			}
		}

		// Extract level
		if lvl, ok := logMap["level"]; ok {
			entry.Level = models.LogLevel(fmt.Sprintf("%v", lvl))
		} else if lvl, ok := logMap["severity"]; ok {
			entry.Level = models.LogLevel(fmt.Sprintf("%v", lvl))
		} else {
			entry.Level = models.LogLevelInfo
		}

		entries = append(entries, entry)
	}

	return &models.LogsResponse{Logs: entries}, nil
}

// StreamLogs streams project logs
func (c *Client) StreamLogs(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()

	// Resolve workspace UUID
	workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
	if err != nil {
		return err
	}

	// Build SDK request options with workspace scoping
	opts := &sdk.LogsOptions{
		Limit:         req.Limit,
		WorkspaceUUID: workspaceUUID,
		App:           "project",
	}

	// For now, just fetch logs (SDK may not have streaming support yet)
	resp, _, err := c.sdkClient.Projects.TailLogs(ctx, req.ProjectID, opts)
	if err != nil {
		return err
	}

	// Convert and callback with each log entry
	for _, logMap := range resp.Data.Logs {
		msg := ""
		if m, ok := logMap["message"]; ok {
			msg = fmt.Sprintf("%v", m)
		} else {
			msg = fmt.Sprintf("%v", logMap)
		}
		streamEntry := &models.StreamLogEntry{
			LogEntry: models.LogEntry{
				Message: msg,
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

	// Use legacy client as fallback since SDK doesn't support containers yet
	token := c.config.OAuth.AccessToken
	return c.legacyClient.GetContainers(token, projectID, addonID)
}

// StartExec starts an exec session
func (c *Client) StartExec(req *models.ExecRequest) (*models.ExecResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// Use legacy client as fallback since SDK doesn't support exec yet
	token := c.config.OAuth.AccessToken
	return c.legacyClient.StartExec(token, req)
}

// StartShell starts a shell session
func (c *Client) StartShell(req *models.ShellRequest) (*models.ShellResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	// Use legacy client as fallback since SDK doesn't support shell yet
	token := c.config.OAuth.AccessToken
	return c.legacyClient.StartShell(token, req)
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
	addons := make([]models.Addon, len(resp.Data))
	for i, a := range resp.Data {
		id := a.UID
		if id == "" {
			id = a.UUID
		}
		if id == "" {
			id = a.ID
		}
		addons[i] = models.Addon{
			ID:          id,
			Name:        a.Name,
			Description: a.Description,
			Category:    a.Category,
			Icon:        a.ImageURL,
			Status:      a.Status,
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
				CreatedAt: timestampToTime(d.CreatedAt),
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
	workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.sdkClient.Servers.List(ctx, workspaceUUID)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to CLI models
	servers := make([]models.Server, len(resp.Data.Servers))
	for i, s := range resp.Data.Servers {
		id := strings.TrimSpace(s.UUID)
		if id == "" {
			id = s.ID
		}
		servers[i] = models.Server{
			ID:        id,
			Name:      s.Name,
			Status:    s.Status,
			Type:      s.Provider, // Provider maps to Type in CLI
			Region:    s.Region,
			CreatedAt: timestampToTime(s.CreatedAt),
			UpdatedAt: timestampToTime(s.UpdatedAt),
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
	workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
	if err != nil {
		return nil, err
	}

	// In the SDK, "server" resources are represented as clusters and require a workspace UUID.
	resp, _, err := c.sdkClient.Servers.Get(ctx, serverID, workspaceUUID)
	if err != nil {
		return nil, err
	}

	id := strings.TrimSpace(resp.Data.Server.UUID)
	if id == "" {
		id = resp.Data.Server.ID
	}

	return &models.Server{
		ID:        id,
		Name:      resp.Data.Server.Name,
		Status:    resp.Data.Server.Status,
		Type:      resp.Data.Server.Provider,
		Region:    resp.Data.Server.Region,
		CreatedAt: timestampToTime(resp.Data.Server.CreatedAt),
		UpdatedAt: timestampToTime(resp.Data.Server.UpdatedAt),
	}, nil
}

// CreateServer creates a new server
func (c *Client) CreateServer(req *models.ServerCreateRequest) (*models.Server, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	return nil, errors.New("server creation is not yet supported via the CLI; use the PipeOps web console")
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
	_, err := c.sdkClient.Servers.Delete(ctx, serverID, "")
	return err
}

// GetWorkspaces retrieves all workspaces for the authenticated user
func (c *Client) GetWorkspaces(ctx context.Context) ([]sdk.Workspace, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Workspaces.List(ctx)
	if err != nil {
		return nil, err
	}

	return resp.Data.Workspaces, nil
}
