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
	p.createProject()
	p.logs()
}
