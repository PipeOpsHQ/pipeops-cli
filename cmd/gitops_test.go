package cmd

import (
	"testing"
)

func TestGitOpsCommandsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "gitops" {
			found = true
			subcommands := map[string]bool{}
			for _, sub := range c.Commands() {
				subcommands[sub.Name()] = true
			}
			for _, name := range []string{"list", "get", "create", "update", "delete", "sync", "status", "diff", "history"} {
				if !subcommands[name] {
					t.Errorf("gitops missing subcommand %q", name)
				}
			}
			createCmd, _, err := c.Find([]string{"create"})
			if err != nil {
				t.Fatalf("find create: %v", err)
			}
			for _, flag := range []string{"name", "repo-url", "branch", "path", "project-id", "environment-id", "manifest-type"} {
				if createCmd.Flag(flag) == nil {
					t.Errorf("gitops create missing --%s flag", flag)
				}
			}
			deleteCmd, _, err := c.Find([]string{"delete"})
			if err != nil {
				t.Fatalf("find delete: %v", err)
			}
			if deleteCmd.Flag("yes") == nil {
				t.Error("gitops delete missing --yes flag")
			}
			syncCmd, _, err := c.Find([]string{"sync"})
			if err != nil {
				t.Fatalf("find sync: %v", err)
			}
			for _, flag := range []string{"revision", "prune", "dry-run"} {
				if syncCmd.Flag(flag) == nil {
					t.Errorf("gitops sync missing --%s flag", flag)
				}
			}
			break
		}
	}
	if !found {
		t.Error("gitops command not registered on root")
	}
}

func TestShortSHA(t *testing.T) {
	if got := shortSHA("abcdefghij"); got != "abcdefgh" {
		t.Errorf("shortSHA long = %q", got)
	}
	if got := shortSHA("abc"); got != "abc" {
		t.Errorf("shortSHA short = %q", got)
	}
}
