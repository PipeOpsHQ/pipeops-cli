package addons

import (
	"testing"
)

func TestBackupsCommandsRegistered(t *testing.T) {
	found := false
	for _, c := range AddonsCmd.Commands() {
		if c.Name() == "backups" {
			found = true
			subcommands := map[string]bool{}
			for _, sub := range c.Commands() {
				subcommands[sub.Name()] = true
			}
			for _, name := range []string{"list", "export", "export-status"} {
				if !subcommands[name] {
					t.Errorf("addons backups missing subcommand %q", name)
				}
			}
			exportCmd, _, err := c.Find([]string{"export"})
			if err != nil {
				t.Fatalf("find export: %v", err)
			}
			for _, flag := range []string{"snapshot-id", "path", "format"} {
				if exportCmd.Flag(flag) == nil {
					t.Errorf("addons backups export missing --%s flag", flag)
				}
			}
			break
		}
	}
	if !found {
		t.Error("backups command not registered under addons")
	}
}

func TestFormatBytes(t *testing.T) {
	if got := formatBytes(512); got != "512 B" {
		t.Errorf("formatBytes(512) = %q", got)
	}
	if got := formatBytes(2048); got != "2.0 KB" {
		t.Errorf("formatBytes(2048) = %q", got)
	}
}
