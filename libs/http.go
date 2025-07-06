package libs

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/go-resty/resty/v2"
)

var (
	ErrInvalidToken           = errors.New("invalid token")
	ErrVerificationFailed     = errors.New("token verification failed")
	PIPEOPS_CONTROL_PLANE_API = ""
)

func init() {
	PIPEOPS_CONTROL_PLANE_API = os.Getenv("PIPEOPS_API_URL")
	if PIPEOPS_CONTROL_PLANE_API == "" {
		PIPEOPS_CONTROL_PLANE_API = "https://api.pipeops.io" // Default API URL
	}
}

type HttpClients interface {
	VerifyToken(token string, operatorID string) (*models.PipeOpsTokenVerificationResponse, error)
	GetProjects(token string) (*models.ProjectsResponse, error)
	GetProject(token string, projectID string) (*models.Project, error)
	CreateProject(token string, req *models.ProjectCreateRequest) (*models.Project, error)
	UpdateProject(token string, projectID string, req *models.ProjectUpdateRequest) (*models.Project, error)
	DeleteProject(token string, projectID string) error
	GetLogs(token string, req *models.LogsRequest) (*models.LogsResponse, error)
	StreamLogs(token string, req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error
	GetServices(token string, projectID string, addonID string) (*models.ListServicesResponse, error)
	StartProxy(token string, req *models.ProxyRequest) (*models.ProxyResponse, error)
	GetContainers(token string, projectID string, addonID string) (*models.ListContainersResponse, error)
	StartExec(token string, req *models.ExecRequest) (*models.ExecResponse, error)
	StartShell(token string, req *models.ShellRequest) (*models.ShellResponse, error)
}

type HttpClient struct {
	client *resty.Client
}

func NewHttpClient() HttpClients {
	r := resty.New()

	// Enable debug mode if environment variable is set
	if os.Getenv("PIPEOPS_DEBUG") == "true" {
		r.Debug = true
	}

	URL := strings.TrimSpace(PIPEOPS_CONTROL_PLANE_API)
	r.SetBaseURL(URL)

	return &HttpClient{
		client: r,
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

	var respData *models.PipeOpsTokenVerificationResponse
	if err := json.Unmarshal(resp.Body(), &respData); err != nil {
		return nil, err
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

	resp, err := h.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get("/projects")
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	if resp.StatusCode() == 401 {
		return nil, ErrInvalidToken
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	var projectsResp *models.ProjectsResponse
	if err := json.Unmarshal(resp.Body(), &projectsResp); err != nil {
		return nil, fmt.Errorf("failed to parse projects response: %w", err)
	}

	return projectsResp, nil
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
