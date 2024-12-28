package deploy

import "github.com/spf13/cobra"

type deployModel struct {
	rootCmd *cobra.Command
}

func NewDeploy(rootCmd *cobra.Command) *deployModel {
	return &deployModel{
		rootCmd: rootCmd,
	}
}

func (k *deployModel) Register() {
	k.newPipeline()
}
