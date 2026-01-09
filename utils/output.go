package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
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

// PrintSuccess prints a success message with emoji and color
func PrintSuccess(message string, opts OutputOptions) {
	if opts.Quiet {
		return
	}
	if opts.Format == OutputFormatJSON {
		return // JSON output doesn't include success messages
	}
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", green("✅"), green(message))
}

// PrintError prints an error message with emoji and color
func PrintError(message string, opts OutputOptions) {
	if opts.Format == OutputFormatJSON {
		errorObj := map[string]interface{}{
			"error":   true,
			"message": message,
		}
		jsonBytes, _ := json.MarshalIndent(errorObj, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s %s\n", red("❌"), red(message))
	}
}

// PrintInfo prints an info message with emoji and color
func PrintInfo(message string, opts OutputOptions) {
	if opts.Quiet {
		return
	}
	if opts.Format == OutputFormatJSON {
		return // JSON output doesn't include info messages
	}
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s %s\n", cyan("🔍"), cyan(message))
}

// PrintWarning prints a warning message with emoji and color
func PrintWarning(message string, opts OutputOptions) {
	if opts.Quiet {
		return
	}
	if opts.Format == OutputFormatJSON {
		return // JSON output doesn't include warning messages
	}
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%s %s\n", yellow("⚠️ "), yellow(message))
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

// PrintTable prints data in a table format using tablewriter
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	// Add colors to header
	headerColors := make([]tablewriter.Colors, len(headers))
	for i := range headerColors {
		headerColors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}
	}
	table.SetHeaderColor(headerColors...)

	table.AppendBulk(rows)
	table.Render()
}

// PrintProjectContextWithOptions prints project context information with output options
func PrintProjectContextWithOptions(projectID string, opts OutputOptions) {
	if opts.Format == OutputFormatJSON || opts.Quiet {
		return
	}

	bold := color.New(color.Bold).SprintFunc()

	if IsLinkedProject() {
		if linkedID, err := GetLinkedProject(); err == nil && linkedID == projectID {
			fmt.Printf("📂 Using linked project: %s\n", bold(projectID))
		} else {
			fmt.Printf("🎯 Using project: %s (overriding linked project)\n", bold(projectID))
		}
	} else {
		fmt.Printf("🎯 Using project: %s\n", bold(projectID))
	}
}

// FormatDate formats a date for display
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDateShort formats a date in short format
func FormatDateShort(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
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

// HandleAuthError handles authentication errors with specific messaging
func HandleAuthError(err error, opts OutputOptions) bool {
	if err == nil {
		return true
	}

	// Check for specific authentication error types
	errorStr := err.Error()

	// Check for token expired
	if strings.Contains(errorStr, "expired") || strings.Contains(errorStr, "expiration") {
		PrintError("Your session has expired. Please run 'pipeops auth login' to authenticate again.", opts)
		return false
	}

	// Check for token revoked
	if strings.Contains(errorStr, "revoked") || strings.Contains(errorStr, "invalidated") {
		PrintError("Your session has been revoked. Please run 'pipeops auth login' to authenticate again.", opts)
		return false
	}

	// Check for invalid token
	if strings.Contains(errorStr, "invalid") || strings.Contains(errorStr, "malformed") {
		PrintError("Your authentication token is invalid. Please run 'pipeops auth login' to authenticate again.", opts)
		return false
	}

	// Check for refresh failed
	if strings.Contains(errorStr, "refresh") && strings.Contains(errorStr, "failed") {
		PrintError("Failed to refresh your session. Please run 'pipeops auth login' to authenticate again.", opts)
		return false
	}

	// Check if it's a general authentication error
	if strings.Contains(errorStr, "authentication") ||
		strings.Contains(errorStr, "unauthorized") ||
		strings.Contains(errorStr, "401") ||
		strings.Contains(errorStr, "invalid token") {
		PrintError("Authentication failed. Please run 'pipeops auth login' to authenticate again.", opts)
		return false
	}

	// Not an authentication error, return true to let caller handle it
	return true
}

// HandleError handles errors consistently across commands
func HandleError(err error, message string, opts OutputOptions) {
	if err != nil {
		PrintError(fmt.Sprintf("%s: %v", message, err), opts)
		os.Exit(1)
	}
}

// PromptUser prompts user for input with a message using promptui
func PromptUser(message string) (string, error) {
	// Clean the message (remove : or space at end)
	message = strings.TrimSuffix(strings.TrimSpace(message), ":")

	prompt := promptui.Prompt{
		Label: message,
	}
	return prompt.Run()
}

// PromptUserWithDefault prompts user for input with a default value using promptui
func PromptUserWithDefault(message, defaultValue string) string {
	// Clean the message
	message = strings.TrimSuffix(strings.TrimSpace(message), ":")

	prompt := promptui.Prompt{
		Label:     message,
		Default:   defaultValue,
		AllowEdit: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return defaultValue
	}
	return result
}

// ConfirmAction asks user for confirmation using promptui
func ConfirmAction(message string) bool {
	// Clean the message
	message = strings.TrimSuffix(strings.TrimSpace(message), "?")

	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	return err == nil
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

// SelectOption prompts user to select from a list of options
func SelectOption(label string, options []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
		Size:  10,
	}

	return prompt.Run()
}

// StartSpinner starts a new spinner with the given message
func StartSpinner(message string, opts OutputOptions) interface{} {
	if opts.Format == OutputFormatJSON || opts.Quiet {
		return nil
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s", message)
	s.Color("cyan")
	s.Start()
	return s
}

// StopSpinner stops the spinner if it exists
func StopSpinner(s interface{}) {
	if s == nil {
		return
	}
	if spin, ok := s.(*spinner.Spinner); ok {
		spin.Stop()
	}
}
