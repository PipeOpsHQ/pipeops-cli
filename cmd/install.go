package cmd

import (
	"github.com/spf13/cobra"
)

// installCmd represents the install command which aliases to agent install
var installCmd = &cobra.Command{
	Use:   "install [pipeops-token]",
	Short: "Alias for 'agent install'",
	Long:  `This command is an alias for 'pipeops agent install'. It installs the PipeOps agent on your Kubernetes cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Find the agent install command
		targetCmd, _, _ := rootCmd.Find([]string{"agent", "install"})
		
		if targetCmd != nil && targetCmd.Run != nil {
			// Run the agent install command using OUR command (which has the flags set)
			targetCmd.Run(cmd, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	
	// Copy flags from agent install command
	// We need to access the flags from agent package, but since we can't easily access the private installCmd there,
	// we'll manually add the common flags here to ensure they appear in help
	installCmd.Flags().String("cluster-name", "", "Name for the cluster (default: pipeops-cluster)")
	installCmd.Flags().String("cluster-type", "", "Kubernetes distribution (k3s|minikube|k3d|kind) (default: k3s)")
	installCmd.Flags().Bool("existing-cluster", false, "Install PipeOps agent on existing Kubernetes cluster")
	installCmd.Flags().Bool("no-monitoring", false, "Skip monitoring setup (agent only)")
	installCmd.Flags().Bool("update", false, "Update PipeOps agent to the latest version")
	installCmd.Flags().Bool("uninstall", false, "Uninstall PipeOps agent and related components")
}
