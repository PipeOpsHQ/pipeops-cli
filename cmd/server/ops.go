package server

import (
	"fmt"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

func serverClient(cmd *cobra.Command, opts utils.OutputOptions) (pipeops.ClientAPI, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}
	client := pipeops.NewClientWithConfig(cfg)
	if !utils.RequireAuth(client, opts) {
		return nil, nil
	}
	return client, nil
}

func GetConnectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connection <server-id>",
		Short: "Get server connection information",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := utils.GetOutputOptions(cmd)
			client, err := serverClient(cmd, opts)
			if err != nil || client == nil {
				return err
			}
			connection, err := client.GetServerConnection(args[0])
			if err != nil {
				return fmt.Errorf("get server connection: %w", err)
			}
			return utils.PrintJSON(connection)
		},
		Args: cobra.ExactArgs(1),
	}
}

func GetCostCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cost <server-id>",
		Short: "Get server cost allocation",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := utils.GetOutputOptions(cmd)
			client, err := serverClient(cmd, opts)
			if err != nil || client == nil {
				return err
			}
			costs, err := client.GetServerCostAllocation(args[0])
			if err != nil {
				return fmt.Errorf("get server cost allocation: %w", err)
			}
			return utils.PrintJSON(costs)
		},
		Args: cobra.ExactArgs(1),
	}
}
