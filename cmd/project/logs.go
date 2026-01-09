package project

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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [project-id]",
	Short: "📄 View logs for a project",
	Long: `📄 The "logs" command retrieves and displays logs for a specific project.
Use --follow to stream logs in real-time.

Examples:
  - View recent logs:
    pipeops project logs proj-123

  - Follow logs in real-time:
    pipeops project logs proj-123 --follow

  - View logs from specific time:
    pipeops project logs proj-123 --since "2024-01-01T10:00:00Z"

  - Get last 100 lines:
    pipeops project logs proj-123 --tail 100
    
  - Interactive project selection:
    pipeops project logs`,
	Run: func(cmd *cobra.Command, args []string) {
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

		var projectID string
		if len(args) == 1 {
			projectID = args[0]
			// Validate project ID
			if err := validation.ValidateProjectID(projectID); err != nil {
				fmt.Printf("❌ Invalid project ID: %v\n", err)
				return
			}
		} else {
			// Interactive project selection
			projectsResp, err := client.GetProjects()
			if err != nil {
				fmt.Printf("❌ Error fetching projects: %v\n", err)
				return
			}

			if len(projectsResp.Projects) == 0 {
				fmt.Println("⚠️ No projects found")
				return
			}

			var options []string
			for _, p := range projectsResp.Projects {
				options = append(options, fmt.Sprintf("%s (%s)", p.Name, p.ID))
			}

			idx, err := selectProject(options)
			if err != nil {
				fmt.Printf("❌ Selection cancelled: %v\n", err)
				return
			}

			projectID = projectsResp.Projects[idx].ID
		}

		// Parse flags
		sinceStr, _ := cmd.Flags().GetString("since")
		untilStr, _ := cmd.Flags().GetString("until")
		limitStr, _ := cmd.Flags().GetString("limit")
		tail, _ := cmd.Flags().GetInt("tail")
		follow, _ := cmd.Flags().GetBool("follow")

		// Build logs request
		req := &models.LogsRequest{
			ProjectID: projectID,
			Tail:      tail,
			Follow:    follow,
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
			fmt.Printf("Streaming logs for project %s", projectID)
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
			fmt.Printf("Fetching logs for project %s...\n", projectID)

			resp, err := client.GetLogs(req)
			if err != nil {
				fmt.Printf("❌ Error fetching logs: %v\n", err)
				return
			}

			if len(resp.Logs) == 0 {
				fmt.Println("No logs found for the specified criteria.")
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

// selectProject prompts user to select from a list of projects
func selectProject(options []string) (int, error) {
	prompt := promptui.Select{
		Label: "Select a project",
		Items: options,
		Size:  10,
	}
	idx, _, err := prompt.Run()
	return idx, err
}

func init() {
	logsCmd.Flags().String("since", "", "Show logs since timestamp (RFC3339 format)")
	logsCmd.Flags().String("until", "", "Show logs until timestamp (RFC3339 format)")
	logsCmd.Flags().String("limit", "", "Maximum number of logs to retrieve")
	logsCmd.Flags().IntP("tail", "t", 100, "Number of recent log lines to show")
	logsCmd.Flags().BoolP("follow", "f", false, "Stream logs in real-time")
}

// printLogEntry formats and prints a log entry
func printLogEntry(entry *models.LogEntry) {
	// Format timestamp
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")

	// Get color for log level
	levelColor := entry.Level.GetColor()
	resetColor := models.ResetColor()

	// Format level with fixed width
	level := fmt.Sprintf("%-5s", strings.ToUpper(string(entry.Level)))

	// Print formatted log entry
	fmt.Printf("%s %s%s%s %s\n",
		timestamp,
		levelColor, level, resetColor,
		entry.Message)
}

func (p *projectModel) logs() {
	p.rootCmd.AddCommand(logsCmd)
}
