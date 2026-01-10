package libs

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/auth"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/go-resty/resty/v2"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrVerificationFailed = errors.New("token verification failed")
)

// Removed init() function - now using secure configuration approach

// validateToken proactively validates a token by making a test request
func (h *HttpClient) validateToken(token string) error {
	if strings.TrimSpace(token) == "" {
		return auth.NewAuthError("token_invalid", "Token is empty", 401, nil)
	}

	// Try to validate token by making a request to a validation endpoint
	// We'll use the userinfo endpoint since it's lightweight and always available
	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/oauth/userinfo")
	if err != nil {
		return auth.NewAuthError("validation_failed",
			fmt.Sprintf("Token validation failed: %v", err), 401, err)
	}

	// Check for authentication errors
	if resp.StatusCode() == 401 {
		return detectAuthError(resp)
	}

	// If we get a successful response, token is valid
	if resp.StatusCode() == 200 {
		return nil
	}

	// For other status codes, assume token is invalid
	return auth.NewAuthError("token_invalid",
		fmt.Sprintf("Token validation failed with status %d", resp.StatusCode()),
		resp.StatusCode(), nil)
}

// detectAuthError analyzes API response to determine the specific authentication error
func detectAuthError(resp *resty.Response) error {
	if resp.StatusCode() != 401 {
		return nil
	}

	// Try to parse error response for more specific error information
	var errorResp struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
		Code        int    `json:"code"`
	}

	body := resp.Body()
	if len(body) > 0 {
		if err := json.Unmarshal(body, &errorResp); err == nil {
			// Check for specific error types in the response
			errorStr := strings.ToLower(errorResp.Error)
			descStr := strings.ToLower(errorResp.Description)

			// Token expired
			if strings.Contains(errorStr, "expired") ||
				strings.Contains(descStr, "expired") ||
				strings.Contains(errorStr, "expiration") ||
				strings.Contains(descStr, "expiration") {
				return auth.NewAuthError("token_expired",
					"Your session has expired. Please run 'pipeops auth login' to authenticate again.",
					401, nil)
			}

			// Token revoked
			if strings.Contains(errorStr, "revoked") ||
				strings.Contains(descStr, "revoked") ||
				strings.Contains(errorStr, "invalidated") ||
				strings.Contains(descStr, "invalidated") {
				return auth.NewAuthError("token_revoked",
					"Your session has been revoked. Please run 'pipeops auth login' to authenticate again.",
					401, nil)
			}

			// Invalid token
			if strings.Contains(errorStr, "invalid") ||
				strings.Contains(descStr, "invalid") ||
				strings.Contains(errorStr, "malformed") ||
				strings.Contains(descStr, "malformed") {
				return auth.NewAuthError("token_invalid",
					"Your authentication token is invalid. Please run 'pipeops auth login' to authenticate again.",
					401, nil)
			}

			// Provide specific error message if available
			if errorResp.Description != "" {
				return auth.NewAuthError("authentication_failed",
					fmt.Sprintf("Authentication failed: %s", errorResp.Description),
					401, nil)
			}
		}
	}

	// Default 401 error
	return auth.NewAuthError("authentication_failed",
		"Authentication failed. Please run 'pipeops auth login' to authenticate again.",
		401, nil)
}

type HttpClients interface {
	VerifyToken(token string, operatorID string) (*models.PipeOpsTokenVerificationResponse, error)
	GetProjects(token string) (*models.ProjectsResponse, error)
	GetProjectsByWorkspace(token string, workspaceUUID string) (*models.ProjectsResponse, error)
	GetProject(token string, projectID string) (*models.Project, error)
	GetProjectByWorkspace(token string, projectID string, workspaceUUID string) (*models.Project, error)
	CreateProject(token string, req *models.ProjectCreateRequest) (*models.Project, error)
	UpdateProject(token string, projectID string, req *models.ProjectUpdateRequest) (*models.Project, error)
	DeleteProject(token string, projectID string) error
	GetLogs(token string, req *models.LogsRequest) (*models.LogsResponse, error)
	GetLogsByWorkspace(token string, req *models.LogsRequest, workspaceUUID string) (*models.LogsResponse, error)
	StreamLogs(token string, req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error
	GetServices(token string, projectID string, addonID string) (*models.ListServicesResponse, error)
	StartProxy(token string, req *models.ProxyRequest) (*models.ProxyResponse, error)
	GetContainers(token string, projectID string, addonID string) (*models.ListContainersResponse, error)
	StartExec(token string, req *models.ExecRequest) (*models.ExecResponse, error)
	StartShell(token string, req *models.ShellRequest) (*models.ShellResponse, error)

	// Addon Management
	GetAddons(token string) (*models.AddonListResponse, error)
	GetAddon(token string, addonID string) (*models.Addon, error)
	GetAddonDeployments(token string, projectID string) ([]models.AddonDeployment, error)
	DeleteAddonDeployment(token string, deploymentID string) error

	// Server Management
	GetServers(token string) (*models.ServersResponse, error)
	GetServer(token string, serverID string) (*models.Server, error)
	CreateServer(token string, req *models.ServerCreateRequest) (*models.Server, error)
	UpdateServer(token string, serverID string, req *models.ServerUpdateRequest) (*models.Server, error)
	DeleteServer(token string, serverID string) error
}

type HttpClient struct {
	client  *resty.Client
	baseURL string
}

func NewHttpClient() HttpClients {
	return NewHttpClientWithURL(config.GetAPIURL())
}

func NewHttpClientWithURL(baseURL string) HttpClients {
	r := resty.New()

	// Enable debug mode if environment variable is set
	if os.Getenv("PIPEOPS_DEBUG") == "true" {
		r.Debug = true
	}

	// Use the provided base URL
	URL := strings.TrimSpace(baseURL)
	r.SetBaseURL(URL)

	return &HttpClient{
		client:  r,
		baseURL: URL,
	}
}

// VerifyToken performs a POST request to verify a token.
func (v *HttpClient) VerifyToken(token string, operatorID string) (*models.PipeOpsTokenVerificationResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	payload := map[string]string{
		"token": token,
	}

	// Add operator_id only if provided
	if strings.TrimSpace(operatorID) != "" {
		payload["operator_id"] = operatorID
	}

	resp, err := v.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 401 || resp.StatusCode() == 400 {
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, ErrVerificationFailed
	}

	body := resp.Body()
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body from server (status: %d)", resp.StatusCode())
	}

	var respData *models.PipeOpsTokenVerificationResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !respData.Valid {
		return nil, ErrInvalidToken
	}

	return respData, nil
}

// GetProjects retrieves all projects for the authenticated user
func (h *HttpClient) GetProjects(token string) (*models.ProjectsResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	// Validate token before making the request
	if err := h.validateToken(token); err != nil {
		return nil, err
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/projects")
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	if resp.StatusCode() == 401 {
		if authErr := detectAuthError(resp); authErr != nil {
			return nil, authErr
		}
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	// Handle empty response body
	body := resp.Body()
	if len(body) == 0 {
		// Return empty projects response when API returns empty body
		return &models.ProjectsResponse{
			Projects: []models.Project{},
			Total:    0,
			Page:     1,
			PerPage:  10,
		}, nil
	}

	var projectsResp *models.ProjectsResponse
	if err := json.Unmarshal(body, &projectsResp); err != nil {
		return nil, fmt.Errorf("failed to parse projects response: %w", err)
	}

	return projectsResp, nil
}

// GetProjectsByWorkspace retrieves projects for a specific workspace
func (h *HttpClient) GetProjectsByWorkspace(token string, workspaceUUID string) (*models.ProjectsResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(workspaceUUID) == "" {
		return nil, errors.New("workspace UUID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/workspace/fetch/" + workspaceUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	if resp.StatusCode() == 401 {
		if authErr := detectAuthError(resp); authErr != nil {
			return nil, authErr
		}
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("workspace not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	// Parse workspace response to extract projects
	var wsResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Workspace struct {
				Projects []struct {
					ID     int    `json:"ID"`
					UUID   string `json:"UUID"`
					Name   string `json:"Name"`
					Status string `json:"Status"`
				} `json:"Projects"`
			} `json:"workspace"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &wsResp); err != nil {
		return nil, fmt.Errorf("failed to parse workspace response: %w", err)
	}

	projects := make([]models.Project, len(wsResp.Data.Workspace.Projects))
	for i, p := range wsResp.Data.Workspace.Projects {
		id := p.UUID
		if id == "" {
			id = fmt.Sprintf("%d", p.ID)
		}
		projects[i] = models.Project{
			ID:     id,
			Name:   p.Name,
			Status: p.Status,
		}
	}

	return &models.ProjectsResponse{
		Projects: projects,
		Total:    len(projects),
		Page:     1,
		PerPage:  len(projects),
	}, nil
}

// GetProjectByWorkspace retrieves a specific project by ID with workspace scope
func (h *HttpClient) GetProjectByWorkspace(token string, projectID string, workspaceUUID string) (*models.Project, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(projectID) == "" {
		return nil, errors.New("project ID is empty")
	}

	if strings.TrimSpace(workspaceUUID) == "" {
		return nil, errors.New("workspace UUID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetQueryParam("workspace_uuid", workspaceUUID).
		Get("/project/fetch/" + projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	if resp.StatusCode() == 401 {
		if authErr := detectAuthError(resp); authErr != nil {
			return nil, authErr
		}
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	// Parse response
	var apiResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Project struct {
				ID          int    `json:"ID"`
				UUID        string `json:"UUID"`
				Name        string `json:"Name"`
				Status      string `json:"Status"`
				Repository  string `json:"Repository"`
				Branch      string `json:"Branch"`
				Language    string `json:"Language"`
				ClusterName string `json:"ClusterName"`
				ClusterUUID string `json:"ClusterUUID"`
			} `json:"project"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse project response: %w", err)
	}

	id := apiResp.Data.Project.UUID
	if id == "" {
		id = fmt.Sprintf("%d", apiResp.Data.Project.ID)
	}

	return &models.Project{
		ID:     id,
		Name:   apiResp.Data.Project.Name,
		Status: apiResp.Data.Project.Status,
	}, nil
}

// GetProject retrieves a specific project by ID
func (h *HttpClient) GetProject(token string, projectID string) (*models.Project, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(projectID) == "" {
		return nil, errors.New("project ID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/projects/" + projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var project *models.Project
	if err := json.Unmarshal(resp.Body(), &project); err != nil {
		return nil, fmt.Errorf("failed to parse project response: %w", err)
	}

	return project, nil
}

// CreateProject creates a new project
func (h *HttpClient) CreateProject(token string, req *models.ProjectCreateRequest) (*models.Project, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post("/projects")
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var project *models.Project
	if err := json.Unmarshal(resp.Body(), &project); err != nil {
		return nil, fmt.Errorf("failed to parse project response: %w", err)
	}

	return project, nil
}

// UpdateProject updates an existing project
func (h *HttpClient) UpdateProject(token string, projectID string, req *models.ProjectUpdateRequest) (*models.Project, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(projectID) == "" {
		return nil, errors.New("project ID is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Put("/projects/" + projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var project *models.Project
	if err := json.Unmarshal(resp.Body(), &project); err != nil {
		return nil, fmt.Errorf("failed to parse project response: %w", err)
	}

	return project, nil
}

// DeleteProject deletes a project
func (h *HttpClient) DeleteProject(token string, projectID string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token is empty")
	}

	if strings.TrimSpace(projectID) == "" {
		return errors.New("project ID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Delete("/projects/" + projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if resp.StatusCode() == 401 {
		return ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("project not found")
	}

	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.String())
	}

	return nil
}

// GetLogs retrieves logs for a project or addon
func (h *HttpClient) GetLogs(token string, req *models.LogsRequest) (*models.LogsResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	if req.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}

	// Build query parameters
	request := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json")

	// Add query parameters
	if req.Level != "" {
		request.SetQueryParam("level", string(req.Level))
	}
	if req.Source != "" {
		request.SetQueryParam("source", req.Source)
	}
	if req.Container != "" {
		request.SetQueryParam("container", req.Container)
	}
	if req.Since != nil {
		request.SetQueryParam("since", req.Since.Format(time.RFC3339))
	}
	if req.Until != nil {
		request.SetQueryParam("until", req.Until.Format(time.RFC3339))
	}
	if req.Limit > 0 {
		request.SetQueryParam("limit", fmt.Sprintf("%d", req.Limit))
	}
	if req.Cursor != "" {
		request.SetQueryParam("cursor", req.Cursor)
	}
	if req.Tail > 0 {
		request.SetQueryParam("tail", fmt.Sprintf("%d", req.Tail))
	}
	if req.AddonID != "" {
		request.SetQueryParam("addon_id", req.AddonID)
	}

	// Determine endpoint based on whether it's for project or addon
	endpoint := fmt.Sprintf("/projects/%s/logs", req.ProjectID)
	if req.AddonID != "" {
		endpoint = fmt.Sprintf("/projects/%s/addons/%s/logs", req.ProjectID, req.AddonID)
	}

	resp, err := request.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project or addon not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var logsResp *models.LogsResponse
	if err := json.Unmarshal(resp.Body(), &logsResp); err != nil {
		return nil, fmt.Errorf("failed to parse logs response: %w", err)
	}

	return logsResp, nil
}

// GetLogsByWorkspace retrieves logs for a project using workspace-scoped endpoint
func (h *HttpClient) GetLogsByWorkspace(token string, req *models.LogsRequest, workspaceUUID string) (*models.LogsResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	if req.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}

	if workspaceUUID == "" {
		return nil, errors.New("workspace UUID is required")
	}

	// Build query parameters - API requires app=project and workspace_uuid
	request := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetQueryParam("app", "project").
		SetQueryParam("workspace_uuid", workspaceUUID)

	// Add optional query parameters
	if req.Since != nil {
		request.SetQueryParam("start_time", req.Since.Format(time.RFC3339))
	}
	if req.Until != nil {
		request.SetQueryParam("end_time", req.Until.Format(time.RFC3339))
	}
	if req.Limit > 0 {
		request.SetQueryParam("limit", fmt.Sprintf("%d", req.Limit))
	}
	if req.Tail > 0 {
		request.SetQueryParam("limit", fmt.Sprintf("%d", req.Tail))
	}

	endpoint := fmt.Sprintf("/project/logs/%s", req.ProjectID)

	resp, err := request.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	if resp.StatusCode() == 401 {
		if authErr := detectAuthError(resp); authErr != nil {
			return nil, authErr
		}
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	// Parse response
	var apiResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    []struct {
			Timestamp string `json:"timestamp"`
			Message   string `json:"message"`
			Level     string `json:"level"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse logs response: %w", err)
	}

	logs := make([]models.LogEntry, 0)
	if apiResp.Data != nil {
		for _, l := range apiResp.Data {
			entry := models.LogEntry{
				Message:   l.Message,
				Level:     models.LogLevel(l.Level),
				Timestamp: parseTimestamp(l.Timestamp),
			}
			logs = append(logs, entry)
		}
	}

	return &models.LogsResponse{
		Logs:       logs,
		TotalCount: len(logs),
		HasMore:    false,
	}, nil
}

// parseTimestamp parses a timestamp string to time.Time
func parseTimestamp(ts string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		t, _ = time.Parse(time.RFC3339, ts)
	}
	return t
}

// StreamLogs streams logs in real-time for a project or addon
func (h *HttpClient) StreamLogs(token string, req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token is empty")
	}

	if req == nil {
		return errors.New("request is nil")
	}

	if req.ProjectID == "" {
		return errors.New("project ID is required")
	}

	if callback == nil {
		return errors.New("callback is required")
	}

	// Force follow mode for streaming
	req.Follow = true

	// Build the streaming endpoint
	endpoint := fmt.Sprintf("/projects/%s/logs/stream", req.ProjectID)
	if req.AddonID != "" {
		endpoint = fmt.Sprintf("/projects/%s/addons/%s/logs/stream", req.ProjectID, req.AddonID)
	}

	// Create request with body for streaming
	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "text/event-stream").
		SetBody(req).
		SetDoNotParseResponse(true).
		Post(endpoint)

	if err != nil {
		return fmt.Errorf("failed to start log stream: %w", err)
	}

	if resp.StatusCode() == 401 {
		return ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("project or addon not found")
	}

	if resp.IsError() {
		return fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	defer resp.RawBody().Close()

	// Parse Server-Sent Events stream
	scanner := bufio.NewScanner(resp.RawBody())
	var currentData string

	for scanner.Scan() {
		line := scanner.Text()

		// Handle Server-Sent Events format
		if strings.HasPrefix(line, "data: ") {
			currentData = strings.TrimPrefix(line, "data: ")
		} else if line == "" && currentData != "" {
			// Empty line indicates end of event
			var streamResp models.LogsStreamResponse
			if err := json.Unmarshal([]byte(currentData), &streamResp); err != nil {
				continue // Skip invalid JSON
			}

			if streamResp.Error != "" {
				return fmt.Errorf("stream error: %s", streamResp.Error)
			}

			if streamResp.Entry != nil {
				if err := callback(streamResp.Entry); err != nil {
					return err
				}

				// Check for EOF
				if streamResp.Entry.EOF {
					break
				}
			}

			currentData = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log stream: %w", err)
	}

	return nil
}

// GetServices retrieves available services for a project or addon
func (h *HttpClient) GetServices(token string, projectID string, addonID string) (*models.ListServicesResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(projectID) == "" {
		return nil, errors.New("project ID is required")
	}

	// Determine endpoint based on whether it's for project or addon
	endpoint := fmt.Sprintf("/projects/%s/services", projectID)
	if addonID != "" {
		endpoint = fmt.Sprintf("/projects/%s/addons/%s/services", projectID, addonID)
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project or addon not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var servicesResp *models.ListServicesResponse
	if err := json.Unmarshal(resp.Body(), &servicesResp); err != nil {
		return nil, fmt.Errorf("failed to parse services response: %w", err)
	}

	return servicesResp, nil
}

// StartProxy starts a proxy session for a project or addon service
func (h *HttpClient) StartProxy(token string, req *models.ProxyRequest) (*models.ProxyResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	if req.Target.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}

	if req.Target.ServiceName == "" {
		return nil, errors.New("service name is required")
	}

	// Determine endpoint based on whether it's for project or addon
	endpoint := fmt.Sprintf("/projects/%s/proxy", req.Target.ProjectID)
	if req.Target.AddonID != "" {
		endpoint = fmt.Sprintf("/projects/%s/addons/%s/proxy", req.Target.ProjectID, req.Target.AddonID)
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to start proxy: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project, addon, or service not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var proxyResp *models.ProxyResponse
	if err := json.Unmarshal(resp.Body(), &proxyResp); err != nil {
		return nil, fmt.Errorf("failed to parse proxy response: %w", err)
	}

	return proxyResp, nil
}

// GetContainers retrieves available containers for a project or addon
func (h *HttpClient) GetContainers(token string, projectID string, addonID string) (*models.ListContainersResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(projectID) == "" {
		return nil, errors.New("project ID is required")
	}

	// Determine endpoint based on whether it's for project or addon
	endpoint := fmt.Sprintf("/projects/%s/containers", projectID)
	if addonID != "" {
		endpoint = fmt.Sprintf("/projects/%s/addons/%s/containers", projectID, addonID)
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get containers: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project or addon not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var containersResp *models.ListContainersResponse
	if err := json.Unmarshal(resp.Body(), &containersResp); err != nil {
		return nil, fmt.Errorf("failed to parse containers response: %w", err)
	}

	return containersResp, nil
}

// StartExec starts an exec session for a project or addon container
func (h *HttpClient) StartExec(token string, req *models.ExecRequest) (*models.ExecResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	if req.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}

	if req.ServiceName == "" {
		return nil, errors.New("service name is required")
	}

	if len(req.Command) == 0 {
		return nil, errors.New("command is required")
	}

	// Determine endpoint based on whether it's for project or addon
	endpoint := fmt.Sprintf("/projects/%s/exec", req.ProjectID)
	if req.AddonID != "" {
		endpoint = fmt.Sprintf("/projects/%s/addons/%s/exec", req.ProjectID, req.AddonID)
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to start exec session: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project, addon, or service not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var execResp *models.ExecResponse
	if err := json.Unmarshal(resp.Body(), &execResp); err != nil {
		return nil, fmt.Errorf("failed to parse exec response: %w", err)
	}

	return execResp, nil
}

// StartShell starts a shell session for a project or addon container
func (h *HttpClient) StartShell(token string, req *models.ShellRequest) (*models.ShellResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	if req.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}

	if req.ServiceName == "" {
		return nil, errors.New("service name is required")
	}

	// Determine endpoint based on whether it's for project or addon
	endpoint := fmt.Sprintf("/projects/%s/shell", req.ProjectID)
	if req.AddonID != "" {
		endpoint = fmt.Sprintf("/projects/%s/addons/%s/shell", req.ProjectID, req.AddonID)
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to start shell session: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project, addon, or service not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var shellResp *models.ShellResponse
	if err := json.Unmarshal(resp.Body(), &shellResp); err != nil {
		return nil, fmt.Errorf("failed to parse shell response: %w", err)
	}

	return shellResp, nil
}

// GetAddons retrieves a list of addons
func (h *HttpClient) GetAddons(token string) (*models.AddonListResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/addons")
	if err != nil {
		return nil, fmt.Errorf("failed to get addons: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var addonsResp *models.AddonListResponse
	if err := json.Unmarshal(resp.Body(), &addonsResp); err != nil {
		return nil, fmt.Errorf("failed to parse addons response: %w", err)
	}

	return addonsResp, nil
}

// GetAddon retrieves a specific addon by ID
func (h *HttpClient) GetAddon(token string, addonID string) (*models.Addon, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(addonID) == "" {
		return nil, errors.New("addon ID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/addons/" + addonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get addon: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("addon not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var addon *models.Addon
	if err := json.Unmarshal(resp.Body(), &addon); err != nil {
		return nil, fmt.Errorf("failed to parse addon response: %w", err)
	}

	return addon, nil
}

// GetAddonDeployments retrieves a list of addon deployments
func (h *HttpClient) GetAddonDeployments(token string, projectID string) ([]models.AddonDeployment, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(projectID) == "" {
		return nil, errors.New("project ID is required")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/addons/deployments")
	if err != nil {
		return nil, fmt.Errorf("failed to get addon deployments: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var deploymentsResp *models.AddonDeploymentsResponse
	if err := json.Unmarshal(resp.Body(), &deploymentsResp); err != nil {
		return nil, fmt.Errorf("failed to parse deployments response: %w", err)
	}

	return deploymentsResp.Deployments, nil
}

// DeleteAddonDeployment deletes an addon deployment
func (h *HttpClient) DeleteAddonDeployment(token string, deploymentID string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token is empty")
	}

	if strings.TrimSpace(deploymentID) == "" {
		return errors.New("deployment ID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Delete("/addons/deployments/" + deploymentID)
	if err != nil {
		return fmt.Errorf("failed to delete addon deployment: %w", err)
	}

	if resp.StatusCode() == 401 {
		return ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("deployment not found")
	}

	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.String())
	}

	return nil
}

// GetServers retrieves all servers for the authenticated user
func (h *HttpClient) GetServers(token string) (*models.ServersResponse, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	// Validate token before making the request
	if err := h.validateToken(token); err != nil {
		return nil, err
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/servers")
	if err != nil {
		return nil, fmt.Errorf("failed to get servers: %w", err)
	}

	if resp.StatusCode() == 401 {
		if authErr := detectAuthError(resp); authErr != nil {
			return nil, authErr
		}
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	// Handle empty response body
	body := resp.Body()
	if len(body) == 0 {
		// Return empty servers response when API returns empty body
		return &models.ServersResponse{
			Servers: []models.Server{},
			Total:   0,
			Page:    1,
			PerPage: 10,
		}, nil
	}

	var serversResp *models.ServersResponse
	if err := json.Unmarshal(body, &serversResp); err != nil {
		return nil, fmt.Errorf("failed to parse servers response: %w", err)
	}

	return serversResp, nil
}

// GetServer retrieves a specific server by ID
func (h *HttpClient) GetServer(token string, serverID string) (*models.Server, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(serverID) == "" {
		return nil, errors.New("server ID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/servers/" + serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("server not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var server *models.Server
	if err := json.Unmarshal(resp.Body(), &server); err != nil {
		return nil, fmt.Errorf("failed to parse server response: %w", err)
	}

	return server, nil
}

// CreateServer creates a new server
func (h *HttpClient) CreateServer(token string, req *models.ServerCreateRequest) (*models.Server, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("server name is required")
	}

	if strings.TrimSpace(req.Type) == "" {
		return nil, errors.New("server type is required")
	}

	if strings.TrimSpace(req.Region) == "" {
		return nil, errors.New("server region is required")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post("/servers")
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var server *models.Server
	if err := json.Unmarshal(resp.Body(), &server); err != nil {
		return nil, fmt.Errorf("failed to parse server response: %w", err)
	}

	return server, nil
}

// UpdateServer updates an existing server
func (h *HttpClient) UpdateServer(token string, serverID string, req *models.ServerUpdateRequest) (*models.Server, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(serverID) == "" {
		return nil, errors.New("server ID is empty")
	}

	if req == nil {
		return nil, errors.New("request is nil")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Put("/servers/" + serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to update server: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("server not found")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var server *models.Server
	if err := json.Unmarshal(resp.Body(), &server); err != nil {
		return nil, fmt.Errorf("failed to parse server response: %w", err)
	}

	return server, nil
}

// DeleteServer deletes a server
func (h *HttpClient) DeleteServer(token string, serverID string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token is empty")
	}

	if strings.TrimSpace(serverID) == "" {
		return errors.New("server ID is empty")
	}

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Delete("/servers/" + serverID)
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	if resp.StatusCode() == 401 {
		return ErrInvalidToken
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("server not found")
	}

	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.String())
	}

	return nil
}
