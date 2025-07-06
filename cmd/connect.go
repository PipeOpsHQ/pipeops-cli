package cmd

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [project-id] [service-name]",
	Short: "🔌 Connect to a database or service shell",
	Long: `🔌 Connect to a database or service shell. Automatically detects the service type
and launches the appropriate client (psql for PostgreSQL, mongosh for MongoDB, redis-cli for Redis, etc.).

If no arguments are provided, uses the linked project and prompts for service selection.

Examples:
  - Connect to a database in linked project (interactive):
    pipeops connect

  - Connect to specific service in linked project:
    pipeops connect postgres

  - Connect to specific service in specific project:
    pipeops connect proj-123 postgres

  - Connect to addon database:
    pipeops connect postgres --addon addon-456`,
	Run: func(cmd *cobra.Command, args []string) {
		var projectID, serviceName string
		var err error

		switch len(args) {
		case 0:
			// No arguments, use linked project and prompt for service
			projectID, err = utils.GetLinkedProject()
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				fmt.Println("💡 Use 'pipeops link <project-id>' to link a project to this directory")
				fmt.Println("   Or provide: pipeops connect <project-id> [service-name]")
				return
			}
			serviceName = "" // Will prompt for selection
		case 1:
			// Only service name provided, use linked project
			serviceName = args[0]
			projectID, err = utils.GetLinkedProject()
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				fmt.Println("💡 Use 'pipeops link <project-id>' to link a project to this directory")
				fmt.Println("   Or provide both: pipeops connect <project-id> <service-name>")
				return
			}
		case 2:
			// Both project ID and service name provided
			projectID = args[0]
			serviceName = args[1]
		default:
			fmt.Println("❌ Too many arguments")
			fmt.Println("Usage: pipeops connect [project-id] [service-name]")
			return
		}

		// Validate project ID
		if err := validation.ValidateProjectID(projectID); err != nil {
			fmt.Printf("❌ Invalid project ID: %v\n", err)
			return
		}

		client := pipeops.NewClient()

		// Load configuration
		if err := client.LoadConfig(); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		// Check if user is authenticated
		if !client.IsAuthenticated() {
			fmt.Println("❌ You are not logged in. Please run 'pipeops auth login' first.")
			return
		}

		// Show project context
		utils.PrintProjectContext(projectID)

		// Parse flags
		addonID, _ := cmd.Flags().GetString("addon")

		// If no service name provided, get available services and prompt
		if serviceName == "" {
			fmt.Printf("🔍 Fetching available services")
			if addonID != "" {
				fmt.Printf(" (addon: %s)", addonID)
			}
			fmt.Println("...")

			services, err := client.GetServices(projectID, addonID)
			if err != nil {
				fmt.Printf("❌ Error fetching services: %v\n", err)
				return
			}

			if len(services.Services) == 0 {
				fmt.Println("📭 No services found.")
				return
			}

			// Filter for database/connectable services
			var connectableServices []models.ServiceInfo
			for _, service := range services.Services {
				if isConnectableService(service.Type) {
					connectableServices = append(connectableServices, service)
				}
			}

			if len(connectableServices) == 0 {
				fmt.Println("📭 No connectable services found.")
				fmt.Println("💡 Connectable services include: database, postgres, mysql, mongodb, redis, etc.")
				return
			}

			// Display connectable services
			fmt.Println("\n🔌 Available services:")
			for i, service := range connectableServices {
				fmt.Printf("  %d. %s (%s) - %s\n", i+1, service.Name, service.Type, service.Description)
			}

			// Get user selection
			fmt.Print("\n🎯 Select a service (1-", len(connectableServices), "): ")
			var selection int
			if _, err := fmt.Scanln(&selection); err != nil || selection < 1 || selection > len(connectableServices) {
				fmt.Println("❌ Invalid selection.")
				return
			}

			serviceName = connectableServices[selection-1].Name
		}

		// Get service details
		fmt.Printf("🔍 Getting connection details for service '%s'...\n", serviceName)
		services, err := client.GetServices(projectID, addonID)
		if err != nil {
			fmt.Printf("❌ Error fetching service details: %v\n", err)
			return
		}

		// Find the specific service
		var targetService *models.ServiceInfo
		for _, service := range services.Services {
			if service.Name == serviceName {
				targetService = &service
				break
			}
		}

		if targetService == nil {
			fmt.Printf("❌ Service '%s' not found\n", serviceName)
			return
		}

		// Determine the connection command based on service type
		command := getConnectionCommand(targetService.Type)
		if len(command) == 0 {
			fmt.Printf("❌ Don't know how to connect to service type '%s'\n", targetService.Type)
			fmt.Println("💡 Supported types: postgres, postgresql, mysql, mongodb, redis, memcached")
			return
		}

		// Parse additional flags
		container, _ := cmd.Flags().GetString("container")
		user, _ := cmd.Flags().GetString("user")

		// Build exec request for interactive connection
		req := &models.ExecRequest{
			ProjectID:   projectID,
			AddonID:     addonID,
			ServiceName: serviceName,
			Container:   container,
			Command:     command,
			Interactive: true,
			User:        user,
		}

		// Start exec session
		fmt.Printf("🚀 Connecting to %s service '%s'...\n", targetService.Type, serviceName)

		execResp, err := client.StartExec(req)
		if err != nil {
			fmt.Printf("❌ Error starting connection: %v\n", err)
			return
		}

		fmt.Printf("💻 Connection session started (ID: %s)\n", execResp.ExecID)

		// Connect to terminal session
		session, err := terminalManager.StartExecSession(execResp.ExecID, execResp.WebSocketURL, true)
		if err != nil {
			fmt.Printf("❌ Error connecting to terminal: %v\n", err)
			return
		}
		defer session.Close()

		fmt.Printf("🔗 Connected to %s. Press Ctrl+C to exit.\n", targetService.Type)
		session.WaitForCompletion()
	},
	Args: cobra.MaximumNArgs(2),
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
	connectCmd.Flags().StringP("addon", "a", "", "Connect to service in a specific addon")
	connectCmd.Flags().StringP("container", "c", "", "Specific container name")
	connectCmd.Flags().StringP("user", "u", "", "User to run connection as")
}
