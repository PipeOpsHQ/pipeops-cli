package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	OutputFormatTable OutputFormat = "table"
	OutputFormatJSON  OutputFormat = "json"
)

// OutputOptions contains options for output formatting
type OutputOptions struct {
	Format  OutputFormat
	Quiet   bool
	Verbose bool
}

// GetOutputOptions extracts output options from command flags
func GetOutputOptions(cmd *cobra.Command) OutputOptions {
	jsonOutput, _ := cmd.Flags().GetBool("json")
	quiet, _ := cmd.Flags().GetBool("quiet")
	verbose, _ := cmd.Flags().GetBool("verbose")

	format := OutputFormatTable
	if jsonOutput {
		format = OutputFormatJSON
	}

	return OutputOptions{
		Format:  format,
		Quiet:   quiet,
		Verbose: verbose,
	}
}

// PrintSuccess prints a success message with emoji
func PrintSuccess(message string, opts OutputOptions) {
	if opts.Quiet {
		return
	}
	if opts.Format == OutputFormatJSON {
		return // JSON output doesn't include success messages
	}
	fmt.Printf("✅ %s\n", message)
}

// PrintError prints an error message with emoji
func PrintError(message string, opts OutputOptions) {
	if opts.Format == OutputFormatJSON {
		errorObj := map[string]interface{}{
			"error":   true,
			"message": message,
		}
		jsonBytes, _ := json.MarshalIndent(errorObj, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("❌ %s\n", message)
	}
}

// PrintInfo prints an info message with emoji
func PrintInfo(message string, opts OutputOptions) {
	if opts.Quiet {
		return
	}
	if opts.Format == OutputFormatJSON {
		return // JSON output doesn't include info messages
	}
	fmt.Printf("🔍 %s\n", message)
}

// PrintWarning prints a warning message with emoji
func PrintWarning(message string, opts OutputOptions) {
	if opts.Quiet {
		return
	}
	if opts.Format == OutputFormatJSON {
		return // JSON output doesn't include warning messages
	}
	fmt.Printf("⚠️  %s\n", message)
}

// PrintJSON prints data as JSON
func PrintJSON(data interface{}) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}

// PrintTable prints data in a table format
func PrintTable(headers []string, rows [][]string, opts OutputOptions) {
	if opts.Format == OutputFormatJSON {
		// Convert table to JSON format
		var jsonData []map[string]interface{}
		for _, row := range rows {
			rowData := make(map[string]interface{})
			for i, header := range headers {
				if i < len(row) {
					rowData[strings.ToLower(header)] = row[i]
				}
			}
			jsonData = append(jsonData, rowData)
		}
		PrintJSON(jsonData)
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	printRow(headers, widths)

	// Print separator
	var separators []string
	for _, width := range widths {
		separators = append(separators, strings.Repeat("-", width))
	}
	printRow(separators, widths)

	// Print rows
	for _, row := range rows {
		printRow(row, widths)
	}
}

// printRow prints a single row with proper spacing
func printRow(row []string, widths []int) {
	var parts []string
	for i, cell := range row {
		if i < len(widths) {
			parts = append(parts, fmt.Sprintf("%-*s", widths[i], cell))
		}
	}
	fmt.Println(strings.Join(parts, " | "))
}

// PrintProjectContextWithOptions prints project context information with output options
func PrintProjectContextWithOptions(projectID string, opts OutputOptions) {
	if opts.Format == OutputFormatJSON || opts.Quiet {
		return
	}

	if IsLinkedProject() {
		if linkedID, err := GetLinkedProject(); err == nil && linkedID == projectID {
			fmt.Printf("📂 Using linked project: %s\n", projectID)
		} else {
			fmt.Printf("🎯 Using project: %s (overriding linked project)\n", projectID)
		}
	} else {
		fmt.Printf("🎯 Using project: %s\n", projectID)
	}
}

// FormatDate formats a date for display
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDateShort formats a date in short format
func FormatDateShort(t time.Time) string {
	return t.Format("2006-01-02")
}

// TruncateString truncates a string to the specified length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// GetStatusIcon returns an emoji icon for status
func GetStatusIcon(status string) string {
	switch strings.ToLower(status) {
	case "active", "running", "healthy", "success":
		return "🟢"
	case "deploying", "building", "starting", "pending":
		return "🟡"
	case "stopped", "inactive", "paused":
		return "⚪"
	case "error", "failed", "crashed":
		return "🔴"
	default:
		return "⚫"
	}
}

// RequireAuth checks if user is authenticated and prints error if not
func RequireAuth(client interface{ IsAuthenticated() bool }, opts OutputOptions) bool {
	if !client.IsAuthenticated() {
		PrintError("You are not logged in. Please run 'pipeops auth login' first.", opts)
		return false
	}
	return true
}

// HandleError handles errors consistently across commands
func HandleError(err error, message string, opts OutputOptions) {
	if err != nil {
		PrintError(fmt.Sprintf("%s: %v", message, err), opts)
		os.Exit(1)
	}
}

// PromptUser prompts user for input with a message
func PromptUser(message string) (string, error) {
	fmt.Print(message)
	var input string
	_, err := fmt.Scanln(&input)
	return input, err
}

// PromptUserWithDefault prompts user for input with a default value
func PromptUserWithDefault(message, defaultValue string) string {
	fmt.Printf("%s [%s]: ", message, defaultValue)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return defaultValue
	}
	return input
}

// ConfirmAction asks user for confirmation
func ConfirmAction(message string) bool {
	fmt.Printf("%s (y/N): ", message)
	var input string
	fmt.Scanln(&input)
	return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}

// WithJSONOutput wraps a function to support JSON output
func WithJSONOutput(fn func() (interface{}, error), opts OutputOptions) error {
	data, err := fn()
	if err != nil {
		HandleError(err, "Operation failed", opts)
		return err
	}

	if opts.Format == OutputFormatJSON {
		return PrintJSON(data)
	}

	return nil
}
