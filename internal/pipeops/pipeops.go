package pipeops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/manifoldco/promptui"
)

// Client represents the PipeOps client wrapping the Go SDK
type Client struct {
	sdkClient *sdk.Client
	config    *config.Config
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

	return &Client{
		sdkClient: sdkClient,
		config:    cfg,
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

	// Ensure the SDK client picks up the latest token and base URL.
	if cfg.OAuth != nil && strings.TrimSpace(cfg.OAuth.AccessToken) != "" {
		c.sdkClient.SetToken(cfg.OAuth.AccessToken)
	}

	baseURL := config.GetAPIURL()
	if cfg.OAuth != nil && strings.TrimSpace(cfg.OAuth.BaseURL) != "" {
		baseURL = cfg.OAuth.BaseURL
	}
	if c.sdkClient == nil || strings.TrimSpace(baseURL) != "" {
		sdkClient, err := sdk.NewClient(baseURL,
			sdk.WithTimeout(30*time.Second),
			sdk.WithMaxRetries(3),
		)
		if err != nil {
			return fmt.Errorf("failed to initialize PipeOps SDK client: %w", err)
		}
		c.sdkClient = sdkClient
	}
	if cfg.OAuth != nil && strings.TrimSpace(cfg.OAuth.AccessToken) != "" {
		c.sdkClient.SetToken(cfg.OAuth.AccessToken)
	}
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

	if !shouldPromptForWorkspace() {
		return "", errors.New("multiple workspaces found; set PIPEOPS_WORKSPACE_UUID or run 'pipeops workspace select' before using non-interactive commands")
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

func shouldPromptForWorkspace() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return shouldPromptForWorkspaceWithMode(0, err)
	}
	return shouldPromptForWorkspaceWithMode(stat.Mode(), nil)
}

func shouldPromptForWorkspaceWithMode(stdinMode os.FileMode, statErr error) bool {
	if envEnabled("PIPEOPS_OUTPUT_JSON") || envEnabled("PIPEOPS_NON_INTERACTIVE") || isCIEnvironment() {
		return false
	}
	if statErr != nil {
		return false
	}
	return stdinMode&os.ModeCharDevice != 0
}

func envEnabled(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func isCIEnvironment() bool {
	return envEnabled("CI") || envEnabled("GITHUB_ACTIONS")
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
			id = p.ID.String()
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

	// Resolve workspace UUID
	workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
	if err != nil {
		return nil, err
	}

	// Use SDK with workspace scoping
	opts := &sdk.ProjectGetOptions{
		WorkspaceUUID: workspaceUUID,
	}
	resp, _, err := c.sdkClient.Projects.Get(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}

	id := strings.TrimSpace(resp.Data.Project.UUID)
	if id == "" {
		id = resp.Data.Project.ID.String()
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
		Name:          req.Name,
		Description:   req.Description,
		ServerID:      req.ServerID,
		EnvironmentID: req.EnvironmentID,
		Repository:    req.Repository,
		Branch:        req.Branch,
		BuildCommand:  req.BuildCommand,
		StartCommand:  req.StartCommand,
		Port:          req.Port,
		Framework:     req.Framework,
		EnvVars:       req.EnvVars,
	}

	resp, _, err := c.sdkClient.Projects.Create(ctx, createReq)
	if err != nil {
		return nil, err
	}

	id := strings.TrimSpace(resp.Data.Project.UUID)
	if id == "" {
		id = resp.Data.Project.ID.String()
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
		Name:         req.Name,
		Description:  req.Description,
		BuildCommand: req.BuildCommand,
		StartCommand: req.StartCommand,
		Port:         req.Port,
	}

	resp, _, err := c.sdkClient.Projects.Update(ctx, projectID, updateReq)
	if err != nil {
		return nil, err
	}

	id := strings.TrimSpace(resp.Data.Project.UUID)
	if id == "" {
		id = resp.Data.Project.ID.String()
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

// RestartProject restarts a project.
func (c *Client) RestartProject(projectID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()
	_, err := c.sdkClient.Projects.Restart(ctx, projectID)
	return err
}

// StopProject stops a project.
func (c *Client) StopProject(projectID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}

	ctx := context.Background()
	_, err := c.sdkClient.Projects.Stop(ctx, projectID)
	return err
}

// GetProjectEnvVariables retrieves environment variables for a project.
func (c *Client) GetProjectEnvVariables(projectID string) ([]sdk.EnvVariable, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Projects.GetEnvVariables(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return resp.Data.EnvVariables, nil
}

// UpdateProjectEnvVariables updates environment variables for a project.
func (c *Client) UpdateProjectEnvVariables(projectID string, envVars []sdk.EnvVariable) ([]sdk.EnvVariable, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Projects.UpdateEnvVariables(ctx, projectID, &sdk.EnvVariablesRequest{EnvVariables: envVars})
	if err != nil {
		return nil, err
	}
	return resp.Data.EnvVariables, nil
}

// ListProjectDeployments lists deployments for a project.
func (c *Client) ListProjectDeployments(projectID string, opts *sdk.ProjectDeploymentListOptions) (*sdk.ProjectDeploymentsResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	if opts == nil {
		opts = &sdk.ProjectDeploymentListOptions{}
	}
	if opts.WorkspaceUUID == "" && opts.WorkspaceID == "" {
		if workspaceUUID, err := c.resolveWorkspaceUUID(ctx); err == nil {
			opts.WorkspaceUUID = workspaceUUID
		}
	}
	resp, _, err := c.sdkClient.Projects.ListDeployments(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListProjectDeploymentHistory lists deployment history for a project.
func (c *Client) ListProjectDeploymentHistory(projectID string, opts *sdk.ProjectDeploymentHistoryOptions) (*sdk.ProjectDeploymentHistoryResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	if opts == nil {
		opts = &sdk.ProjectDeploymentHistoryOptions{}
	}
	if opts.WorkspaceUUID == "" && opts.WorkspaceID == "" {
		if workspaceUUID, err := c.resolveWorkspaceUUID(ctx); err == nil {
			opts.WorkspaceUUID = workspaceUUID
		}
	}
	resp, _, err := c.sdkClient.Projects.ListDeploymentHistory(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
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

	// If follow mode, poll for new logs every 2 seconds
	if req.Follow {
		seenLogs := make(map[string]bool)
		for {
			resp, _, err := c.sdkClient.Projects.TailLogs(ctx, req.ProjectID, opts)
			if err != nil {
				return fmt.Errorf("failed to fetch logs: %w", err)
			}

			// Process new logs
			for _, logMap := range resp.Data.Logs {
				msg := ""
				if m, ok := logMap["message"]; ok {
					msg = fmt.Sprintf("%v", m)
				} else {
					msg = fmt.Sprintf("%v", logMap)
				}

				// Create a unique key for this log entry
				logKey := msg
				if ts, ok := logMap["timestamp"]; ok {
					logKey = fmt.Sprintf("%v:%s", ts, msg)
				}

				// Skip already-seen logs
				if seenLogs[logKey] {
					continue
				}
				seenLogs[logKey] = true

				streamEntry := &models.StreamLogEntry{
					LogEntry: models.LogEntry{
						Message: msg,
					},
				}
				if err := callback(streamEntry); err != nil {
					return err
				}
			}

			// Wait before polling again
			time.Sleep(2 * time.Second)
		}
	}

	// Non-follow mode: just fetch logs once
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

	return nil, errors.New("container listing is not supported by the PipeOps Go SDK")
}

// StartExec starts an exec session
func (c *Client) StartExec(req *models.ExecRequest) (*models.ExecResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	return nil, errors.New("container exec is not supported by the PipeOps Go SDK")
}

// StartShell starts a shell session
func (c *Client) StartShell(req *models.ShellRequest) (*models.ShellResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	return nil, errors.New("container shell is not supported by the PipeOps Go SDK")
}

// GetAddons retrieves a list of addons
func (c *Client) GetAddons() (*models.AddonListResponse, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	// Use limit=100 to get more addons
	opts := &sdk.ListAddOnsOptions{Limit: 100}
	resp, _, err := c.sdkClient.AddOns.List(ctx, opts)
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

	// Get ID from UID, UUID, or ID field
	id := resp.Data.UID
	if id == "" {
		id = resp.Data.UUID
	}
	if id == "" {
		id = resp.Data.ID
	}

	return &models.Addon{
		ID:          id,
		Name:        resp.Data.Name,
		Description: resp.Data.Description,
		Category:    resp.Data.Category,
		Icon:        resp.Data.ImageURL,
		Status:      resp.Data.Status,
	}, nil
}

// DeployAddon deploys an addon.
func (c *Client) DeployAddon(req *sdk.DeployAddOnRequest) (*models.AddonDeployment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	if req.Workspace == "" {
		workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
		if err != nil {
			return nil, err
		}
		req.Workspace = workspaceUUID
	}

	resp, _, err := c.sdkClient.AddOns.Deploy(ctx, req)
	if err != nil {
		return nil, err
	}
	deployment := addonDeploymentFromSDK(resp.Data.Deployment)
	return &deployment, nil
}

// GetAddonDeployments retrieves a list of addon deployments for the workspace.
func (c *Client) GetAddonDeployments() ([]models.AddonDeployment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()

	// Resolve workspace UUID
	workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
	if err != nil {
		return nil, err
	}

	// Build options with workspace scoping
	opts := &sdk.ListDeploymentsOptions{
		WorkspaceUUID: workspaceUUID,
	}

	resp, _, err := c.sdkClient.AddOns.ListDeployments(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Convert SDK response to CLI models
	deployments := make([]models.AddonDeployment, 0)
	for _, d := range resp.Data {
		deployments = append(deployments, addonDeploymentFromSDK(d))
	}

	return deployments, nil
}

// GetAddonDeployment retrieves a single addon deployment.
func (c *Client) GetAddonDeployment(deploymentID string) (*models.AddonDeployment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.AddOns.GetDeployment(ctx, deploymentID)
	if err != nil {
		deployment, fallbackErr := c.findAddonDeployment(deploymentID)
		if fallbackErr == nil {
			return deployment, nil
		}
		return nil, err
	}
	deployment := addonDeploymentFromSDK(resp.Data.Deployment)
	return &deployment, nil
}

func (c *Client) findAddonDeployment(deploymentID string) (*models.AddonDeployment, error) {
	deployments, err := c.GetAddonDeployments()
	if err != nil {
		return nil, err
	}
	for _, deployment := range deployments {
		if deployment.ID == deploymentID {
			deploymentCopy := deployment
			return &deploymentCopy, nil
		}
	}
	return nil, fmt.Errorf("addon deployment %q not found", deploymentID)
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

// ListAddonCategories lists addon categories.
func (c *Client) ListAddonCategories() ([]sdk.AddOnCategory, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.AddOns.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Data.Categories, nil
}

// GetAddonDeploymentSession gets an addon deployment session.
func (c *Client) GetAddonDeploymentSession(sessionID string) (map[string]interface{}, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.AddOns.GetDeploymentSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return resp.Data.Session, nil
}

// ViewAddonDeploymentConfigs retrieves addon deployment configs.
func (c *Client) ViewAddonDeploymentConfigs(deploymentID string) (map[string]interface{}, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.AddOns.ViewDeploymentConfigs(ctx, deploymentID)
	if err != nil {
		return nil, err
	}
	return resp.Data.Configs, nil
}

func addonDeploymentFromSDK(d sdk.AddOnDeployment) models.AddonDeployment {
	id := d.UID
	if id == "" {
		id = d.DeploymentName
	}
	name := d.Name
	if name == "" {
		name = d.DeploymentName
	}
	version := d.Version
	if version == "" {
		version = d.CurrentVersion
	}
	return models.AddonDeployment{
		ID:            id,
		Name:          name,
		DeploymentURL: d.DeploymentURL,
		Category:      d.Category,
		Status:        d.Status,
		Environment:   d.Environment,
		Version:       version,
		CreatedAt:     timestampToTime(d.CreatedAt),
		UpdatedAt:     timestampToTime(d.UpdatedAt),
	}
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
		provider := s.Provider
		if provider == "" {
			provider = detectProviderFromName(s.Name)
		}
		servers[i] = models.Server{
			ID:        id,
			Name:      s.Name,
			Status:    s.Status,
			Type:      provider, // Provider maps to Type in CLI
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

	resp, _, err := c.sdkClient.Servers.Get(ctx, serverID, workspaceUUID)
	if err != nil {
		return nil, err
	}

	server := resp.Data.Server
	id := strings.TrimSpace(server.UUID)
	if id == "" {
		id = server.ID
	}
	if id == "" {
		return nil, errors.New("no cluster data returned")
	}

	provider := server.Provider
	if provider == "" {
		provider = detectProviderFromName(server.Name)
	}

	return &models.Server{
		ID:        id,
		Name:      server.Name,
		Status:    server.Status,
		Type:      provider,
		Region:    server.Region,
		CreatedAt: timestampToTime(server.CreatedAt),
		UpdatedAt: timestampToTime(server.UpdatedAt),
	}, nil
}

// GetServerConnection retrieves server connection information.
func (c *Client) GetServerConnection(serverID string) (map[string]interface{}, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Servers.GetClusterConnection(ctx, serverID)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Data.Connection) == 0 {
		return nil, errors.New("no server connection information returned")
	}
	return resp.Data.Connection, nil
}

// GetServerCostAllocation retrieves server cost allocation.
func (c *Client) GetServerCostAllocation(serverID string) (map[string]interface{}, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}

	ctx := context.Background()
	resp, _, err := c.sdkClient.Servers.GetClusterCostAllocation(ctx, serverID)
	if err != nil {
		return nil, err
	}
	return resp.Data.Costs, nil
}

// detectProviderFromName attempts to detect the cloud provider from the server name.
func detectProviderFromName(name string) string {
	name = strings.ToLower(name)
	if strings.Contains(name, "linode") {
		return "Linode"
	}
	if strings.Contains(name, "hetzner") {
		return "Hetzner"
	}
	if strings.Contains(name, "do") || strings.Contains(name, "digitalocean") {
		return "DigitalOcean"
	}
	if strings.Contains(name, "gcp") || strings.Contains(name, "google") {
		return "GCP"
	}
	if strings.Contains(name, "aws") || strings.Contains(name, "amazon") {
		return "AWS"
	}
	if strings.Contains(name, "azure") || strings.Contains(name, "microsoft") {
		return "Azure"
	}
	return "Unknown"
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

// GetWorkspace retrieves a workspace.
func (c *Client) GetWorkspace(ctx context.Context, workspaceID string) (*sdk.Workspace, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Workspaces.Get(ctx, workspaceID)
	if err != nil {
		workspace, fallbackErr := c.findWorkspace(ctx, workspaceID)
		if fallbackErr == nil {
			return workspace, nil
		}
		return nil, err
	}
	return &resp.Data.Workspace, nil
}

func (c *Client) findWorkspace(ctx context.Context, workspaceID string) (*sdk.Workspace, error) {
	workspaces, err := c.GetWorkspaces(ctx)
	if err != nil {
		return nil, err
	}
	for _, workspace := range workspaces {
		if workspace.UUID == workspaceID || workspace.ID == workspaceID || workspace.Name == workspaceID {
			workspaceCopy := workspace
			return &workspaceCopy, nil
		}
	}
	return nil, fmt.Errorf("workspace %q not found", workspaceID)
}

// CreateWorkspace creates a workspace.
func (c *Client) CreateWorkspace(ctx context.Context, req *sdk.CreateWorkspaceRequest) (*sdk.Workspace, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Workspaces.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Workspace, nil
}

// UpdateWorkspace updates a workspace.
func (c *Client) UpdateWorkspace(ctx context.Context, workspaceID string, req *sdk.UpdateWorkspaceRequest) (*sdk.Workspace, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Workspaces.Update(ctx, workspaceID, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Workspace, nil
}

// DeleteWorkspace deletes a workspace.
func (c *Client) DeleteWorkspace(ctx context.Context, workspaceID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := c.sdkClient.Workspaces.Delete(ctx, workspaceID)
	return err
}

// ListEnvironments lists environments.
func (c *Client) ListEnvironments(ctx context.Context) ([]sdk.Environment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Environments.List(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Data.Environments, nil
}

// GetEnvironment retrieves an environment.
func (c *Client) GetEnvironment(ctx context.Context, environmentID string) (*sdk.Environment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Environments.Get(ctx, environmentID)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Environment, nil
}

// CreateEnvironment creates an environment.
func (c *Client) CreateEnvironment(ctx context.Context, req *sdk.CreateEnvironmentRequest) (*sdk.Environment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if req.WorkspaceUUID == "" && req.WorkspaceID == "" {
		workspaceUUID, err := c.resolveWorkspaceUUID(ctx)
		if err != nil {
			return nil, err
		}
		req.WorkspaceUUID = workspaceUUID
	}

	resp, _, err := c.sdkClient.Environments.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Environment, nil
}

// UpdateEnvironment updates an environment.
func (c *Client) UpdateEnvironment(ctx context.Context, environmentID string, req *sdk.UpdateEnvironmentRequest) (*sdk.Environment, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.Environments.Update(ctx, environmentID, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Environment, nil
}

// DeleteEnvironment deletes an environment.
func (c *Client) DeleteEnvironment(ctx context.Context, environmentID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := c.sdkClient.Environments.Delete(ctx, environmentID)
	return err
}

// SetEnvironmentVariables sets environment variables for an environment.
func (c *Client) SetEnvironmentVariables(ctx context.Context, environmentID string, envVars []sdk.EnvVariable) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := c.sdkClient.Environments.SetEnvVariables(ctx, environmentID, &sdk.SetEnvironmentVariablesRequest{EnvVariables: envVars})
	return err
}

// ListServiceAccountTokens lists service account tokens.
func (c *Client) ListServiceAccountTokens(ctx context.Context) ([]sdk.ServiceAccountToken, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.ServiceTokens.ListServiceAccountTokens(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Data.Tokens, nil
}

// GetServiceAccountToken retrieves a service account token.
func (c *Client) GetServiceAccountToken(ctx context.Context, tokenID string) (*sdk.ServiceAccountToken, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.ServiceTokens.GetServiceAccountToken(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Token, nil
}

// CreateServiceAccountToken creates a service account token.
func (c *Client) CreateServiceAccountToken(ctx context.Context, req *sdk.ServiceAccountTokenRequest) (*sdk.ServiceAccountToken, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.ServiceTokens.CreateServiceAccountToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Token, nil
}

// UpdateServiceAccountToken updates a service account token.
func (c *Client) UpdateServiceAccountToken(ctx context.Context, tokenID string, req *sdk.ServiceAccountTokenUpdateRequest) (*sdk.ServiceAccountToken, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, _, err := c.sdkClient.ServiceTokens.UpdateServiceAccountToken(ctx, tokenID, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Token, nil
}

// RevokeServiceAccountToken revokes a service account token.
func (c *Client) RevokeServiceAccountToken(ctx context.Context, tokenID string) error {
	if !c.IsAuthenticated() {
		return errors.New("not authenticated")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := c.sdkClient.ServiceTokens.RevokeServiceAccountToken(ctx, tokenID)
	return err
}
