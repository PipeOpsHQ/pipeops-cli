package pipeops

import (
	"context"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

// MockClient is a mock implementation of ClientAPI
type MockClient struct {
	IsAuthenticatedFunc       func() bool
	LoadConfigFunc            func() error
	SaveConfigFunc            func() error
	GetConfigFunc             func() *config.Config
	GetProjectsFunc           func() (*models.ProjectsResponse, error)
	GetProjectFunc            func(projectID string) (*models.Project, error)
	CreateProjectFunc         func(req *models.ProjectCreateRequest) (*models.Project, error)
	UpdateProjectFunc         func(projectID string, req *models.ProjectUpdateRequest) (*models.Project, error)
	DeleteProjectFunc         func(projectID string) error
	GetLogsFunc               func(req *models.LogsRequest) (*models.LogsResponse, error)
	StreamLogsFunc            func(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error
	GetServicesFunc           func(projectID string, addonID string) (*models.ListServicesResponse, error)
	StartProxyFunc            func(req *models.ProxyRequest) (*models.ProxyResponse, error)
	GetContainersFunc         func(projectID string, addonID string) (*models.ListContainersResponse, error)
	StartExecFunc             func(req *models.ExecRequest) (*models.ExecResponse, error)
	StartShellFunc            func(req *models.ShellRequest) (*models.ShellResponse, error)
	GetAddonsFunc             func() (*models.AddonListResponse, error)
	GetAddonFunc              func(addonID string) (*models.Addon, error)
	DeployAddonFunc           func(req *models.AddonDeployRequest) (*models.AddonDeployResponse, error)
	GetAddonDeploymentsFunc   func(projectID string) ([]models.AddonDeployment, error)
	DeleteAddonDeploymentFunc func(deploymentID string) error
	GetServersFunc            func() (*models.ServersResponse, error)
	GetServerFunc             func(serverID string) (*models.Server, error)
	CreateServerFunc          func(req *models.ServerCreateRequest) (*models.Server, error)
	UpdateServerFunc          func(serverID string, req *models.ServerUpdateRequest) (*models.Server, error)
	DeleteServerFunc          func(serverID string) error
	VerifyTokenFunc           func() (*models.PipeOpsTokenVerificationResponse, error)
	GetWorkspacesFunc         func(ctx context.Context) ([]sdk.Workspace, error)
}

func (m *MockClient) IsAuthenticated() bool {
	if m.IsAuthenticatedFunc != nil {
		return m.IsAuthenticatedFunc()
	}
	return true
}

func (m *MockClient) GetProjects() (*models.ProjectsResponse, error) {
	if m.GetProjectsFunc != nil {
		return m.GetProjectsFunc()
	}
	return &models.ProjectsResponse{}, nil
}

func (m *MockClient) GetProject(projectID string) (*models.Project, error) {
	if m.GetProjectFunc != nil {
		return m.GetProjectFunc(projectID)
	}
	return nil, nil
}

func (m *MockClient) CreateProject(req *models.ProjectCreateRequest) (*models.Project, error) {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(req)
	}
	return nil, nil
}

func (m *MockClient) UpdateProject(projectID string, req *models.ProjectUpdateRequest) (*models.Project, error) {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(projectID, req)
	}
	return nil, nil
}

func (m *MockClient) DeleteProject(projectID string) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(projectID)
	}
	return nil
}

func (m *MockClient) GetLogs(req *models.LogsRequest) (*models.LogsResponse, error) {
	if m.GetLogsFunc != nil {
		return m.GetLogsFunc(req)
	}
	return nil, nil
}

func (m *MockClient) StreamLogs(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error {
	if m.StreamLogsFunc != nil {
		return m.StreamLogsFunc(req, callback)
	}
	return nil
}

func (m *MockClient) GetServices(projectID string, addonID string) (*models.ListServicesResponse, error) {
	if m.GetServicesFunc != nil {
		return m.GetServicesFunc(projectID, addonID)
	}
	return nil, nil
}

func (m *MockClient) StartProxy(req *models.ProxyRequest) (*models.ProxyResponse, error) {
	if m.StartProxyFunc != nil {
		return m.StartProxyFunc(req)
	}
	return nil, nil
}

func (m *MockClient) GetContainers(projectID string, addonID string) (*models.ListContainersResponse, error) {
	if m.GetContainersFunc != nil {
		return m.GetContainersFunc(projectID, addonID)
	}
	return nil, nil
}

func (m *MockClient) StartExec(req *models.ExecRequest) (*models.ExecResponse, error) {
	if m.StartExecFunc != nil {
		return m.StartExecFunc(req)
	}
	return nil, nil
}

func (m *MockClient) StartShell(req *models.ShellRequest) (*models.ShellResponse, error) {
	if m.StartShellFunc != nil {
		return m.StartShellFunc(req)
	}
	return nil, nil
}

func (m *MockClient) GetAddons() (*models.AddonListResponse, error) {
	if m.GetAddonsFunc != nil {
		return m.GetAddonsFunc()
	}
	return nil, nil
}

func (m *MockClient) GetAddon(addonID string) (*models.Addon, error) {
	if m.GetAddonFunc != nil {
		return m.GetAddonFunc(addonID)
	}
	return nil, nil
}

func (m *MockClient) DeployAddon(req *models.AddonDeployRequest) (*models.AddonDeployResponse, error) {
	if m.DeployAddonFunc != nil {
		return m.DeployAddonFunc(req)
	}
	return nil, nil
}

func (m *MockClient) GetAddonDeployments(projectID string) ([]models.AddonDeployment, error) {
	if m.GetAddonDeploymentsFunc != nil {
		return m.GetAddonDeploymentsFunc(projectID)
	}
	return nil, nil
}

func (m *MockClient) DeleteAddonDeployment(deploymentID string) error {
	if m.DeleteAddonDeploymentFunc != nil {
		return m.DeleteAddonDeploymentFunc(deploymentID)
	}
	return nil
}

func (m *MockClient) GetServers() (*models.ServersResponse, error) {
	if m.GetServersFunc != nil {
		return m.GetServersFunc()
	}
	return nil, nil
}

func (m *MockClient) GetServer(serverID string) (*models.Server, error) {
	if m.GetServerFunc != nil {
		return m.GetServerFunc(serverID)
	}
	return nil, nil
}

func (m *MockClient) CreateServer(req *models.ServerCreateRequest) (*models.Server, error) {
	if m.CreateServerFunc != nil {
		return m.CreateServerFunc(req)
	}
	return nil, nil
}

func (m *MockClient) UpdateServer(serverID string, req *models.ServerUpdateRequest) (*models.Server, error) {
	if m.UpdateServerFunc != nil {
		return m.UpdateServerFunc(serverID, req)
	}
	return nil, nil
}

func (m *MockClient) DeleteServer(serverID string) error {
	if m.DeleteServerFunc != nil {
		return m.DeleteServerFunc(serverID)
	}
	return nil
}

func (m *MockClient) VerifyToken() (*models.PipeOpsTokenVerificationResponse, error) {
	if m.VerifyTokenFunc != nil {
		return m.VerifyTokenFunc()
	}
	return &models.PipeOpsTokenVerificationResponse{Valid: true}, nil
}

func (m *MockClient) LoadConfig() error {
	if m.LoadConfigFunc != nil {
		return m.LoadConfigFunc()
	}
	return nil
}

func (m *MockClient) SaveConfig() error {
	if m.SaveConfigFunc != nil {
		return m.SaveConfigFunc()
	}
	return nil
}

func (m *MockClient) GetConfig() *config.Config {
	if m.GetConfigFunc != nil {
		return m.GetConfigFunc()
	}
	return config.DefaultConfig()
}

func (m *MockClient) GetWorkspaces(ctx context.Context) ([]sdk.Workspace, error) {
	if m.GetWorkspacesFunc != nil {
		return m.GetWorkspacesFunc(ctx)
	}
	return nil, nil
}
