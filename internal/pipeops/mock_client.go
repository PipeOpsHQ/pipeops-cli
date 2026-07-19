package pipeops

import (
	"context"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

// MockClient is a mock implementation of ClientAPI
type MockClient struct {
	IsAuthenticatedFunc              func() bool
	LoadConfigFunc                   func() error
	SaveConfigFunc                   func() error
	GetConfigFunc                    func() *config.Config
	GetProjectsFunc                  func() (*models.ProjectsResponse, error)
	GetProjectFunc                   func(projectID string) (*models.Project, error)
	CreateProjectFunc                func(req *models.ProjectCreateRequest) (*models.Project, error)
	UpdateProjectFunc                func(projectID string, req *models.ProjectUpdateRequest) (*models.Project, error)
	DeleteProjectFunc                func(projectID string) error
	DeployProjectFunc                func(projectID string) error
	RestartProjectFunc               func(projectID string) error
	StopProjectFunc                  func(projectID string) error
	GetProjectEnvVariablesFunc       func(projectID string) ([]sdk.EnvVariable, error)
	UpdateProjectEnvVariablesFunc    func(projectID string, envVars []sdk.EnvVariable) ([]sdk.EnvVariable, error)
	ListProjectDeploymentsFunc       func(projectID string, opts *sdk.ProjectDeploymentListOptions) (*sdk.ProjectDeploymentsResponse, error)
	ListProjectDeploymentHistoryFunc func(projectID string, opts *sdk.ProjectDeploymentHistoryOptions) (*sdk.ProjectDeploymentHistoryResponse, error)
	GetLogsFunc                      func(req *models.LogsRequest) (*models.LogsResponse, error)
	StreamLogsFunc                   func(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error
	GetServicesFunc                  func(projectID string, addonID string) (*models.ListServicesResponse, error)
	StartProxyFunc                   func(req *models.ProxyRequest) (*models.ProxyResponse, error)
	GetContainersFunc                func(projectID string, addonID string) (*models.ListContainersResponse, error)
	StartExecFunc                    func(req *models.ExecRequest) (*models.ExecResponse, error)
	StartShellFunc                   func(req *models.ShellRequest) (*models.ShellResponse, error)
	GetAddonsFunc                    func() (*models.AddonListResponse, error)
	GetAddonFunc                     func(addonID string) (*models.Addon, error)
	DeployAddonFunc                  func(req *sdk.DeployAddOnRequest) (*models.AddonDeployment, error)
	GetAddonDeploymentsFunc          func() ([]models.AddonDeployment, error)
	GetAddonDeploymentFunc           func(deploymentID string) (*models.AddonDeployment, error)
	DeleteAddonDeploymentFunc        func(deploymentID string) error
	ListAddonCategoriesFunc          func() ([]sdk.AddOnCategory, error)
	GetAddonDeploymentSessionFunc    func(sessionID string) (map[string]interface{}, error)
	ViewAddonDeploymentConfigsFunc   func(deploymentID string) (map[string]interface{}, error)
	GetServersFunc                   func() (*models.ServersResponse, error)
	GetServerFunc                    func(serverID string) (*models.Server, error)
	GetServerConnectionFunc          func(serverID string) (map[string]interface{}, error)
	GetServerCostAllocationFunc      func(serverID string) (map[string]interface{}, error)
	CreateServerFunc                 func(req *models.ServerCreateRequest) (*models.Server, error)
	UpdateServerFunc                 func(serverID string, req *models.ServerUpdateRequest) (*models.Server, error)
	DeleteServerFunc                 func(serverID string) error
	VerifyTokenFunc                  func() (*models.PipeOpsTokenVerificationResponse, error)
	GetWorkspacesFunc                func(ctx context.Context) ([]sdk.Workspace, error)
	GetWorkspaceFunc                 func(ctx context.Context, workspaceID string) (*sdk.Workspace, error)
	CreateWorkspaceFunc              func(ctx context.Context, req *sdk.CreateWorkspaceRequest) (*sdk.Workspace, error)
	UpdateWorkspaceFunc              func(ctx context.Context, workspaceID string, req *sdk.UpdateWorkspaceRequest) (*sdk.Workspace, error)
	DeleteWorkspaceFunc              func(ctx context.Context, workspaceID string) error
	ListEnvironmentsFunc             func(ctx context.Context) ([]sdk.Environment, error)
	GetEnvironmentFunc               func(ctx context.Context, environmentID string) (*sdk.Environment, error)
	CreateEnvironmentFunc            func(ctx context.Context, req *sdk.CreateEnvironmentRequest) (*sdk.Environment, error)
	UpdateEnvironmentFunc            func(ctx context.Context, environmentID string, req *sdk.UpdateEnvironmentRequest) (*sdk.Environment, error)
	DeleteEnvironmentFunc            func(ctx context.Context, environmentID string) error
	SetEnvironmentVariablesFunc      func(ctx context.Context, environmentID string, envVars []sdk.EnvVariable) error
	ListServiceAccountTokensFunc     func(ctx context.Context) ([]sdk.ServiceAccountToken, error)
	GetServiceAccountTokenFunc       func(ctx context.Context, tokenID string) (*sdk.ServiceAccountToken, error)
	CreateServiceAccountTokenFunc    func(ctx context.Context, req *sdk.ServiceAccountTokenRequest) (*sdk.ServiceAccountToken, error)
	UpdateServiceAccountTokenFunc    func(ctx context.Context, tokenID string, req *sdk.ServiceAccountTokenUpdateRequest) (*sdk.ServiceAccountToken, error)
	RevokeServiceAccountTokenFunc    func(ctx context.Context, tokenID string) error
	ListVolumesFunc                  func(ctx context.Context, opts *sdk.VolumeListOptions) (*sdk.VolumeListResponse, error)
	GetVolumeFunc                    func(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) (*sdk.Volume, error)
	RemountVolumeFunc                func(ctx context.Context, volumeUUID string, body *sdk.RemountVolumeRequest, opts *sdk.VolumeListOptions) (*sdk.RemountVolumeResponse, error)
	DeleteVolumeFunc                 func(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) error
	StartVolumeExportFunc            func(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) (*sdk.VolumeExportResponse, error)
	GetVolumeExportFunc              func(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) (*sdk.VolumeExportResponse, error)
	ListAddonBackupsFunc             func(ctx context.Context, deploymentUID string) (*sdk.AddonBackupListResponse, error)
	StartAddonBackupExportFunc       func(ctx context.Context, deploymentUID string, body *sdk.AddonBackupExportRequest) (*sdk.AddonBackupExportResponse, error)
	GetAddonBackupExportFunc         func(ctx context.Context, deploymentUID, exportID string) (*sdk.AddonBackupExportResponse, error)
	ListGitOpsFunc          func(ctx context.Context, opts *sdk.GitOpsListOptions) (*sdk.GitOpsListResponse, error)
	GetGitOpsFunc           func(ctx context.Context, uuid string) (*sdk.GitOpsConfig, error)
	CreateGitOpsFunc        func(ctx context.Context, body *sdk.CreateGitOpsConfigRequest) (*sdk.GitOpsConfig, error)
	UpdateGitOpsFunc        func(ctx context.Context, uuid string, body *sdk.UpdateGitOpsConfigRequest) (*sdk.GitOpsConfig, error)
	DeleteGitOpsFunc        func(ctx context.Context, uuid string) error
	TriggerGitOpsSyncFunc   func(ctx context.Context, uuid string, body *sdk.TriggerGitOpsSyncRequest) (*sdk.GitOpsSyncTriggerResponse, error)
	GetGitOpsSyncStatusFunc func(ctx context.Context, uuid string) (*sdk.GitOpsSyncStatusResponse, error)
	GetGitOpsDiffFunc       func(ctx context.Context, uuid string) (*sdk.GitOpsDiffResponse, error)
	GetGitOpsHistoryFunc    func(ctx context.Context, uuid string, opts *sdk.GitOpsListOptions) (*sdk.GitOpsSyncHistoryResponse, error)
	ListProjectGroupsFunc           func(ctx context.Context, opts *sdk.ProjectGroupListOptions) (*sdk.ProjectGroupListResponse, error)
	GetProjectGroupFunc             func(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroup, error)
	CreateProjectGroupFunc          func(ctx context.Context, body *sdk.CreateProjectGroupRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroup, error)
	UpdateProjectGroupFunc          func(ctx context.Context, uuid string, body *sdk.UpdateProjectGroupRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroup, error)
	DeleteProjectGroupFunc          func(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) error
	AttachProjectGroupMemberFunc    func(ctx context.Context, uuid string, body *sdk.AttachProjectGroupMemberRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupAttachResponse, error)
	DetachProjectGroupMemberFunc    func(ctx context.Context, uuid, memberType, memberUUID string, opts *sdk.ProjectGroupDetachOptions) error
	GetProjectGroupTopologyFunc     func(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupTopologyResponse, error)
	GetProjectGroupSharedEnvFunc    func(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupSharedEnvResponse, error)
	PutProjectGroupSharedEnvFunc    func(ctx context.Context, uuid string, body *sdk.UpsertProjectGroupSharedEnvRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupSharedEnvResponse, error)
	InjectProjectGroupSharedEnvFunc func(ctx context.Context, uuid string, body *sdk.InjectProjectGroupSharedEnvRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupInjectSharedEnvResponse, error)
	ConnectProjectGroupServicesFunc func(ctx context.Context, uuid string, body *sdk.ConnectProjectGroupServicesRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupConnectResponse, error)
	RedeployProjectGroupAppsFunc    func(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupRedeployAppsResponse, error)
	ResolveProjectGroupMemberFunc   func(ctx context.Context, opts *sdk.ProjectGroupResolveOptions) (*sdk.ProjectGroupResolveResponse, error)
	ListProjectGroupCandidatesFunc  func(ctx context.Context, opts *sdk.ProjectGroupCandidatesOptions) (*sdk.ProjectGroupCandidatesResponse, error)
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

func (m *MockClient) DeployProject(projectID string) error {
	if m.DeployProjectFunc != nil {
		return m.DeployProjectFunc(projectID)
	}
	return nil
}

func (m *MockClient) RestartProject(projectID string) error {
	if m.RestartProjectFunc != nil {
		return m.RestartProjectFunc(projectID)
	}
	return nil
}

func (m *MockClient) StopProject(projectID string) error {
	if m.StopProjectFunc != nil {
		return m.StopProjectFunc(projectID)
	}
	return nil
}

func (m *MockClient) GetProjectEnvVariables(projectID string) ([]sdk.EnvVariable, error) {
	if m.GetProjectEnvVariablesFunc != nil {
		return m.GetProjectEnvVariablesFunc(projectID)
	}
	return nil, nil
}

func (m *MockClient) UpdateProjectEnvVariables(projectID string, envVars []sdk.EnvVariable) ([]sdk.EnvVariable, error) {
	if m.UpdateProjectEnvVariablesFunc != nil {
		return m.UpdateProjectEnvVariablesFunc(projectID, envVars)
	}
	return envVars, nil
}

func (m *MockClient) ListProjectDeployments(projectID string, opts *sdk.ProjectDeploymentListOptions) (*sdk.ProjectDeploymentsResponse, error) {
	if m.ListProjectDeploymentsFunc != nil {
		return m.ListProjectDeploymentsFunc(projectID, opts)
	}
	return &sdk.ProjectDeploymentsResponse{}, nil
}

func (m *MockClient) ListProjectDeploymentHistory(projectID string, opts *sdk.ProjectDeploymentHistoryOptions) (*sdk.ProjectDeploymentHistoryResponse, error) {
	if m.ListProjectDeploymentHistoryFunc != nil {
		return m.ListProjectDeploymentHistoryFunc(projectID, opts)
	}
	return &sdk.ProjectDeploymentHistoryResponse{}, nil
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

func (m *MockClient) DeployAddon(req *sdk.DeployAddOnRequest) (*models.AddonDeployment, error) {
	if m.DeployAddonFunc != nil {
		return m.DeployAddonFunc(req)
	}
	return &models.AddonDeployment{}, nil
}

func (m *MockClient) GetAddonDeployments() ([]models.AddonDeployment, error) {
	if m.GetAddonDeploymentsFunc != nil {
		return m.GetAddonDeploymentsFunc()
	}
	return nil, nil
}

func (m *MockClient) GetAddonDeployment(deploymentID string) (*models.AddonDeployment, error) {
	if m.GetAddonDeploymentFunc != nil {
		return m.GetAddonDeploymentFunc(deploymentID)
	}
	return &models.AddonDeployment{}, nil
}

func (m *MockClient) DeleteAddonDeployment(deploymentID string) error {
	if m.DeleteAddonDeploymentFunc != nil {
		return m.DeleteAddonDeploymentFunc(deploymentID)
	}
	return nil
}

func (m *MockClient) ListAddonCategories() ([]sdk.AddOnCategory, error) {
	if m.ListAddonCategoriesFunc != nil {
		return m.ListAddonCategoriesFunc()
	}
	return nil, nil
}

func (m *MockClient) GetAddonDeploymentSession(sessionID string) (map[string]interface{}, error) {
	if m.GetAddonDeploymentSessionFunc != nil {
		return m.GetAddonDeploymentSessionFunc(sessionID)
	}
	return map[string]interface{}{}, nil
}

func (m *MockClient) ViewAddonDeploymentConfigs(deploymentID string) (map[string]interface{}, error) {
	if m.ViewAddonDeploymentConfigsFunc != nil {
		return m.ViewAddonDeploymentConfigsFunc(deploymentID)
	}
	return map[string]interface{}{}, nil
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

func (m *MockClient) GetServerConnection(serverID string) (map[string]interface{}, error) {
	if m.GetServerConnectionFunc != nil {
		return m.GetServerConnectionFunc(serverID)
	}
	return map[string]interface{}{}, nil
}

func (m *MockClient) GetServerCostAllocation(serverID string) (map[string]interface{}, error) {
	if m.GetServerCostAllocationFunc != nil {
		return m.GetServerCostAllocationFunc(serverID)
	}
	return map[string]interface{}{}, nil
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

func (m *MockClient) GetWorkspace(ctx context.Context, workspaceID string) (*sdk.Workspace, error) {
	if m.GetWorkspaceFunc != nil {
		return m.GetWorkspaceFunc(ctx, workspaceID)
	}
	return &sdk.Workspace{}, nil
}

func (m *MockClient) CreateWorkspace(ctx context.Context, req *sdk.CreateWorkspaceRequest) (*sdk.Workspace, error) {
	if m.CreateWorkspaceFunc != nil {
		return m.CreateWorkspaceFunc(ctx, req)
	}
	return &sdk.Workspace{}, nil
}

func (m *MockClient) UpdateWorkspace(ctx context.Context, workspaceID string, req *sdk.UpdateWorkspaceRequest) (*sdk.Workspace, error) {
	if m.UpdateWorkspaceFunc != nil {
		return m.UpdateWorkspaceFunc(ctx, workspaceID, req)
	}
	return &sdk.Workspace{}, nil
}

func (m *MockClient) DeleteWorkspace(ctx context.Context, workspaceID string) error {
	if m.DeleteWorkspaceFunc != nil {
		return m.DeleteWorkspaceFunc(ctx, workspaceID)
	}
	return nil
}

func (m *MockClient) ListEnvironments(ctx context.Context) ([]sdk.Environment, error) {
	if m.ListEnvironmentsFunc != nil {
		return m.ListEnvironmentsFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetEnvironment(ctx context.Context, environmentID string) (*sdk.Environment, error) {
	if m.GetEnvironmentFunc != nil {
		return m.GetEnvironmentFunc(ctx, environmentID)
	}
	return &sdk.Environment{}, nil
}

func (m *MockClient) CreateEnvironment(ctx context.Context, req *sdk.CreateEnvironmentRequest) (*sdk.Environment, error) {
	if m.CreateEnvironmentFunc != nil {
		return m.CreateEnvironmentFunc(ctx, req)
	}
	return &sdk.Environment{}, nil
}

func (m *MockClient) UpdateEnvironment(ctx context.Context, environmentID string, req *sdk.UpdateEnvironmentRequest) (*sdk.Environment, error) {
	if m.UpdateEnvironmentFunc != nil {
		return m.UpdateEnvironmentFunc(ctx, environmentID, req)
	}
	return &sdk.Environment{}, nil
}

func (m *MockClient) DeleteEnvironment(ctx context.Context, environmentID string) error {
	if m.DeleteEnvironmentFunc != nil {
		return m.DeleteEnvironmentFunc(ctx, environmentID)
	}
	return nil
}

func (m *MockClient) SetEnvironmentVariables(ctx context.Context, environmentID string, envVars []sdk.EnvVariable) error {
	if m.SetEnvironmentVariablesFunc != nil {
		return m.SetEnvironmentVariablesFunc(ctx, environmentID, envVars)
	}
	return nil
}

func (m *MockClient) ListServiceAccountTokens(ctx context.Context) ([]sdk.ServiceAccountToken, error) {
	if m.ListServiceAccountTokensFunc != nil {
		return m.ListServiceAccountTokensFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetServiceAccountToken(ctx context.Context, tokenID string) (*sdk.ServiceAccountToken, error) {
	if m.GetServiceAccountTokenFunc != nil {
		return m.GetServiceAccountTokenFunc(ctx, tokenID)
	}
	return &sdk.ServiceAccountToken{}, nil
}

func (m *MockClient) CreateServiceAccountToken(ctx context.Context, req *sdk.ServiceAccountTokenRequest) (*sdk.ServiceAccountToken, error) {
	if m.CreateServiceAccountTokenFunc != nil {
		return m.CreateServiceAccountTokenFunc(ctx, req)
	}
	return &sdk.ServiceAccountToken{}, nil
}

func (m *MockClient) UpdateServiceAccountToken(ctx context.Context, tokenID string, req *sdk.ServiceAccountTokenUpdateRequest) (*sdk.ServiceAccountToken, error) {
	if m.UpdateServiceAccountTokenFunc != nil {
		return m.UpdateServiceAccountTokenFunc(ctx, tokenID, req)
	}
	return &sdk.ServiceAccountToken{}, nil
}

func (m *MockClient) RevokeServiceAccountToken(ctx context.Context, tokenID string) error {
	if m.RevokeServiceAccountTokenFunc != nil {
		return m.RevokeServiceAccountTokenFunc(ctx, tokenID)
	}
	return nil
}

func (m *MockClient) ListVolumes(ctx context.Context, opts *sdk.VolumeListOptions) (*sdk.VolumeListResponse, error) {
	if m.ListVolumesFunc != nil {
		return m.ListVolumesFunc(ctx, opts)
	}
	return &sdk.VolumeListResponse{}, nil
}

func (m *MockClient) GetVolume(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) (*sdk.Volume, error) {
	if m.GetVolumeFunc != nil {
		return m.GetVolumeFunc(ctx, volumeUUID, opts)
	}
	return &sdk.Volume{}, nil
}

func (m *MockClient) RemountVolume(ctx context.Context, volumeUUID string, body *sdk.RemountVolumeRequest, opts *sdk.VolumeListOptions) (*sdk.RemountVolumeResponse, error) {
	if m.RemountVolumeFunc != nil {
		return m.RemountVolumeFunc(ctx, volumeUUID, body, opts)
	}
	return &sdk.RemountVolumeResponse{}, nil
}

func (m *MockClient) DeleteVolume(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) error {
	if m.DeleteVolumeFunc != nil {
		return m.DeleteVolumeFunc(ctx, volumeUUID, opts)
	}
	return nil
}

func (m *MockClient) StartVolumeExport(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) (*sdk.VolumeExportResponse, error) {
	if m.StartVolumeExportFunc != nil {
		return m.StartVolumeExportFunc(ctx, volumeUUID, opts)
	}
	return &sdk.VolumeExportResponse{}, nil
}

func (m *MockClient) GetVolumeExport(ctx context.Context, volumeUUID string, opts *sdk.VolumeListOptions) (*sdk.VolumeExportResponse, error) {
	if m.GetVolumeExportFunc != nil {
		return m.GetVolumeExportFunc(ctx, volumeUUID, opts)
	}
	return &sdk.VolumeExportResponse{}, nil
}

func (m *MockClient) ListAddonBackups(ctx context.Context, deploymentUID string) (*sdk.AddonBackupListResponse, error) {
	if m.ListAddonBackupsFunc != nil {
		return m.ListAddonBackupsFunc(ctx, deploymentUID)
	}
	return &sdk.AddonBackupListResponse{}, nil
}

func (m *MockClient) StartAddonBackupExport(ctx context.Context, deploymentUID string, body *sdk.AddonBackupExportRequest) (*sdk.AddonBackupExportResponse, error) {
	if m.StartAddonBackupExportFunc != nil {
		return m.StartAddonBackupExportFunc(ctx, deploymentUID, body)
	}
	return &sdk.AddonBackupExportResponse{}, nil
}

func (m *MockClient) GetAddonBackupExport(ctx context.Context, deploymentUID, exportID string) (*sdk.AddonBackupExportResponse, error) {
	if m.GetAddonBackupExportFunc != nil {
		return m.GetAddonBackupExportFunc(ctx, deploymentUID, exportID)
	}
	return &sdk.AddonBackupExportResponse{}, nil
}

func (m *MockClient) ListGitOps(ctx context.Context, opts *sdk.GitOpsListOptions) (*sdk.GitOpsListResponse, error) {
	if m.ListGitOpsFunc != nil {
		return m.ListGitOpsFunc(ctx, opts)
	}
	return &sdk.GitOpsListResponse{}, nil
}


func (m *MockClient) GetGitOps(ctx context.Context, uuid string) (*sdk.GitOpsConfig, error) {
	if m.GetGitOpsFunc != nil {
		return m.GetGitOpsFunc(ctx, uuid)
	}
	return &sdk.GitOpsConfig{}, nil
}


func (m *MockClient) CreateGitOps(ctx context.Context, body *sdk.CreateGitOpsConfigRequest) (*sdk.GitOpsConfig, error) {
	if m.CreateGitOpsFunc != nil {
		return m.CreateGitOpsFunc(ctx, body)
	}
	return &sdk.GitOpsConfig{}, nil
}


func (m *MockClient) UpdateGitOps(ctx context.Context, uuid string, body *sdk.UpdateGitOpsConfigRequest) (*sdk.GitOpsConfig, error) {
	if m.UpdateGitOpsFunc != nil {
		return m.UpdateGitOpsFunc(ctx, uuid, body)
	}
	return &sdk.GitOpsConfig{}, nil
}


func (m *MockClient) DeleteGitOps(ctx context.Context, uuid string) error {
	if m.DeleteGitOpsFunc != nil {
		return m.DeleteGitOpsFunc(ctx, uuid)
	}
	return nil
}


func (m *MockClient) TriggerGitOpsSync(ctx context.Context, uuid string, body *sdk.TriggerGitOpsSyncRequest) (*sdk.GitOpsSyncTriggerResponse, error) {
	if m.TriggerGitOpsSyncFunc != nil {
		return m.TriggerGitOpsSyncFunc(ctx, uuid, body)
	}
	return &sdk.GitOpsSyncTriggerResponse{}, nil
}


func (m *MockClient) GetGitOpsSyncStatus(ctx context.Context, uuid string) (*sdk.GitOpsSyncStatusResponse, error) {
	if m.GetGitOpsSyncStatusFunc != nil {
		return m.GetGitOpsSyncStatusFunc(ctx, uuid)
	}
	return &sdk.GitOpsSyncStatusResponse{}, nil
}


func (m *MockClient) GetGitOpsDiff(ctx context.Context, uuid string) (*sdk.GitOpsDiffResponse, error) {
	if m.GetGitOpsDiffFunc != nil {
		return m.GetGitOpsDiffFunc(ctx, uuid)
	}
	return &sdk.GitOpsDiffResponse{}, nil
}


func (m *MockClient) GetGitOpsHistory(ctx context.Context, uuid string, opts *sdk.GitOpsListOptions) (*sdk.GitOpsSyncHistoryResponse, error) {
	if m.GetGitOpsHistoryFunc != nil {
		return m.GetGitOpsHistoryFunc(ctx, uuid, opts)
	}
	return &sdk.GitOpsSyncHistoryResponse{}, nil
}


func (m *MockClient) ListProjectGroups(ctx context.Context, opts *sdk.ProjectGroupListOptions) (*sdk.ProjectGroupListResponse, error) {
	if m.ListProjectGroupsFunc != nil {
		return m.ListProjectGroupsFunc(ctx, opts)
	}
	return &sdk.ProjectGroupListResponse{}, nil
}


func (m *MockClient) GetProjectGroup(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroup, error) {
	if m.GetProjectGroupFunc != nil {
		return m.GetProjectGroupFunc(ctx, uuid, opts)
	}
	return &sdk.ProjectGroup{}, nil
}


func (m *MockClient) CreateProjectGroup(ctx context.Context, body *sdk.CreateProjectGroupRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroup, error) {
	if m.CreateProjectGroupFunc != nil {
		return m.CreateProjectGroupFunc(ctx, body, opts)
	}
	return &sdk.ProjectGroup{}, nil
}


func (m *MockClient) UpdateProjectGroup(ctx context.Context, uuid string, body *sdk.UpdateProjectGroupRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroup, error) {
	if m.UpdateProjectGroupFunc != nil {
		return m.UpdateProjectGroupFunc(ctx, uuid, body, opts)
	}
	return &sdk.ProjectGroup{}, nil
}


func (m *MockClient) DeleteProjectGroup(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) error {
	if m.DeleteProjectGroupFunc != nil {
		return m.DeleteProjectGroupFunc(ctx, uuid, opts)
	}
	return nil
}


func (m *MockClient) AttachProjectGroupMember(ctx context.Context, uuid string, body *sdk.AttachProjectGroupMemberRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupAttachResponse, error) {
	if m.AttachProjectGroupMemberFunc != nil {
		return m.AttachProjectGroupMemberFunc(ctx, uuid, body, opts)
	}
	return &sdk.ProjectGroupAttachResponse{}, nil
}


func (m *MockClient) DetachProjectGroupMember(ctx context.Context, uuid, memberType, memberUUID string, opts *sdk.ProjectGroupDetachOptions) error {
	if m.DetachProjectGroupMemberFunc != nil {
		return m.DetachProjectGroupMemberFunc(ctx, uuid, memberType, memberUUID, opts)
	}
	return nil
}


func (m *MockClient) GetProjectGroupTopology(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupTopologyResponse, error) {
	if m.GetProjectGroupTopologyFunc != nil {
		return m.GetProjectGroupTopologyFunc(ctx, uuid, opts)
	}
	return &sdk.ProjectGroupTopologyResponse{}, nil
}


func (m *MockClient) GetProjectGroupSharedEnv(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupSharedEnvResponse, error) {
	if m.GetProjectGroupSharedEnvFunc != nil {
		return m.GetProjectGroupSharedEnvFunc(ctx, uuid, opts)
	}
	return &sdk.ProjectGroupSharedEnvResponse{}, nil
}


func (m *MockClient) PutProjectGroupSharedEnv(ctx context.Context, uuid string, body *sdk.UpsertProjectGroupSharedEnvRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupSharedEnvResponse, error) {
	if m.PutProjectGroupSharedEnvFunc != nil {
		return m.PutProjectGroupSharedEnvFunc(ctx, uuid, body, opts)
	}
	return &sdk.ProjectGroupSharedEnvResponse{}, nil
}


func (m *MockClient) InjectProjectGroupSharedEnv(ctx context.Context, uuid string, body *sdk.InjectProjectGroupSharedEnvRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupInjectSharedEnvResponse, error) {
	if m.InjectProjectGroupSharedEnvFunc != nil {
		return m.InjectProjectGroupSharedEnvFunc(ctx, uuid, body, opts)
	}
	return &sdk.ProjectGroupInjectSharedEnvResponse{}, nil
}


func (m *MockClient) ConnectProjectGroupServices(ctx context.Context, uuid string, body *sdk.ConnectProjectGroupServicesRequest, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupConnectResponse, error) {
	if m.ConnectProjectGroupServicesFunc != nil {
		return m.ConnectProjectGroupServicesFunc(ctx, uuid, body, opts)
	}
	return &sdk.ProjectGroupConnectResponse{}, nil
}


func (m *MockClient) RedeployProjectGroupApps(ctx context.Context, uuid string, opts *sdk.ProjectGroupWorkspaceOptions) (*sdk.ProjectGroupRedeployAppsResponse, error) {
	if m.RedeployProjectGroupAppsFunc != nil {
		return m.RedeployProjectGroupAppsFunc(ctx, uuid, opts)
	}
	return &sdk.ProjectGroupRedeployAppsResponse{}, nil
}


func (m *MockClient) ResolveProjectGroupMember(ctx context.Context, opts *sdk.ProjectGroupResolveOptions) (*sdk.ProjectGroupResolveResponse, error) {
	if m.ResolveProjectGroupMemberFunc != nil {
		return m.ResolveProjectGroupMemberFunc(ctx, opts)
	}
	return &sdk.ProjectGroupResolveResponse{}, nil
}


func (m *MockClient) ListProjectGroupCandidates(ctx context.Context, opts *sdk.ProjectGroupCandidatesOptions) (*sdk.ProjectGroupCandidatesResponse, error) {
	if m.ListProjectGroupCandidatesFunc != nil {
		return m.ListProjectGroupCandidatesFunc(ctx, opts)
	}
	return &sdk.ProjectGroupCandidatesResponse{}, nil
}
