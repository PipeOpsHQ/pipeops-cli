package addons

import (
	"github.com/spf13/cobra"
)

// AddonsCmd represents the addons command group
var AddonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "Manage addons",
	Long: `Manage addons in your PipeOps account.

Addons are pre-built services like databases, caches, and message queues
that can be deployed alongside your projects.

Examples:
  - List all available addons:
    pipeops addons ls

  - View addon details:
    pipeops addons info <addon-id>

  - List addon deployments:
    pipeops addons deployments --project <project-id>`,
}

func init() {
	// Subcommands are added from their respective files
}
