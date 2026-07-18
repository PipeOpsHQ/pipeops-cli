package pipeops

import (
	"testing"

	"github.com/PipeOpsHQ/pipeops-cli/models"
)

func TestBuildSDKCreateProjectRequest_WebDefaults(t *testing.T) {
	req := &models.ProjectCreateRequest{
		Name:            "my-app",
		ClusterUUID:     "cluster-1",
		EnvironmentUUID: "env-1",
		Repository:      "https://github.com/acme/app",
		Branch:          "main",
		BuildCommand:    "npm run build",
		StartCommand:    "npm start",
		Framework:       "nextjs",
		EnvVariables: []models.ProjectEnvVar{
			{Key: "PORT", Value: "3000"},
		},
		WorkspaceUUID: "ws-1",
	}

	got := BuildSDKCreateProjectRequest(req)
	if got.Name != "my-app" {
		t.Fatalf("Name = %q", got.Name)
	}
	if got.ClusterUUID != "cluster-1" {
		t.Fatalf("ClusterUUID = %q", got.ClusterUUID)
	}
	if got.EnvironmentUUID != "env-1" {
		t.Fatalf("EnvironmentUUID = %q", got.EnvironmentUUID)
	}
	if got.Environment != "development" {
		t.Fatalf("Environment default = %q, want development", got.Environment)
	}
	if got.Source != "github" {
		t.Fatalf("Source default = %q, want github", got.Source)
	}
	if got.Username != "acme" {
		t.Fatalf("Username parsed = %q, want acme", got.Username)
	}
	if got.BuildSettings.BuildMethod != "nodejs" {
		t.Fatalf("BuildMethod = %q, want nodejs", got.BuildSettings.BuildMethod)
	}
	if got.BuildSettings.BuildCommand != "npm run build" {
		t.Fatalf("BuildCommand = %q", got.BuildSettings.BuildCommand)
	}
	if got.BuildSettings.RunCommand != "npm start" {
		t.Fatalf("RunCommand = %q", got.BuildSettings.RunCommand)
	}
	if got.BuildSettings.Worker == nil || *got.BuildSettings.Worker {
		t.Fatalf("Worker = %v, want false", got.BuildSettings.Worker)
	}
	if len(got.NetworkSettings) != 1 || got.NetworkSettings[0].Port != 3000 || got.NetworkSettings[0].Protocol != "HTTP" {
		t.Fatalf("NetworkSettings = %+v, want port 3000 HTTP", got.NetworkSettings)
	}
	if got.EnvVariables == nil {
		t.Fatal("EnvVariables must not be nil")
	}
	if len(got.EnvVariables) != 1 || got.EnvVariables[0].Key != "PORT" {
		t.Fatalf("EnvVariables = %+v", got.EnvVariables)
	}
	if got.JobDetails.Enable == nil || *got.JobDetails.Enable {
		t.Fatalf("JobDetails.Enable = %v, want false", got.JobDetails.Enable)
	}
	if got.JobDetails.Suspended == nil || *got.JobDetails.Suspended {
		t.Fatalf("JobDetails.Suspended = %v, want false", got.JobDetails.Suspended)
	}
	if got.Framework != "nextjs" {
		t.Fatalf("Framework = %q", got.Framework)
	}
	if got.RepositoryLanguage != "nextjs" {
		t.Fatalf("RepositoryLanguage fallback = %q, want nextjs", got.RepositoryLanguage)
	}
}

func TestBuildSDKCreateProjectRequest_Worker(t *testing.T) {
	got := BuildSDKCreateProjectRequest(&models.ProjectCreateRequest{
		Name:         "worker",
		ClusterUUID:  "c1",
		Repository:   "acme/jobs",
		Worker:       true,
		StartCommand: "node worker.js",
		BuildMethod:  "nodejs",
	})
	if got.BuildSettings.Worker == nil || !*got.BuildSettings.Worker {
		t.Fatalf("Worker = %v, want true", got.BuildSettings.Worker)
	}
	if len(got.NetworkSettings) != 0 {
		t.Fatalf("worker NetworkSettings should be empty, got %+v", got.NetworkSettings)
	}
	if got.Username != "acme" {
		t.Fatalf("Username = %q, want acme", got.Username)
	}
	if got.EnvVariables == nil {
		t.Fatal("EnvVariables must be empty slice, not nil")
	}
	if len(got.EnvVariables) != 0 {
		t.Fatalf("EnvVariables len = %d", len(got.EnvVariables))
	}
}

func TestBuildSDKCreateProjectRequest_DockerfileAndLanguage(t *testing.T) {
	got := BuildSDKCreateProjectRequest(&models.ProjectCreateRequest{
		Name:               "api",
		ClusterUUID:        "c1",
		RepositoryLanguage: "docker",
		Port:               8080,
		Environment:        "staging",
		Source:             "gitlab",
		Username:           "custom",
	})
	if got.BuildSettings.BuildMethod != "dockerfile" {
		t.Fatalf("BuildMethod = %q, want dockerfile", got.BuildSettings.BuildMethod)
	}
	if len(got.NetworkSettings) != 1 || got.NetworkSettings[0].Port != 8080 {
		t.Fatalf("NetworkSettings = %+v", got.NetworkSettings)
	}
	if got.Environment != "staging" {
		t.Fatalf("Environment = %q", got.Environment)
	}
	if got.Source != "gitlab" {
		t.Fatalf("Source = %q", got.Source)
	}
	if got.Username != "custom" {
		t.Fatalf("Username = %q", got.Username)
	}
}

func TestBuildSDKCreateProjectRequest_LegacyEnvVarsMap(t *testing.T) {
	got := BuildSDKCreateProjectRequest(&models.ProjectCreateRequest{
		Name:        "legacy",
		ClusterUUID: "c1",
		EnvVars: map[string]interface{}{
			"FOO": "bar",
		},
	})
	if len(got.EnvVariables) != 1 || got.EnvVariables[0].Key != "FOO" || got.EnvVariables[0].Value != "bar" {
		t.Fatalf("EnvVariables from map = %+v", got.EnvVariables)
	}
}

func TestUsernameFromRepository(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"acme/app", "acme"},
		{"https://github.com/acme/app", "acme"},
		{"https://github.com/acme/app.git", "acme"},
		{"git@github.com:acme/app.git", "acme"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := usernameFromRepository(tc.in); got != tc.want {
			t.Errorf("usernameFromRepository(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
