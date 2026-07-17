package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	mcpEndpoint        = "https://mcp.pipeops.app/mcp"
	mcpDocsURL         = "https://docs.pipeops.io/docs/integrations/pipeops-mcp"
	mcpServiceTokenURL = "https://console.pipeops.io/dashboard/integrations?cloudIntegrations=tokens"
)

type mcpSetupInfo struct {
	Endpoint        string `json:"endpoint"`
	Authentication  string `json:"authentication"`
	ServiceTokenURL string `json:"service_token_url"`
	DocsURL         string `json:"docs_url"`
	CodexCommand    string `json:"codex_command"`
}

func newMCPCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Connect AI assistants to PipeOps with MCP",
		Long: `Connect Codex and other MCP clients to your PipeOps account through the
hosted PipeOps MCP server.

Create a least-privilege service token in the PipeOps Console, store it in the
PIPEOPS_TOKEN environment variable, and configure your MCP client with the
hosted endpoint. This command never reads or prints your PipeOps login token.`,
		Example: `  export PIPEOPS_TOKEN="your-service-token"
  codex mcp add pipeops --url https://mcp.pipeops.app/mcp --bearer-token-env-var PIPEOPS_TOKEN
  codex mcp get pipeops`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := mcpSetupInfo{
				Endpoint:        mcpEndpoint,
				Authentication:  "Bearer service token via PIPEOPS_TOKEN",
				ServiceTokenURL: mcpServiceTokenURL,
				DocsURL:         mcpDocsURL,
				CodexCommand:    "codex mcp add pipeops --url " + mcpEndpoint + " --bearer-token-env-var PIPEOPS_TOKEN",
			}

			jsonOutput, _ := cmd.Flags().GetBool("json")
			if jsonOutput {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(info)
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), `PipeOps MCP

Hosted endpoint: %s
Authentication: Bearer service token

1. Create a service token with api:read (and api:write only when needed):
   %s

2. Store it in your environment:
   export PIPEOPS_TOKEN="your-service-token"

3. Connect Codex:
   %s

4. Verify the connection:
   codex mcp get pipeops

Documentation: %s
`, info.Endpoint, info.ServiceTokenURL, info.CodexCommand, info.DocsURL)
			return err
		},
	}
}

func init() {
	rootCmd.AddCommand(newMCPCommand())
}
