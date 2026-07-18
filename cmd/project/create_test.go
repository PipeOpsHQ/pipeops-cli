package project

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestProjectCreateRequestFromFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	addProjectCreateFlags(cmd)

	if err := cmd.Flags().Set("name", "api"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("server", "cluster-from-server"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("environment", "env-1"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("environment-name", "staging"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("repository", "acme/api"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("branch", "main"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("source", "github"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("build-command", "npm run build"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("start-command", "npm start"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("build-method", "nodejs"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("port", "8080"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("framework", "express"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("language", "nodejs"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("workspace", "ws-1"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("commit-url", "https://github.com/acme/api/commit/abc"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("commit-sha", "abc"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("env", "FOO=bar"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("env", "BAZ=qux"); err != nil {
		t.Fatal(err)
	}

	req, err := projectCreateRequestFromFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if req.Name != "api" {
		t.Fatalf("Name = %q", req.Name)
	}
	if req.ClusterUUID != "cluster-from-server" {
		t.Fatalf("ClusterUUID = %q", req.ClusterUUID)
	}
	if req.EnvironmentUUID != "env-1" {
		t.Fatalf("EnvironmentUUID = %q", req.EnvironmentUUID)
	}
	if req.Environment != "staging" {
		t.Fatalf("Environment = %q", req.Environment)
	}
	if req.Port != 8080 {
		t.Fatalf("Port = %d", req.Port)
	}
	if req.RepositoryLanguage != "nodejs" {
		t.Fatalf("RepositoryLanguage = %q", req.RepositoryLanguage)
	}
	if len(req.EnvVariables) != 2 {
		t.Fatalf("EnvVariables = %+v", req.EnvVariables)
	}
	if req.WorkspaceUUID != "ws-1" || req.CommitSha != "abc" {
		t.Fatalf("workspace/commit = %+v", req)
	}
}

func TestProjectCreateRequestFromFlags_ClusterAlias(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	addProjectCreateFlags(cmd)
	_ = cmd.Flags().Set("name", "api")
	_ = cmd.Flags().Set("cluster", "cluster-only")

	req, err := projectCreateRequestFromFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if req.ClusterUUID != "cluster-only" {
		t.Fatalf("ClusterUUID from --cluster = %q", req.ClusterUUID)
	}
}

func TestProjectCreateRequestFromFlags_ServerWinsOverCluster(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	addProjectCreateFlags(cmd)
	_ = cmd.Flags().Set("name", "api")
	_ = cmd.Flags().Set("server", "from-server")
	_ = cmd.Flags().Set("cluster", "from-cluster")

	req, err := projectCreateRequestFromFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if req.ClusterUUID != "from-server" {
		t.Fatalf("ClusterUUID = %q, want from-server", req.ClusterUUID)
	}
}

func TestProjectCreateRequestFromFlags_InvalidEnv(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	addProjectCreateFlags(cmd)
	_ = cmd.Flags().Set("name", "api")
	_ = cmd.Flags().Set("env", "not-a-pair")

	_, err := projectCreateRequestFromFlags(cmd)
	if err == nil {
		t.Fatal("expected error for invalid env pair")
	}
}

func TestProjectCreateRequestFromFlags_Worker(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	addProjectCreateFlags(cmd)
	_ = cmd.Flags().Set("name", "worker")
	_ = cmd.Flags().Set("worker", "true")

	req, err := projectCreateRequestFromFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !req.Worker {
		t.Fatal("Worker should be true")
	}
}
