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

var logsCmd = &cobra.Command{
	Use:   "logs [project-id]",
	Short: "📋 View project logs",
	Long: `📋 View and stream logs from your project services.

This command allows you to view historical logs and stream real-time logs from your deployed services.

Examples:
  - View logs for linked project:
    pipeops logs

  - View logs for specific project:
    pipeops logs proj-123

  - Stream logs in real-time:
    pipeops logs --follow

  - View last 100 lines:
    pipeops logs --lines 100

  - Filter logs by service:
    pipeops logs --service web-service`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := utils.GetOutputOptions(cmd)

		// Get project ID
		var projectID string
		if len(args) == 1 {
			projectID = args[0]
		} else {
			projectContext, err := utils.LoadProjectContext()
			if err != nil {
				utils.HandleError(err, "Error loading project context", opts)
				return
			}

			projectID = projectContext.ProjectID
			if projectID == "" {
				utils.HandleError(fmt.Errorf("project ID is required"), "Project ID is required. Use --project flag or link a project with 'pipeops link'", opts)
				return
			}
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

		// Parse flags
		serviceName, _ := cmd.Flags().GetString("service")
		containerName, _ := cmd.Flags().GetString("container")
		follow, _ := cmd.Flags().GetBool("follow")
		lines, _ := cmd.Flags().GetInt("lines")
		sinceStr, _ := cmd.Flags().GetString("since")
		untilStr, _ := cmd.Flags().GetString("until")

		// Build logs request
		req := &models.LogsRequest{
			ProjectID: projectID,
			Source:    serviceName,
			Container: containerName,
			Tail:      lines,
			Follow:    follow,
		}

		// Parse time filters
		if sinceStr != "" {
			since, err := time.Parse(time.RFC3339, sinceStr)
			if err != nil {
				utils.HandleError(fmt.Errorf("invalid since time format. Use RFC3339 format: %v", err), "Invalid time format", opts)
				return
			}
			req.Since = &since
		}

		if untilStr != "" {
			until, err := time.Parse(time.RFC3339, untilStr)
			if err != nil {
				utils.HandleError(fmt.Errorf("invalid until time format. Use RFC3339 format: %v", err), "Invalid time format", opts)
				return
			}
			req.Until = &until
		}

		if follow {
			// Stream logs in real-time
			utils.PrintInfo("Starting log stream... (Press Ctrl+C to stop)", opts)

			err := client.StreamLogs(req, func(entry *models.StreamLogEntry) error {
				if opts.Format == utils.OutputFormatJSON {
					utils.PrintJSON(entry)
				} else {
					timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
					fmt.Printf("%s [%s] %s\n", timestamp, entry.Level, entry.Message)
				}
				return nil
			})

			if err != nil {
				utils.HandleError(err, "Error streaming logs", opts)
				return
			}
		} else {
			// Get historical logs
			utils.PrintInfo("Fetching logs...", opts)

			logsResp, err := client.GetLogs(req)
			if err != nil {
				utils.HandleError(err, "Error fetching logs", opts)
				return
			}

			if opts.Format == utils.OutputFormatJSON {
				utils.PrintJSON(logsResp)
			} else {
				if len(logsResp.Logs) == 0 {
					utils.PrintWarning("No logs found", opts)
					return
				}

				for _, log := range logsResp.Logs {
					timestamp := log.Timestamp.Format("2006-01-02 15:04:05")
					fmt.Printf("%s [%s] %s\n", timestamp, log.Level, log.Message)
				}

				utils.PrintSuccess(fmt.Sprintf("Found %d log entries", len(logsResp.Logs)), opts)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
}

// printLogEntry formats and prints a log entry with colors (shared with project logs)
func printLogEntry(entry *models.LogEntry) {
	// Format timestamp
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")

	// Get color for log level
	levelColor := entry.Level.GetColor()
	resetColor := models.ResetColor()

	// Format level with fixed width
	level := fmt.Sprintf("%-5s", strings.ToUpper(string(entry.Level)))

	// Build source info
	sourceInfo := ""
	if entry.Source != "" {
		sourceInfo = fmt.Sprintf("[%s]", entry.Source)
	}
	if entry.Container != "" {
		if sourceInfo != "" {
			sourceInfo += fmt.Sprintf("[%s]", entry.Container)
		} else {
			sourceInfo = fmt.Sprintf("[%s]", entry.Container)
		}
	}

	// Print formatted log entry
	fmt.Printf("%s %s%s%s %s %s\n",
		timestamp,
		levelColor, level, resetColor,
		sourceInfo,
		entry.Message)
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// Add flags
	logsCmd.Flags().StringP("service", "s", "", "Filter logs by service name")
	logsCmd.Flags().StringP("container", "c", "", "Filter logs by container name")
	logsCmd.Flags().BoolP("follow", "f", false, "Stream logs in real-time")
	logsCmd.Flags().IntP("lines", "n", 100, "Number of lines to show")
	logsCmd.Flags().String("since", "", "Show logs since timestamp (RFC3339)")
	logsCmd.Flags().String("until", "", "Show logs until timestamp (RFC3339)")
}
