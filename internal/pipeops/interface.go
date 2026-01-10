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
	GetLogs(req *models.LogsRequest) (*models.LogsResponse, error)
	StreamLogs(req *models.LogsRequest, callback func(*models.StreamLogEntry) error) error
	GetServices(projectID string, addonID string) (*models.ListServicesResponse, error)
	StartProxy(req *models.ProxyRequest) (*models.ProxyResponse, error)
	GetContainers(projectID string, addonID string) (*models.ListContainersResponse, error)
	StartExec(req *models.ExecRequest) (*models.ExecResponse, error)
	StartShell(req *models.ShellRequest) (*models.ShellResponse, error)
	GetAddons() (*models.AddonListResponse, error)
	GetAddon(addonID string) (*models.Addon, error)
	GetAddonDeployments(projectID string) ([]models.AddonDeployment, error)
	DeleteAddonDeployment(deploymentID string) error
	GetServers() (*models.ServersResponse, error)
	GetServer(serverID string) (*models.Server, error)
	CreateServer(req *models.ServerCreateRequest) (*models.Server, error)
	UpdateServer(serverID string, req *models.ServerUpdateRequest) (*models.Server, error)
	DeleteServer(serverID string) error
	VerifyToken() (*models.PipeOpsTokenVerificationResponse, error)
	GetWorkspaces(ctx context.Context) ([]sdk.Workspace, error)
	LoadConfig() error
	SaveConfig() error
	GetConfig() *config.Config
}
