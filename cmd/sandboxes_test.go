package cmd

import (
	"testing"
)

func TestSandboxesCommandsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "sandboxes" {
			found = true
			subcommands := map[string]bool{}
			for _, sub := range c.Commands() {
				subcommands[sub.Name()] = true
			}
			for _, name := range []string{"list", "get", "create", "start", "stop", "restart", "delete", "session", "usage"} {
				if !subcommands[name] {
					t.Errorf("sandboxes missing subcommand %q", name)
				}
			}
			listCmd, _, err := c.Find([]string{"list"})
			if err != nil {
				t.Fatalf("find list: %v", err)
			}
			if listCmd.Flag("workspace") == nil {
				t.Error("sandboxes list missing --workspace flag")
			}
			deleteCmd, _, err := c.Find([]string{"delete"})
			if err != nil {
				t.Fatalf("find delete: %v", err)
			}
			if deleteCmd.Flag("yes") == nil {
				t.Error("sandboxes delete missing --yes flag")
			}
			break
		}
	}
	if !found {
		t.Error("sandboxes command not registered on root")
	}
}
