// cmd/version.go
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "📦 Display the PipeOps CLI version",
	Long: `📦 The version command shows the current version of the PipeOps CLI.
This can be useful for debugging or verifying that you're using the expected version.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("🚀 PipeOps CLI Version: %s\n", GetConfig().Version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
