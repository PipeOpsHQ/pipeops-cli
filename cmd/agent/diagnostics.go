package agent

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var diagnosticsCmd = &cobra.Command{
	Use:     "diagnostics",
	Aliases: []string{"diag", "report"},
	Short:   "Run diagnostic checks for the PipeOps agent",
	Long: `Run a series of diagnostic checks and print a unified report about your PipeOps agent.

This command is particularly useful when troubleshooting issues or when requested by PipeOps Support.
It collects:
- Agent pod status
- Recent agent events
- Node readiness
- Basic cluster info
- Helm releases`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for kubectl
		if _, err := exec.LookPath("kubectl"); err != nil {
			return fmt.Errorf("kubectl is required for diagnostics but was not found in PATH")
		}

		fmt.Println("========================================")
		fmt.Println("      PipeOps Agent Diagnostics         ")
		fmt.Println("========================================")
		fmt.Println()

		runDiagnosticCommand("1. Agent Pod Status", "kubectl", "get", "pods", "-n", "pipeops-system", "-l", "app=pipeops-agent")
		runDiagnosticCommand("2. Recent Agent Events", "kubectl", "get", "events", "-n", "pipeops-system", "--sort-by=.metadata.creationTimestamp")
		runDiagnosticCommand("3. Node Status", "kubectl", "get", "nodes")

		if _, err := exec.LookPath("helm"); err == nil {
			runDiagnosticCommand("4. Helm Releases (pipeops-system)", "helm", "list", "-n", "pipeops-system")
		} else {
			fmt.Println("4. Helm Releases (pipeops-system)")
			fmt.Println("   [SKIPPED - helm binary not found]")
			fmt.Println()
		}

		fmt.Println("========================================")
		fmt.Println("      End of Diagnostics Report         ")
		fmt.Println("========================================")
		fmt.Println()

		return nil
	},
}

func runDiagnosticCommand(title string, cmdName string, cmdArgs ...string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("%s\n", strings.Repeat("-", len(title)))

	execCmd := exec.Command(cmdName, cmdArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
	fmt.Println()
}

func (a *agentModel) diagnostics() {
	a.rootCmd.AddCommand(diagnosticsCmd)
}
