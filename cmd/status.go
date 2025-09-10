package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [project-id]",
	Short: "📊 Show project or addon status",
	Long: `📊 Show the status of a project or addon.

View detailed information about your project's health, deployments, and services.
Can also show information about specific addons.

Examples:
  - Show status for linked project:
    pipeops status

  - Show status for specific project:
    pipeops status proj-123

  - Show addon information:
    pipeops status --addon redis

  - Show status in JSON format:
    pipeops status --json`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Parse flags
		addonID, _ := cmd.Flags().GetString("addon")

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

		if addonID != "" {
			// Show addon status
			showAddonStatus(client, addonID, opts)
		} else {
			// Show project status (existing behavior)
			showProjectStatus(client, args, opts)
		}
	},
	Args: cobra.MaximumNArgs(1),
}

func showAddonStatus(client *pipeops.Client, addonID string, opts utils.OutputOptions) {
	utils.PrintInfo(fmt.Sprintf("Getting addon '%s' information...", addonID), opts)

	addon, err := client.GetAddon(addonID)
	if err != nil {
		utils.HandleError(err, "Error fetching addon information", opts)
		return
	}

	if opts.Format == utils.OutputFormatJSON {
		utils.PrintJSON(addon)
	} else {
		// Display addon information
		utils.PrintInfo(fmt.Sprintf("📦 Addon: %s", addon.Name), opts)

		fmt.Printf("\n📊 ADDON DETAILS\n")
		fmt.Printf("├─ ID: %s\n", addon.ID)
		fmt.Printf("├─ Name: %s\n", addon.Name)
		fmt.Printf("├─ Category: %s\n", addon.Category)
		fmt.Printf("├─ Version: %s\n", addon.Version)
		fmt.Printf("├─ Status: %s %s\n", utils.GetStatusIcon(addon.Status), addon.Status)
		fmt.Printf("└─ Image: %s\n", addon.Image)

		if addon.Description != "" {
			fmt.Printf("\n📝 DESCRIPTION\n")
			fmt.Printf("%s\n", addon.Description)
		}

		if len(addon.Tags) > 0 {
			fmt.Printf("\n🏷️  TAGS\n")
			for i, tag := range addon.Tags {
				if i == len(addon.Tags)-1 {
					fmt.Printf("└─ %s\n", tag)
				} else {
					fmt.Printf("├─ %s\n", tag)
				}
			}
		}

		if len(addon.Ports) > 0 {
			fmt.Printf("\n🌐 PORTS\n")
			for i, port := range addon.Ports {
				if i == len(addon.Ports)-1 {
					fmt.Printf("└─ %d\n", port)
				} else {
					fmt.Printf("├─ %d\n", port)
				}
			}
		}

		if len(addon.EnvVars) > 0 {
			fmt.Printf("\n🔧 ENVIRONMENT VARIABLES\n")
			i := 0
			for key, value := range addon.EnvVars {
				if i == len(addon.EnvVars)-1 {
					fmt.Printf("└─ %s=%s\n", key, value)
				} else {
					fmt.Printf("├─ %s=%s\n", key, value)
				}
				i++
			}
		}

		fmt.Printf("\n⏰ TIMESTAMPS\n")
		fmt.Printf("├─ Created: %s\n", utils.FormatDate(addon.CreatedAt))
		fmt.Printf("└─ Updated: %s\n", utils.FormatDate(addon.UpdatedAt))

		// Show helpful tips
		if !opts.Quiet {
			fmt.Printf("\n💡 NEXT STEPS\n")
			fmt.Printf("├─ Deploy addon: pipeops deploy --addon %s --project <project-id>\n", addon.ID)
			fmt.Printf("├─ List all addons: pipeops list --addons\n")
			fmt.Printf("└─ View addon deployments: pipeops list --deployments --project <project-id>\n")
		}
	}
}

func showProjectStatus(client *pipeops.Client, args []string, opts utils.OutputOptions) {
	// Get project ID
	var projectID string
	var isLinkedProject bool

	if len(args) == 1 {
		projectID = args[0]
	} else {
		// Try to get from linked project
		projectContext, err := utils.LoadProjectContext()
		if err != nil || projectContext.ProjectID == "" {
			utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use 'pipeops link' to link a project or provide project ID as argument", opts)
			return
		}
		projectID = projectContext.ProjectID
		isLinkedProject = true
	}

	// Get project details
	utils.PrintInfo(fmt.Sprintf("Getting project '%s' status...", projectID), opts)

	project, err := client.GetProject(projectID)
	if err != nil {
		utils.HandleError(err, "Error fetching project", opts)
		return
	}

	// Get services for the project
	services, err := client.GetServices(projectID, "")
	if err != nil {
		// Services might not be available for all projects, don't fail
		services = &models.ListServicesResponse{Services: []models.ServiceInfo{}}
	}

	// Get addon deployments for the project
	addonDeployments, err := client.GetAddonDeployments(projectID)
	if err != nil {
		// Addon deployments might not be available, don't fail
		addonDeployments = []models.AddonDeployment{}
	}

	if opts.Format == utils.OutputFormatJSON {
		statusData := map[string]interface{}{
			"project":          project,
			"services":         services,
			"addon_deployments": addonDeployments,
			"is_linked":        isLinkedProject,
		}
		utils.PrintJSON(statusData)
	} else {
		// Display enhanced project information
		fmt.Printf("\n")
		if isLinkedProject {
			utils.PrintInfo(fmt.Sprintf("🔗 Linked Project: %s", project.Name), opts)
		} else {
			utils.PrintInfo(fmt.Sprintf("🚀 Project: %s", project.Name), opts)
		}

		// Project Overview
		fmt.Printf("\n📊 PROJECT OVERVIEW\n")
		fmt.Printf("├─ ID: %s\n", project.ID)
		fmt.Printf("├─ Name: %s\n", project.Name)
		fmt.Printf("├─ Status: %s %s\n", getStatusIcon(project.Status), project.Status)
		
		// Add description if available
		if project.Description != "" {
			fmt.Printf("├─ Description: %s\n", utils.TruncateString(project.Description, 60))
		}
		
		fmt.Printf("├─ Created: %s\n", utils.FormatDate(project.CreatedAt))
		fmt.Printf("└─ Last Updated: %s\n", utils.FormatDate(project.UpdatedAt))

		// Health Status Summary
		healthyServices := 0
		unhealthyServices := 0
		unknownServices := 0
		
		for _, service := range services.Services {
			switch strings.ToLower(service.Health) {
			case "healthy":
				healthyServices++
			case "unhealthy":
				unhealthyServices++
			default:
				unknownServices++
			}
		}
		
		if len(services.Services) > 0 {
			fmt.Printf("\n🏥 HEALTH STATUS\n")
			fmt.Printf("├─ Total Services: %d\n", len(services.Services))
			if healthyServices > 0 {
				fmt.Printf("├─ 🟢 Healthy: %d\n", healthyServices)
			}
			if unhealthyServices > 0 {
				fmt.Printf("├─ 🔴 Unhealthy: %d\n", unhealthyServices)
			}
			if unknownServices > 0 {
				fmt.Printf("└─ 🟡 Unknown: %d\n", unknownServices)
			}
		}

		// Show services with more details
		if len(services.Services) > 0 {
			fmt.Printf("\n🔧 SERVICES (%d)\n", len(services.Services))
			for i, service := range services.Services {
				symbol := "├─"
				if i == len(services.Services)-1 {
					symbol = "└─"
				}
				
				// Enhanced service display
				healthIcon := getHealthIcon(service.Health)
				fmt.Printf("%s %s %s\n", symbol, healthIcon, service.Name)
				
				// Add sub-details for each service
				subSymbol := "│  "
				if i == len(services.Services)-1 {
					subSymbol = "   "
				}
				
				fmt.Printf("%s ├─ Status: %s\n", subSymbol, service.Health)
				if service.Type != "" {
					fmt.Printf("%s ├─ Type: %s\n", subSymbol, service.Type)
				}
				if service.Protocol != "" {
					fmt.Printf("%s ├─ Protocol: %s\n", subSymbol, service.Protocol)
				}
				if service.Port != 0 {
					fmt.Printf("%s └─ Port: %d\n", subSymbol, service.Port)
				} else {
					fmt.Printf("%s └─ Port: N/A\n", subSymbol)
				}
			}
		}

		// Show addon deployments
		if len(addonDeployments) > 0 {
			fmt.Printf("\n📦 ADDON DEPLOYMENTS (%d)\n", len(addonDeployments))
			for i, addon := range addonDeployments {
				symbol := "├─"
				if i == len(addonDeployments)-1 {
					symbol = "└─"
				}
				
				statusIcon := utils.GetStatusIcon(addon.Status)
				fmt.Printf("%s %s %s\n", symbol, statusIcon, addon.Name)
				
				// Add sub-details for each addon
				subSymbol := "│  "
				if i == len(addonDeployments)-1 {
					subSymbol = "   "
				}
				
				fmt.Printf("%s ├─ ID: %s\n", subSymbol, addon.ID)
				fmt.Printf("%s ├─ Status: %s\n", subSymbol, addon.Status)
				if addon.URL != "" {
					fmt.Printf("%s ├─ URL: %s\n", subSymbol, addon.URL)
				}
				fmt.Printf("%s └─ Created: %s\n", subSymbol, utils.FormatDateShort(addon.CreatedAt))
			}
		}

		// Recent Activity
		fmt.Printf("\n📅 RECENT ACTIVITY\n")
		fmt.Printf("├─ Last deployment: %s\n", utils.FormatDate(project.UpdatedAt))
		fmt.Printf("└─ Project age: %s\n", getProjectAge(project.CreatedAt))

		// Show helpful tips based on project state
		if !opts.Quiet {
			fmt.Printf("\n💡 ACTIONS\n")
			
			// Context-aware actions
			if isLinkedProject {
				fmt.Printf("├─ Deploy changes: pipeops deploy\n")
				fmt.Printf("├─ View logs: pipeops logs\n")
				fmt.Printf("├─ Unlink project: pipeops unlink\n")
			} else {
				fmt.Printf("├─ Link to directory: pipeops link %s\n", projectID)
				fmt.Printf("├─ View logs: pipeops logs --project %s\n", projectID)
				fmt.Printf("├─ Deploy: pipeops deploy --project %s\n", projectID)
			}
			
			// Common actions
			if len(addonDeployments) == 0 {
				fmt.Printf("├─ Add addon: pipeops deploy --addon <addon-id> --project %s\n", projectID)
			} else {
				fmt.Printf("├─ Manage addons: pipeops list --deployments --project %s\n", projectID)
			}
			
			if len(services.Services) > 0 {
				fmt.Printf("├─ Connect to service: pipeops connect --project %s\n", projectID)
				fmt.Printf("├─ Execute command: pipeops exec --project %s\n", projectID)
			}
			
			fmt.Printf("└─ Open dashboard: https://app.pipeops.io/projects/%s\n", projectID)
		}
	}
}

// getProjectAge calculates and formats the age of a project
func getProjectAge(createdAt time.Time) string {
	duration := time.Since(createdAt)
	days := int(duration.Hours() / 24)
	
	if days == 0 {
		hours := int(duration.Hours())
		if hours == 0 {
			return "Less than an hour"
		}
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	
	if days == 1 {
		return "1 day"
	}
	
	if days < 30 {
		return fmt.Sprintf("%d days", days)
	}
	
	months := days / 30
	if months == 1 {
		return "1 month"
	}
	
	if months < 12 {
		return fmt.Sprintf("%d months", months)
	}
	
	years := months / 12
	if years == 1 {
		return "1 year"
	}
	
	return fmt.Sprintf("%d years", years)
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

	// Add flags
	statusCmd.Flags().StringP("addon", "a", "", "Show addon status instead of project status")
}
