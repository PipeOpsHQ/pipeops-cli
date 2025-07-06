package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/internal/pipeops"
	"github.com/PipeOpsHQ/pipeops-cli/internal/validation"
	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [project-id]",
	Short: "📄 View logs for the current or specified project",
	Long: `📄 View logs for a project. If no project ID is provided, uses the linked project
from the current directory (set with 'pipeops link').

Examples:
  - View logs for linked project:
    pipeops logs

  - View logs for specific project:
    pipeops logs proj-123

  - Follow logs in real-time:
    pipeops logs --follow

  - Filter by log level:
    pipeops logs --level error

  - View addon logs:
    pipeops logs --addon addon-456`,
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
				return
			}
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
		level, _ := cmd.Flags().GetString("level")
		source, _ := cmd.Flags().GetString("source")
		container, _ := cmd.Flags().GetString("container")
		sinceStr, _ := cmd.Flags().GetString("since")
		untilStr, _ := cmd.Flags().GetString("until")
		limitStr, _ := cmd.Flags().GetString("limit")
		tail, _ := cmd.Flags().GetInt("tail")
		follow, _ := cmd.Flags().GetBool("follow")
		addonID, _ := cmd.Flags().GetString("addon")

		// Build logs request
		req := &models.LogsRequest{
			ProjectID: projectID,
			AddonID:   addonID,
			Tail:      tail,
			Follow:    follow,
		}

		// Parse level
		if level != "" {
			req.Level = models.LogLevel(level)
		}

		// Parse source and container
		if source != "" {
			req.Source = source
		}
		if container != "" {
			req.Container = container
		}

		// Parse time filters
		if sinceStr != "" {
			since, err := time.Parse(time.RFC3339, sinceStr)
			if err != nil {
				fmt.Printf("❌ Invalid since time format. Use RFC3339 format (e.g., 2024-01-01T10:00:00Z): %v\n", err)
				return
			}
			req.Since = &since
		}

		if untilStr != "" {
			until, err := time.Parse(time.RFC3339, untilStr)
			if err != nil {
				fmt.Printf("❌ Invalid until time format. Use RFC3339 format (e.g., 2024-01-01T10:00:00Z): %v\n", err)
				return
			}
			req.Until = &until
		}

		// Parse limit
		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				fmt.Printf("❌ Invalid limit: %v\n", err)
				return
			}
			req.Limit = limit
		}

		if follow {
			// Stream logs in real-time
			fmt.Printf("🔄 Streaming logs")
			if addonID != "" {
				fmt.Printf(" (addon: %s)", addonID)
			}
			fmt.Println("... (Press Ctrl+C to stop)")

			// Set up signal handling
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			// Channel to signal completion
			doneChan := make(chan error, 1)

			// Start streaming in a goroutine
			go func() {
				doneChan <- client.StreamLogs(req, func(entry *models.StreamLogEntry) error {
					printLogEntry(&entry.LogEntry)
					return nil
				})
			}()

			// Wait for completion or signal
			select {
			case err := <-doneChan:
				if err != nil {
					fmt.Printf("\n❌ Error streaming logs: %v\n", err)
				} else {
					fmt.Println("\n✅ Log stream ended.")
				}
			case <-sigChan:
				fmt.Println("\n🛑 Log streaming stopped by user.")
			}
		} else {
			// Get logs once
			fmt.Printf("🔍 Fetching logs")
			if addonID != "" {
				fmt.Printf(" (addon: %s)", addonID)
			}
			fmt.Println("...")

			resp, err := client.GetLogs(req)
			if err != nil {
				fmt.Printf("❌ Error fetching logs: %v\n", err)
				return
			}

			if len(resp.Logs) == 0 {
				fmt.Println("📭 No logs found for the specified criteria.")
				return
			}

			// Display logs
			for _, entry := range resp.Logs {
				printLogEntry(&entry)
			}

			fmt.Printf("\n✅ Found %d log entries", len(resp.Logs))
			if resp.HasMore {
				fmt.Printf(" (more available - use --limit to get more or --follow to stream)")
			}
			fmt.Println()
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
	logsCmd.Flags().StringP("level", "l", "", "Filter by log level (debug, info, warn, error, fatal)")
	logsCmd.Flags().StringP("source", "s", "", "Filter by log source")
	logsCmd.Flags().StringP("container", "c", "", "Filter by container name")
	logsCmd.Flags().String("since", "", "Show logs since timestamp (RFC3339 format)")
	logsCmd.Flags().String("until", "", "Show logs until timestamp (RFC3339 format)")
	logsCmd.Flags().String("limit", "", "Maximum number of logs to retrieve")
	logsCmd.Flags().IntP("tail", "t", 100, "Number of recent log lines to show")
	logsCmd.Flags().BoolP("follow", "f", false, "Stream logs in real-time")
	logsCmd.Flags().StringP("addon", "a", "", "Get logs for a specific addon")
}
