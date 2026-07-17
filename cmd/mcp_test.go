package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestMCPCommandTextOutput(t *testing.T) {
	cmd := newMCPCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute mcp command: %v", err)
	}

	for _, expected := range []string{
		"PipeOps MCP",
		mcpEndpoint,
		"api:read",
		"--bearer-token-env-var PIPEOPS_TOKEN",
		mcpServiceTokenURL,
		mcpDocsURL,
	} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("output does not contain %q\n%s", expected, output.String())
		}
	}
}

func TestMCPCommandJSONOutput(t *testing.T) {
	cmd := newMCPCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.Flags().Bool("json", false, "Output in JSON format")
	if err := cmd.Flags().Set("json", "true"); err != nil {
		t.Fatalf("set json flag: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute mcp command: %v", err)
	}

	var info mcpSetupInfo
	if err := json.Unmarshal(output.Bytes(), &info); err != nil {
		t.Fatalf("decode JSON output: %v\n%s", err, output.String())
	}
	if info.Endpoint != mcpEndpoint {
		t.Errorf("endpoint = %q, want %q", info.Endpoint, mcpEndpoint)
	}
	if info.ServiceTokenURL != mcpServiceTokenURL {
		t.Errorf("service token URL = %q, want %q", info.ServiceTokenURL, mcpServiceTokenURL)
	}
	if !strings.Contains(info.CodexCommand, "PIPEOPS_TOKEN") {
		t.Errorf("Codex command does not reference PIPEOPS_TOKEN: %q", info.CodexCommand)
	}
}

func TestRootCommandIncludesMCP(t *testing.T) {
	command, _, err := rootCmd.Find([]string{"mcp"})
	if err != nil {
		t.Fatalf("find mcp command: %v", err)
	}
	if command.Name() != "mcp" {
		t.Fatalf("root command resolved %q, want mcp", command.Name())
	}
}
