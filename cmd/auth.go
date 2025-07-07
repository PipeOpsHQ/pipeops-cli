package cmd

import (
	"github.com/PipeOpsHQ/pipeops-cli/cmd/auth"
)

func init() {
	// Add the auth command as a subcommand of the root command
	rootCmd.AddCommand(auth.New())
}
