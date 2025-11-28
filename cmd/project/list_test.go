package project

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
)

func TestListProjects(t *testing.T) {
	// Save original factory and restore after test
	origFactory := pipeops.NewClientWithConfigFunc
	defer func() { pipeops.NewClientWithConfigFunc = origFactory }()

	// Setup mock client
	mockClient := &pipeops.MockClient{
		IsAuthenticatedFunc: func() bool {
			return true
		},
		GetProjectsFunc: func() (*models.ProjectsResponse, error) {
			return &models.ProjectsResponse{
				Projects: []models.Project{
					{
						ID:        "proj-123",
						Name:      "Test Project",
						Status:    "active",
						CreatedAt: time.Now(),
					},
				},
			}, nil
		},
	}

	// Replace factory
	pipeops.NewClientWithConfigFunc = func(cfg *config.Config) pipeops.ClientAPI {
		return mockClient
	}

	// Capture stdout
	output := captureOutput(func() {
		// Execute command
		listCmd.SetArgs([]string{})
		// We ignore error here because we want to check output even if it fails (though it shouldn't)
		_ = listCmd.Execute()
	})

	if !strings.Contains(output, "Test Project") {
		t.Errorf("Expected output to contain 'Test Project', got: %s", output)
	}
	if !strings.Contains(output, "proj-123") {
		t.Errorf("Expected output to contain 'proj-123', got: %s", output)
	}
}

func TestListProjectsEmpty(t *testing.T) {
	// Save original factory and restore after test
	origFactory := pipeops.NewClientWithConfigFunc
	defer func() { pipeops.NewClientWithConfigFunc = origFactory }()

	// Setup mock client
	mockClient := &pipeops.MockClient{
		IsAuthenticatedFunc: func() bool {
			return true
		},
		GetProjectsFunc: func() (*models.ProjectsResponse, error) {
			return &models.ProjectsResponse{
				Projects: []models.Project{},
			}, nil
		},
	}

	// Replace factory
	pipeops.NewClientWithConfigFunc = func(cfg *config.Config) pipeops.ClientAPI {
		return mockClient
	}

	// Capture stdout
	output := captureOutput(func() {
		// Execute command
		listCmd.SetArgs([]string{})
		_ = listCmd.Execute()
	})

	if !strings.Contains(output, "No projects found yet") {
		t.Errorf("Expected output to contain 'No projects found yet', got: %s", output)
	}
}

// captureOutput captures stdout from a function
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
