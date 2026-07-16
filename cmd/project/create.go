package project

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project",
	Long: `Create a new PipeOps project.

Examples:
  pipeops project create --name api --server <server-id> --environment <env-id>
  pipeops project create --name api --repository owner/repo --branch main --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := authenticatedClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}

		req, err := projectCreateRequestFromFlags(cmd)
		if err != nil {
			return err
		}
		project, err := client.CreateProject(req)
		if err != nil {
			return fmt.Errorf("create project: %w", err)
		}

		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Project created", opts)
		}
		printProject(project, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

func (p *projectModel) createProject() *cobra.Command {
	addProjectCreateFlags(createCmd)
	p.rootCmd.AddCommand(createCmd)
	return createCmd
}
