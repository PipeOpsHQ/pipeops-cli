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
	var err error

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
		utils.HandleError(err, "Error fetching services", opts)
		return
	}

	if opts.Format == utils.OutputFormatJSON {
		statusData := map[string]interface{}{
			"project":  project,
			"services": services,
		}
		utils.PrintJSON(statusData)
	} else {
		// Display project information
		utils.PrintInfo(fmt.Sprintf("🚀 Project: %s", project.Name), opts)

		fmt.Printf("\n📊 PROJECT STATUS\n")
		fmt.Printf("├─ ID: %s\n", project.ID)
		fmt.Printf("├─ Name: %s\n", project.Name)
		fmt.Printf("├─ Status: %s %s\n", utils.GetStatusIcon(project.Status), project.Status)
		fmt.Printf("├─ Created: %s\n", utils.FormatDate(project.CreatedAt))
		fmt.Printf("└─ Updated: %s\n", utils.FormatDate(project.UpdatedAt))

		// Show services
		if len(services.Services) > 0 {
			fmt.Printf("\n🔧 SERVICES (%d)\n", len(services.Services))
			for i, service := range services.Services {
				symbol := "├─"
				if i == len(services.Services)-1 {
					symbol = "└─"
				}
				fmt.Printf("%s %s %s (%s)\n", symbol, utils.GetStatusIcon(service.Health), service.Name, service.Health)
			}
		}

		// Show helpful tips
		if !opts.Quiet {
			fmt.Printf("\n💡 NEXT STEPS\n")
			fmt.Printf("├─ View logs: pipeops logs --project %s\n", projectID)
			fmt.Printf("├─ Deploy: pipeops deploy --project %s\n", projectID)
			fmt.Printf("├─ Connect: pipeops connect --project %s\n", projectID)
			fmt.Printf("└─ List deployments: pipeops list --deployments --project %s\n", projectID)
		}
	}
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
