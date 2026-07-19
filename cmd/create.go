package cmd

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/cmd/project"
	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a PipeOps project",
	Long: `Create a new PipeOps project.

This is a convenience alias for project creation:
  pipeops create --name api --server <cluster-uuid> --environment <env-uuid>
  pipeops create --name api --repository owner/repo --branch main --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load configuration: %w", err)
		}
		client := pipeops.NewClientWithConfigFunc(cfg)
		if !utils.RequireAuth(client, opts) {
			return nil
		}
		req, err := project.ProjectCreateRequestFromFlags(cmd)
		if err != nil {
			return err
		}
		created, err := client.CreateProject(req)
		if err != nil {
			return fmt.Errorf("create project: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(created)
		}
		utils.PrintSuccess("Project created", opts)
		utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
			{"ID", created.ID},
			{"Name", created.Name},
			{"Status", created.Status},
			{"Description", created.Description},
		}, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	project.AddProjectCreateFlags(createCmd)
	rootCmd.AddCommand(createCmd)
}
