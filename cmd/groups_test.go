package cmd

import (
	"testing"
)

func TestGroupsCommandsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "groups" {
			found = true
			subcommands := map[string]bool{}
			for _, sub := range c.Commands() {
				subcommands[sub.Name()] = true
			}
			for _, name := range []string{
				"list", "get", "create", "update", "delete", "topology",
				"members", "env", "connect", "redeploy", "resolve", "candidates",
			} {
				if !subcommands[name] {
					t.Errorf("groups missing subcommand %q", name)
				}
			}

			membersCmd, _, err := c.Find([]string{"members"})
			if err != nil {
				t.Fatalf("find members: %v", err)
			}
			memberSubs := map[string]bool{}
			for _, sub := range membersCmd.Commands() {
				memberSubs[sub.Name()] = true
			}
			for _, name := range []string{"attach", "detach"} {
				if !memberSubs[name] {
					t.Errorf("groups members missing subcommand %q", name)
				}
			}

			envCmd, _, err := c.Find([]string{"env"})
			if err != nil {
				t.Fatalf("find env: %v", err)
			}
			envSubs := map[string]bool{}
			for _, sub := range envCmd.Commands() {
				envSubs[sub.Name()] = true
			}
			for _, name := range []string{"get", "put", "inject"} {
				if !envSubs[name] {
					t.Errorf("groups env missing subcommand %q", name)
				}
			}

			listCmd, _, err := c.Find([]string{"list"})
			if err != nil {
				t.Fatalf("find list: %v", err)
			}
			if listCmd.Flag("workspace") == nil {
				t.Error("groups list missing --workspace flag")
			}

			deleteCmd, _, err := c.Find([]string{"delete"})
			if err != nil {
				t.Fatalf("find delete: %v", err)
			}
			if deleteCmd.Flag("yes") == nil {
				t.Error("groups delete missing --yes flag")
			}

			attachCmd, _, err := c.Find([]string{"members", "attach"})
			if err != nil {
				t.Fatalf("find members attach: %v", err)
			}
			for _, flag := range []string{"type", "member-uuid", "workspace"} {
				if attachCmd.Flag(flag) == nil {
					t.Errorf("groups members attach missing --%s flag", flag)
				}
			}
			break
		}
	}
	if !found {
		t.Error("groups command not registered on root")
	}
}

func TestMapMemberType(t *testing.T) {
	if got := mapMemberType("addon"); got != "addon_deployment" {
		t.Errorf("addon = %q", got)
	}
	if got := mapMemberType("project"); got != "project" {
		t.Errorf("project = %q", got)
	}
	if got := mapMemberType("addon_deployment"); got != "addon_deployment" {
		t.Errorf("addon_deployment = %q", got)
	}
}

func TestDisplayOr(t *testing.T) {
	if got := displayOr("a", "b"); got != "a" {
		t.Errorf("displayOr primary = %q", got)
	}
	if got := displayOr("", "b"); got != "b" {
		t.Errorf("displayOr fallback = %q", got)
	}
}
