package cmd

import (
	"fmt"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [project-id]",
	Short: "📊 Show project status and information",
	Long: `📊 Display detailed information about the current or specified project,
including services, deployment status, and recent activity.

If no project ID is provided, uses the linked project from the current directory.

Examples:
  - Show status for linked project:
    pipeops status

  - Show status for specific project:
    pipeops status proj-123`,
	Run: func(cmd *cobra.Command, args []string) {
		var projectID string
		var err error

		if len(args) == 1 {
			projectID = args[0]
		} else {
			// Try to get linked project
			projectID, err = utils.GetLinkedProject()
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				fmt.Println("💡 Use 'pipeops link <project-id>' to link a project to this directory")
				fmt.Println("   Or provide: pipeops status <project-id>")
				return
			}
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

		// Get project details
		fmt.Printf("🔍 Fetching project information...\n\n")

		project, err := client.GetProject(projectID)
		if err != nil {
			fmt.Printf("❌ Error fetching project: %v\n", err)
			return
		}

		// Display project information
		fmt.Printf("📂 PROJECT INFORMATION\n")
		fmt.Printf("├─ Name: %s\n", project.Name)
		fmt.Printf("├─ ID: %s\n", project.ID)
		fmt.Printf("├─ Status: %s\n", getStatusIcon(project.Status)+project.Status)
		fmt.Printf("├─ Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
		if project.Description != "" {
			fmt.Printf("└─ Description: %s\n", project.Description)
		} else {
			fmt.Printf("└─ Description: (none)\n")
		}

		// Show if this is the linked project
		if utils.IsLinkedProject() {
			if linkedID, err := utils.GetLinkedProject(); err == nil && linkedID == projectID {
				fmt.Printf("\n🔗 This project is linked to the current directory\n")
			}
		}

		// Get and display services
		fmt.Printf("\n🚀 SERVICES\n")
		services, err := client.GetServices(projectID, "")
		if err != nil {
			fmt.Printf("❌ Error fetching services: %v\n", err)
		} else if len(services.Services) == 0 {
			fmt.Printf("└─ No services found\n")
		} else {
			for i, service := range services.Services {
				isLast := i == len(services.Services)-1
				prefix := "├─"
				if isLast {
					prefix = "└─"
				}

				healthIcon := getHealthIcon(service.Health)
				fmt.Printf("%s %s %s (%s) - %s:%d\n",
					prefix, healthIcon, service.Name, service.Type, service.Protocol, service.Port)
			}
		}

		// Get proxy information
		fmt.Printf("\n🌐 ACTIVE PROXIES\n")
		proxies := proxyManager.ListProxies()
		if len(proxies.Proxies) == 0 {
			fmt.Printf("└─ No active proxies\n")
		} else {
			activeCount := 0
			for _, proxy := range proxies.Proxies {
				if proxy.Status == "active" {
					activeCount++
				}
			}
			fmt.Printf("└─ %d active proxy connections\n", activeCount)
		}

		fmt.Printf("\n✨ QUICK ACTIONS\n")
		fmt.Printf("├─ View logs: pipeops logs\n")
		fmt.Printf("├─ Start shell: pipeops shell <service-name>\n")
		fmt.Printf("├─ Connect to DB: pipeops connect\n")
		fmt.Printf("└─ Start proxy: pipeops proxy start <service-name>\n")
	},
	Args: cobra.MaximumNArgs(1),
}

// getStatusIcon returns an icon for project status
func getStatusIcon(status string) string {
	switch strings.ToLower(status) {
	case "active", "running", "healthy":
		return "🟢 "
	case "deploying", "building", "starting":
		return "🟡 "
	case "stopped", "inactive":
		return "⚪ "
	case "error", "failed", "crashed":
		return "🔴 "
	default:
		return "⚫ "
	}
}

// getHealthIcon returns an icon for service health
func getHealthIcon(health string) string {
	switch strings.ToLower(health) {
	case "healthy":
		return "🟢"
	case "unhealthy":
		return "🔴"
	case "unknown":
		return "🟡"
	default:
		return "⚫"
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
