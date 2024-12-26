package cmd

import "github.com/spf13/cobra"

type k3sModel struct {
	rootCmd *cobra.Command
}

func NewK3s(rootCmd *cobra.Command) *k3sModel {
	return &k3sModel{
		rootCmd: rootCmd,
	}
}

func (k *k3sModel) Register() {
	k.install()
	k.join()
	k.kill()
	k.restart()
}
