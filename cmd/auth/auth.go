package auth

import "github.com/spf13/cobra"

type authModel struct {
	rootCmd *cobra.Command
}

func NewAuth(rootCmd *cobra.Command) *authModel {
	return &authModel{
		rootCmd: rootCmd,
	}
}

func (k *authModel) Register() {
	k.me()
}
