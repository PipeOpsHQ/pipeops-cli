package cmd

import (
	"testing"
)

func TestVolumesCommandsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "volumes" {
			found = true
			subcommands := map[string]bool{}
			for _, sub := range c.Commands() {
				subcommands[sub.Name()] = true
			}
			for _, name := range []string{"list", "get", "remount", "delete", "export", "export-status"} {
				if !subcommands[name] {
					t.Errorf("volumes missing subcommand %q", name)
				}
			}
			if c.Flags().Lookup("workspace") != nil {
				// workspace is on subcommands, not the group
			}
			listCmd, _, err := c.Find([]string{"list"})
			if err != nil {
				t.Fatalf("find list: %v", err)
			}
			if listCmd.Flag("workspace") == nil {
				t.Error("volumes list missing --workspace flag")
			}
			remountCmd, _, err := c.Find([]string{"remount"})
			if err != nil {
				t.Fatalf("find remount: %v", err)
			}
			for _, flag := range []string{"target-type", "target-uuid", "mount-path", "workspace"} {
				if remountCmd.Flag(flag) == nil {
					t.Errorf("volumes remount missing --%s flag", flag)
				}
			}
			deleteCmd, _, err := c.Find([]string{"delete"})
			if err != nil {
				t.Fatalf("find delete: %v", err)
			}
			if deleteCmd.Flag("yes") == nil {
				t.Error("volumes delete missing --yes flag")
			}
			break
		}
	}
	if !found {
		t.Error("volumes command not registered on root")
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
