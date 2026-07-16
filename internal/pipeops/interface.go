package pipeops

import (
	"context"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

// ClientAPI defines the interface for PipeOps API operations
type ClientAPI interface {
	IsAuthenticated() bool
	GetProjects() (*models.ProjectsResponse, error)
	GetProject(projectID string) (*models.Project, error)
	CreateProject(req *models.ProjectCreateRequest) (*models.Project, error)
	UpdateProject(projectID string, req *models.ProjectUpdateRequest) (*models.Project, error)
	DeleteProject(projectID string) error
	DeployProject(projectID string) error
	RestartProject(projectID string) error
	StopProject(projectID string) error
	GetProjectEnvVariables(projectID string) ([]sdk.EnvVariable, error)
	UpdateProjectEnvVariables(projectID string, envVars []sdk.EnvVariable) ([]sdk.EnvVariable, error)
	ListProjectDeployments(projectID string, opts *sdk.ProjectDeploymentListOptions) (*sdk.ProjectDeploymentsResponse, error)
	ListProjectDeploymentHistory(projectID string, opts *sdk.ProjectDeploymentHistoryOptions) (*sdk.ProjectDeploymentHistoryResponse, error)
	GetLogs(req *models.LogsRequest) (*models.LogsResponse, error)
	StreamLogs(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error
	GetServices(projectID string, addonID string) (*models.ListServicesResponse, error)
	StartProxy(req *models.ProxyRequest) (*models.ProxyResponse, error)
	GetContainers(projectID string, addonID string) (*models.ListContainersResponse, error)
	StartExec(req *models.ExecRequest) (*models.ExecResponse, error)
	StartShell(req *models.ShellRequest) (*models.ShellResponse, error)
	GetAddons() (*models.AddonListResponse, error)
	GetAddon(addonID string) (*models.Addon, error)
	DeployAddon(req *sdk.DeployAddOnRequest) (*models.AddonDeployment, error)
	GetAddonDeployments() ([]models.AddonDeployment, error)
	GetAddonDeployment(deploymentID string) (*models.AddonDeployment, error)
	DeleteAddonDeployment(deploymentID string) error
	ListAddonCategories() ([]sdk.AddOnCategory, error)
	GetAddonDeploymentSession(sessionID string) (map[string]interface{}, error)
	ViewAddonDeploymentConfigs(deploymentID string) (map[string]interface{}, error)
	GetServers() (*models.ServersResponse, error)
	GetServer(serverID string) (*models.Server, error)
	GetServerConnection(serverID string) (map[string]interface{}, error)
	GetServerCostAllocation(serverID string) (map[string]interface{}, error)
	CreateServer(req *models.ServerCreateRequest) (*models.Server, error)
	UpdateServer(serverID string, req *models.ServerUpdateRequest) (*models.Server, error)
	DeleteServer(serverID string) error
	VerifyToken() (*models.PipeOpsTokenVerificationResponse, error)
	GetWorkspaces(ctx context.Context) ([]sdk.Workspace, error)
	GetWorkspace(ctx context.Context, workspaceID string) (*sdk.Workspace, error)
	CreateWorkspace(ctx context.Context, req *sdk.CreateWorkspaceRequest) (*sdk.Workspace, error)
	UpdateWorkspace(ctx context.Context, workspaceID string, req *sdk.UpdateWorkspaceRequest) (*sdk.Workspace, error)
	DeleteWorkspace(ctx context.Context, workspaceID string) error
	ListEnvironments(ctx context.Context) ([]sdk.Environment, error)
	GetEnvironment(ctx context.Context, environmentID string) (*sdk.Environment, error)
	CreateEnvironment(ctx context.Context, req *sdk.CreateEnvironmentRequest) (*sdk.Environment, error)
	UpdateEnvironment(ctx context.Context, environmentID string, req *sdk.UpdateEnvironmentRequest) (*sdk.Environment, error)
	DeleteEnvironment(ctx context.Context, environmentID string) error
	SetEnvironmentVariables(ctx context.Context, environmentID string, envVars []sdk.EnvVariable) error
	ListServiceAccountTokens(ctx context.Context) ([]sdk.ServiceAccountToken, error)
	GetServiceAccountToken(ctx context.Context, tokenID string) (*sdk.ServiceAccountToken, error)
	CreateServiceAccountToken(ctx context.Context, req *sdk.ServiceAccountTokenRequest) (*sdk.ServiceAccountToken, error)
	UpdateServiceAccountToken(ctx context.Context, tokenID string, req *sdk.ServiceAccountTokenUpdateRequest) (*sdk.ServiceAccountToken, error)
	RevokeServiceAccountToken(ctx context.Context, tokenID string) error
	LoadConfig() error
	SaveConfig() error
	GetConfig() *config.Config
}
