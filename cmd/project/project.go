package project

import "github.com/spf13/cobra"

type k3sModel struct {
	rootCmd *cobra.Command
}

func NewProject(rootCmd *cobra.Command) *k3sModel {
	return &k3sModel{
		rootCmd: rootCmd,
	}
}

func (k *k3sModel) Register() {
}
