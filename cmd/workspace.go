package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/workspace"
)

func init() {
	// Add the workspace command as a subcommand of the root command
	rootCmd.AddCommand(workspace.New())
}
