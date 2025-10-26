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
}
