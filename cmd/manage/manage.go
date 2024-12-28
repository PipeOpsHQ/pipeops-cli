package manage

import "github.com/spf13/cobra"

type manageModel struct {
	rootCmd *cobra.Command
}

func NewManagement(rootCmd *cobra.Command) *manageModel {
	return &manageModel{
		rootCmd: rootCmd,
	}
}

func (k *manageModel) Register() {
}
