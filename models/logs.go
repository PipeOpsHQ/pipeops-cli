package models

import "time"

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// LogEntry represents a single log entry
type LogEntry struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Level     LogLevel          `json:"level"`
	Message   string            `json:"message"`
	Source    string            `json:"source"`           // e.g., "app", "nginx", "database"
	Container string            `json:"container"`        // container name/id
	Pod       string            `json:"pod"`              // kubernetes pod name
	Node      string            `json:"node"`             // kubernetes node name
	Labels    map[string]string `json:"labels,omitempty"` // additional metadata
}

// LogsResponse represents the response from the logs API
type LogsResponse struct {
	Logs       []LogEntry `json:"logs"`
	TotalCount int        `json:"total_count"`
	HasMore    bool       `json:"has_more"`
	NextCursor string     `json:"next_cursor,omitempty"`
}

// LogsRequest represents a request for logs with filtering options
type LogsRequest struct {
	ProjectID string     `json:"project_id"`
	AddonID   string     `json:"addon_id,omitempty"`  // optional, for addon logs
	Level     LogLevel   `json:"level,omitempty"`     // filter by minimum level
	Source    string     `json:"source,omitempty"`    // filter by source
	Container string     `json:"container,omitempty"` // filter by container
	Since     *time.Time `json:"since,omitempty"`     // logs since this time
	Until     *time.Time `json:"until,omitempty"`     // logs until this time
	Limit     int        `json:"limit,omitempty"`     // max number of logs to return
	Cursor    string     `json:"cursor,omitempty"`    // pagination cursor
	Follow    bool       `json:"follow,omitempty"`    // stream logs in real-time
	Tail      int        `json:"tail,omitempty"`      // get last N lines
}

// StreamLogEntry represents a log entry in a streaming context
type StreamLogEntry struct {
	LogEntry
	StreamID string `json:"stream_id"` // unique identifier for this stream
	EOF      bool   `json:"eof"`       // indicates end of stream
}

// LogsStreamResponse represents a streaming logs response
type LogsStreamResponse struct {
	Entry *StreamLogEntry `json:"entry,omitempty"`
	Error string          `json:"error,omitempty"`
}

// GetLogLevelColor returns ANSI color code for log levels
func (l LogLevel) GetColor() string {
	switch l {
	case LogLevelDebug:
		return "\033[36m" // Cyan
	case LogLevelInfo:
		return "\033[32m" // Green
	case LogLevelWarn:
		return "\033[33m" // Yellow
	case LogLevelError:
		return "\033[31m" // Red
	case LogLevelFatal:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// ResetColor returns ANSI reset code
func ResetColor() string {
	return "\033[0m"
}
