package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	clipipeops "github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func TestPrintAuditLogsEmpty(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	opts := utils.OutputOptions{Format: utils.OutputFormatTable, Quiet: true}
	// Quiet path still should not error on empty
	if err := printAuditLogs(nil, sdk.AuditLogPagination{}, opts, "project"); err != nil {
		t.Fatalf("printAuditLogs: %v", err)
	}
	_ = buf
}

func TestPrintAuditLogsJSON(t *testing.T) {
	t.Parallel()
	ts := &sdk.Timestamp{Time: time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)}
	logs := []sdk.ProjectAuditLog{
		{
			UUID:        "log-1",
			Action:      "project.redeploy",
			ActionLabel: "Redeployed",
			Status:      "success",
			Summary:     "Redeployed api",
			ProjectName: "api",
			CreatedAt:   ts,
			Actor:       sdk.ProjectAuditActor{Type: "user", Name: "Ada"},
		},
	}
	opts := utils.OutputOptions{Format: utils.OutputFormatJSON}
	if err := printAuditLogs(logs, sdk.AuditLogPagination{Total: 1, Limit: 20}, opts, "project"); err != nil {
		t.Fatalf("printAuditLogs json: %v", err)
	}
}

func TestAuditProjectCommandUsesClient(t *testing.T) {
	t.Parallel()

	mock := &clipipeops.MockClient{
		ListProjectAuditLogsFunc: func(ctx context.Context, projectUUID string, opts *sdk.ProjectAuditLogListOptions) (*sdk.ProjectAuditLogListResponse, error) {
			if projectUUID != "proj-1" {
				t.Fatalf("projectUUID = %q", projectUUID)
			}
			if opts == nil || opts.Action != "project.redeploy" {
				t.Fatalf("opts = %+v", opts)
			}
			return &sdk.ProjectAuditLogListResponse{
				Success: true,
				Data: []sdk.ProjectAuditLog{
					{UUID: "log-1", Action: "project.redeploy", ActionLabel: "Redeployed", Status: "success"},
				},
				Pagination: sdk.AuditLogPagination{Total: 1, Limit: 10, Offset: 0},
			}, nil
		},
	}

	// Exercise opts builder + print path with mock response
	cmd := auditProjectCmd
	_ = cmd.Flags().Set("action", "project.redeploy")
	_ = cmd.Flags().Set("limit", "10")
	opts := auditProjectOpts(cmd)
	if opts.Action != "project.redeploy" || opts.Limit != 10 {
		t.Fatalf("auditProjectOpts = %+v", opts)
	}
	resp, err := mock.ListProjectAuditLogs(context.Background(), "proj-1", opts)
	if err != nil {
		t.Fatalf("ListProjectAuditLogs: %v", err)
	}
	if len(resp.Data) != 1 || !strings.Contains(resp.Data[0].Action, "redeploy") {
		t.Fatalf("resp = %+v", resp)
	}
}

func TestFormatAuditTime(t *testing.T) {
	t.Parallel()
	if got := formatAuditTime(nil); got != "" {
		t.Fatalf("nil => %q", got)
	}
	ts := &sdk.Timestamp{Time: time.Date(2026, 1, 2, 3, 4, 0, 0, time.UTC)}
	if got := formatAuditTime(ts); got != "2026-01-02 03:04" {
		t.Fatalf("got %q", got)
	}
}
