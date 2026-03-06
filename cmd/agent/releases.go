package agent

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var releasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "List all Helm releases in the cluster",
	Long: `List all Helm releases across all namespaces in the cluster.
This is useful to see what applications and services are currently installed via Helm.

This command wraps 'helm list --all-namespaces'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for helm
		if _, err := exec.LookPath("helm"); err != nil {
			return fmt.Errorf("helm is required to list releases but was not found in PATH")
		}

		allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
		namespace, _ := cmd.Flags().GetString("namespace")

		var helmArgs []string
		helmArgs = append(helmArgs, "list")

		if allNamespaces {
			helmArgs = append(helmArgs, "--all-namespaces")
		} else if namespace != "" {
			helmArgs = append(helmArgs, "-n", namespace)
		}

		fmt.Printf("Fetching Helm releases...\n\n")

		execCmd := exec.Command("helm", helmArgs...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("failed to list helm releases: %w", err)
		}

		return nil
	},
}

func (a *agentModel) releases() {
	releasesCmd.Flags().BoolP("all-namespaces", "A", true, "List releases across all namespaces")
	releasesCmd.Flags().StringP("namespace", "n", "", "Namespace scope for this request")
	a.rootCmd.AddCommand(releasesCmd)
}
