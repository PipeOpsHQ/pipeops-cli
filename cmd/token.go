package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:     "token",
	Aliases: []string{"tokens", "service-token", "service-tokens"},
	Short:   "Manage service account tokens",
}

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List service account tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		tokens, err := client.ListServiceAccountTokens(context.Background())
		if err != nil {
			return fmt.Errorf("list service account tokens: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(tokens)
		}
		rows := make([][]string, 0, len(tokens))
		for _, token := range tokens {
			rows = append(rows, []string{token.UUID, token.Name, strings.Join(token.Permissions, ","), boolString(token.IsActive)})
		}
		utils.PrintTable([]string{"ID", "NAME", "PERMISSIONS", "ACTIVE"}, rows, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

var tokenGetCmd = &cobra.Command{
	Use:   "get <token-id>",
	Short: "Get service account token details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		token, err := client.GetServiceAccountToken(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("get service account token: %w", err)
		}
		return printToken(token, opts)
	},
	Args: cobra.ExactArgs(1),
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a service account token",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		permissions, _ := cmd.Flags().GetStringArray("permission")
		expiresAt, _ := cmd.Flags().GetString("expires-at")
		token, err := client.CreateServiceAccountToken(context.Background(), &sdk.ServiceAccountTokenRequest{
			Name:        name,
			Description: description,
			Permissions: permissions,
			ExpiresAt:   expiresAt,
		})
		if err != nil {
			return fmt.Errorf("create service account token: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Service account token created", opts)
		}
		return printToken(token, opts)
	},
	Args: cobra.NoArgs,
}

var tokenUpdateCmd = &cobra.Command{
	Use:   "update <token-id>",
	Short: "Update a service account token",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		permissions, _ := cmd.Flags().GetStringArray("permission")
		activeText, _ := cmd.Flags().GetString("active")
		active, err := parseOptionalBool(activeText)
		if err != nil {
			return fmt.Errorf("invalid --active value: %w", err)
		}
		token, err := client.UpdateServiceAccountToken(context.Background(), args[0], &sdk.ServiceAccountTokenUpdateRequest{
			Name:        name,
			Description: description,
			Permissions: permissions,
			IsActive:    active,
		})
		if err != nil {
			return fmt.Errorf("update service account token: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Service account token updated", opts)
		}
		return printToken(token, opts)
	},
	Args: cobra.ExactArgs(1),
}

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke <token-id>",
	Short: "Revoke a service account token",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			return fmt.Errorf("--force is required to revoke a service account token")
		}
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.RevokeServiceAccountToken(context.Background(), args[0]); err != nil {
			return fmt.Errorf("revoke service account token: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "revoked", "token_id": args[0]})
		}
		utils.PrintSuccess("Service account token revoked", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func printToken(token *sdk.ServiceAccountToken, opts utils.OutputOptions) error {
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(token)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"ID", token.UUID},
		{"Name", token.Name},
		{"Description", token.Description},
		{"Permissions", strings.Join(token.Permissions, ",")},
		{"Active", boolString(token.IsActive)},
		{"Token", token.Token},
	}, opts)
	return nil
}

func parseOptionalBool(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}
	switch strings.ToLower(value) {
	case "true", "t", "1", "yes", "y":
		parsed := true
		return &parsed, nil
	case "false", "f", "0", "no", "n":
		parsed := false
		return &parsed, nil
	default:
		return nil, fmt.Errorf("expected true or false")
	}
}

func boolString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func init() {
	tokenCreateCmd.Flags().String("name", "", "Token name")
	tokenCreateCmd.Flags().String("description", "", "Token description")
	tokenCreateCmd.Flags().StringArray("permission", nil, "Permission; repeatable")
	tokenCreateCmd.Flags().String("expires-at", "", "Expiration timestamp")
	_ = tokenCreateCmd.MarkFlagRequired("name")

	tokenUpdateCmd.Flags().String("name", "", "Token name")
	tokenUpdateCmd.Flags().String("description", "", "Token description")
	tokenUpdateCmd.Flags().StringArray("permission", nil, "Permission; repeatable")
	tokenUpdateCmd.Flags().String("active", "", "Set token active state true/false")
	tokenRevokeCmd.Flags().Bool("force", false, "Confirm token revocation")

	tokenCmd.AddCommand(tokenListCmd, tokenGetCmd, tokenCreateCmd, tokenUpdateCmd, tokenRevokeCmd)
	rootCmd.AddCommand(tokenCmd)
}
