package agent

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAgentCommands(t *testing.T) {
	// Create a dummy root command
	rootCmd := &cobra.Command{Use: "pipeops"}
	
	// Initialize agent commands
	agentModel := NewAgent(rootCmd)
	agentModel.Register()

	// Verify that commands are registered
	commands := rootCmd.Commands()
	foundInstall := false
	foundUninstall := false
	foundUpdate := false

	for _, cmd := range commands {
		if cmd.Name() == "install" {
			foundInstall = true
			// Check flags
			if cmd.Flag("cluster-name") == nil {
				t.Error("install command missing 'cluster-name' flag")
			}
			if cmd.Flag("update") == nil {
				t.Error("install command missing 'update' flag")
			}
		}
		if cmd.Name() == "uninstall" {
			foundUninstall = true
			// Check aliases
			hasRemoveAlias := false
			hasRmAlias := false
			for _, alias := range cmd.Aliases {
				if alias == "remove" {
					hasRemoveAlias = true
				}
				if alias == "rm" {
					hasRmAlias = true
				}
			}
			if !hasRemoveAlias {
				t.Error("uninstall command missing 'remove' alias")
			}
			if !hasRmAlias {
				t.Error("uninstall command missing 'rm' alias")
			}
			// Check flags
			if cmd.Flag("force") == nil {
				t.Error("uninstall command missing 'force' flag")
			}
		}
		if cmd.Name() == "update" {
			foundUpdate = true
			// Check aliases
			hasUpgradeAlias := false
			for _, alias := range cmd.Aliases {
				if alias == "upgrade" {
					hasUpgradeAlias = true
					break
				}
			}
			if !hasUpgradeAlias {
				t.Error("update command missing 'upgrade' alias")
			}
			// Check flags
			if cmd.Flag("cluster-name") == nil {
				t.Error("update command missing 'cluster-name' flag")
			}
		}
		if cmd.Name() == "logs" {
			// Check flags
			if cmd.Flag("follow") == nil {
				t.Error("logs command missing 'follow' flag")
			}
			if cmd.Flag("tail") == nil {
				t.Error("logs command missing 'tail' flag")
			}
		}
	}

	if !foundInstall {
		t.Error("install command not registered")
	}
	if !foundUninstall {
		t.Error("uninstall command not registered")
	}
	if !foundUpdate {
		t.Error("update command not registered")
	}
}
