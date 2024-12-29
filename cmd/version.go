// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "📦 Display the PipeOps CLI version",
	Long: `📦 The version command shows the current version of the PipeOps CLI.
This can be useful for debugging or verifying that you're using the expected version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🚀 PipeOps CLI Version: %s\n", GetConfig().Version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
