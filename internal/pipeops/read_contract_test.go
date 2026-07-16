package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func TestGetEnvironmentUsesWorkspaceScopedListFallback(t *testing.T) {
	const workspaceUUID = "workspace-123"
	const environmentUUID = "environment-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/environment/fetch" {
			t.Fatalf("path = %q, want /environment/fetch", r.URL.Path)
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != workspaceUUID {
			t.Fatalf("workspace_uuid = %q, want %q", got, workspaceUUID)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"environments": []map[string]interface{}{
					{
						"uuid":         environmentUUID,
						"name":         "Production",
						"workspace_id": workspaceUUID,
					},
				},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, workspaceUUID)
	env, err := client.GetEnvironment(context.Background(), environmentUUID)
	if err != nil {
		t.Fatalf("GetEnvironment() error = %v", err)
	}
	if env.UUID != environmentUUID {
		t.Fatalf("env.UUID = %q, want %q", env.UUID, environmentUUID)
	}
}

func TestGetServiceAccountTokenIncludesWorkspaceUUID(t *testing.T) {
	const workspaceUUID = "workspace-123"
	const tokenUUID = "token-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/service-account-tokens/"+tokenUUID {
			t.Fatalf("path = %q, want /api/v1/service-account-tokens/%s", r.URL.Path, tokenUUID)
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != workspaceUUID {
			t.Fatalf("workspace_uuid = %q, want %q", got, workspaceUUID)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":           tokenUUID,
				"name":         "CI",
				"scopes":       []string{"api:read"},
				"token_prefix": "sat_test",
				"is_revoked":   false,
				"is_expired":   false,
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, workspaceUUID)
	token, err := client.GetServiceAccountToken(context.Background(), tokenUUID)
	if err != nil {
		t.Fatalf("GetServiceAccountToken() error = %v", err)
	}
	if token.UUID != tokenUUID {
		t.Fatalf("token.UUID = %q, want %q", token.UUID, tokenUUID)
	}
	if len(token.Permissions) != 1 || token.Permissions[0] != "api:read" {
		t.Fatalf("token.Permissions = %#v, want [api:read]", token.Permissions)
	}
	if !token.IsActive {
		t.Fatalf("token.IsActive = false, want true")
	}
}

func newTestClient(t *testing.T, baseURL, workspaceUUID string) *Client {
	t.Helper()

	sdkClient, err := sdk.NewClient(baseURL)
	if err != nil {
		t.Fatalf("sdk.NewClient() error = %v", err)
	}
	sdkClient.SetToken("sat_test")

	return &Client{
		sdkClient: sdkClient,
		config: &config.Config{
			OAuth: &config.OAuthConfig{
				BaseURL:     baseURL,
				AccessToken: "sat_test",
			},
			Settings: &config.Settings{
				DefaultWorkspaceUUID: workspaceUUID,
			},
		},
	}
}
