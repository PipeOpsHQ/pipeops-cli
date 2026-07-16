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
  - List all deployable addons:
    pipeops addons
    pipeops addons available

  - View addon details:
    pipeops addons info <addon-id>

  - List deployed addons:
    pipeops addons list`,
}

func init() {
	// Subcommands are added from their respective files
}
