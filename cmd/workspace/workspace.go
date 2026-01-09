package workspace

import (
	"github.com/spf13/cobra"
)

// workspaceModel represents the workspace command model
type workspaceModel struct {
	rootCmd *cobra.Command
}

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage your PipeOps workspaces",
	Long: `Manage your PipeOps workspaces.

Common commands:
  pipeops workspace list      List all workspaces
  pipeops workspace select    Select a default workspace

Get started by running: pipeops workspace list`,
}

// New initializes and returns workspace command
func New() *cobra.Command {
	wm := &workspaceModel{
		rootCmd: workspaceCmd,
	}

	wm.list()
	wm.selectWorkspace()

	return workspaceCmd
}
