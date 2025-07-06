package project

import "github.com/spf13/cobra"

type projectModel struct {
	rootCmd *cobra.Command
}

func NewProject(rootCmd *cobra.Command) *projectModel {
	return &projectModel{
		rootCmd: rootCmd,
	}
}

func (p *projectModel) Register() {
	p.listProjects()
	p.createProject() // Re-enabled to show disabled message
	p.logs()
}
