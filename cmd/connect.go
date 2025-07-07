package cmd

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [service-name]",
	Short: "🔗 Connect to a service",
	Long: `🔗 Connect to a service in your project.

This command helps you connect to various services like databases, caches, and other infrastructure components.

Examples:
  - Connect to a database:
    pipeops connect postgres --project proj-123

  - Connect to a service by name:
    pipeops connect web-service --project proj-123`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Get project context
		projectContext, err := utils.LoadProjectContext()
		if err != nil {
			utils.HandleError(err, "Error loading project context", opts)
			return
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			utils.HandleError(err, "Error loading configuration", opts)
			return
		}

		// Check if user is authenticated
		if !utils.RequireAuth(client, opts) {
			return
		}

		// Get project ID
		projectID := projectContext.ProjectID
		if projectID == "" {
			flagProjectID, _ := cmd.Flags().GetString("project")
			if flagProjectID == "" {
				utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use --project flag or link a project with 'pipeops link'", opts)
				return
			}
			projectID = flagProjectID
		}

		// Get service name
		var serviceName string
		if len(args) > 0 {
			serviceName = args[0]
		}

		// If no service name provided, list available services
		if serviceName == "" {
			utils.PrintInfo("Fetching available services...", opts)

			services, err := client.GetServices(projectID, "")
			if err != nil {
				utils.HandleError(err, "Error fetching services", opts)
				return
			}

			if len(services.Services) == 0 {
				utils.PrintWarning("No services found for this project", opts)
				return
			}

			utils.PrintInfo("Available services:", opts)
			for _, service := range services.Services {
				fmt.Printf("  - %s (%s)\n", service.Name, service.Type)
			}
			utils.PrintInfo("Use: pipeops connect <service-name>", opts)
			return
		}

		// Get service information
		utils.PrintInfo(fmt.Sprintf("Connecting to service '%s'...", serviceName), opts)

		services, err := client.GetServices(projectID, "")
		if err != nil {
			utils.HandleError(err, "Error fetching services", opts)
			return
		}

		// Find the service
		var targetService *models.ServiceInfo
		for _, service := range services.Services {
			if service.Name == serviceName {
				targetService = &service
				break
			}
		}

		if targetService == nil {
			utils.HandleError(fmt.Errorf("service not found"), fmt.Sprintf("Service '%s' not found in project", serviceName), opts)
			return
		}

		// Start proxy for the service
		req := &models.ProxyRequest{
			Target: models.ProxyTarget{
				ProjectID:   projectID,
				ServiceName: serviceName,
				Port:        targetService.Port,
			},
			LocalPort: 0, // Auto-assign
		}

		proxyResp, err := client.StartProxy(req)
		if err != nil {
			utils.HandleError(err, "Error starting connection", opts)
			return
		}

		if opts.Format == utils.OutputFormatJSON {
			utils.PrintJSON(proxyResp)
		} else {
			utils.PrintSuccess(fmt.Sprintf("Connected to %s service", serviceName), opts)
			utils.PrintInfo(fmt.Sprintf("Local connection: localhost:%d", proxyResp.LocalPort), opts)
			utils.PrintInfo(fmt.Sprintf("Remote endpoint: %s:%d", proxyResp.RemoteHost, proxyResp.RemotePort), opts)
			utils.PrintInfo("Connection is active. Press Ctrl+C to disconnect", opts)
		}
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
