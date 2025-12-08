package cmd

import (
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [service-name]",
	Short: "Connect to a service",
	Long: `Connect to a service in your project.

This command helps you connect to various services like databases, caches, and other infrastructure components.

Examples:
  - Connect to a database:
    pipeops connect postgres --project proj-123

  - Connect to a service by name:
    pipeops connect web-service --project proj-123`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)
		utils.PrintWarning("The 'connect' command is coming soon! Please check our documentation for updates.", opts)
		return
	},
	Args: cobra.MaximumNArgs(1),
}

// isConnectableService determines if a service type can be connected to
func isConnectableService(serviceType string) bool {
	connectableTypes := []string{
		"database", "postgres", "postgresql", "mysql", "mariadb",
		"mongodb", "mongo", "redis", "memcached", "cassandra",
		"elasticsearch", "clickhouse", "influxdb",
	}

	serviceType = strings.ToLower(serviceType)
	for _, connectable := range connectableTypes {
		if serviceType == connectable {
			return true
		}
	}
	return false
}

// getConnectionCommand returns the appropriate command to connect to a service type
func getConnectionCommand(serviceType string) []string {
	serviceType = strings.ToLower(serviceType)

	switch serviceType {
	case "postgres", "postgresql", "database":
		return []string{"psql"}
	case "mysql", "mariadb":
		return []string{"mysql"}
	case "mongodb", "mongo":
		return []string{"mongosh"}
	case "redis":
		return []string{"redis-cli"}
	case "memcached":
		return []string{"telnet", "localhost", "11211"}
	case "cassandra":
		return []string{"cqlsh"}
	case "elasticsearch":
		return []string{"curl", "-X", "GET", "localhost:9200"}
	case "clickhouse":
		return []string{"clickhouse-client"}
	case "influxdb":
		return []string{"influx"}
	default:
		return nil
	}
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Add flags
	connectCmd.Flags().StringP("project", "p", "", "Project ID")
}
