// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the PipeOps CLI version",
	Long: `The version command shows the current version of the PipeOps CLI.
This can be useful for debugging or verifying that you're using the expected version.`,
	Run: func(cmd *cobra.Command, args []string) {
		// ASCII Art Logo
		logo := `
 ____  _             ___
|  _ \(_)_ __   ___ / _ \ _ __  ___
| |_) | | '_ \ / _ \ | | | '_ \/ __|
|  __/| | |_) |  __/ |_| | |_) \__ \
|_|   |_| .__/ \___|\___/| .__/|___/
        |_|              |_|
`
		fmt.Println(logo)
		fmt.Printf("PipeOps CLI Version: %s\n", Version)
		fmt.Printf("Documentation:       https://docs.pipeops.io\n")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
