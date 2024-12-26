// cmd/version.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "v1.0.0" // Update this as needed

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the PipeOps CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PipeOps CLI Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
