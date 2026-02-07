package agent

import "github.com/spf13/cobra"

type agentModel struct {
	rootCmd *cobra.Command
}

func NewAgent(rootCmd *cobra.Command) *agentModel {
	return &agentModel{
		rootCmd: rootCmd,
	}
}

func (a *agentModel) Register() {
	a.install()
	a.join()
	a.info()
	a.funnel()
	a.uninstall()
	a.registerUpdate()
	a.logs()
}

func (a *agentModel) registerUpdate() {
	updateCmd.Flags().String("cluster-name", "", "Name of the cluster to update")
	updateCmd.Flags().String("cluster-type", "", "Type of the cluster (k3s|minikube|k3d|kind|auto)")
	a.rootCmd.AddCommand(updateCmd)
}
