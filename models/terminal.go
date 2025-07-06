package models

// ExecRequest represents a request to execute a command in a container
type ExecRequest struct {
	ProjectID   string            `json:"project_id"`
	AddonID     string            `json:"addon_id,omitempty"`    // optional, for addon containers
	ServiceName string            `json:"service_name"`          // service to execute in
	Container   string            `json:"container,omitempty"`   // specific container name (if service has multiple)
	Command     []string          `json:"command"`               // command to execute
	Interactive bool              `json:"interactive"`           // whether to allocate a TTY
	Environment map[string]string `json:"environment,omitempty"` // environment variables
	WorkingDir  string            `json:"working_dir,omitempty"` // working directory
	User        string            `json:"user,omitempty"`        // user to run command as
}

// ExecResponse represents the response when starting an exec session
type ExecResponse struct {
	ExecID       string `json:"exec_id"`       // unique identifier for this exec session
	WebSocketURL string `json:"websocket_url"` // WebSocket URL for interactive session
	Status       string `json:"status"`        // "starting", "running", "completed", "error"
	StartedAt    string `json:"started_at"`    // when the exec session was started
}

// ExecMessage represents a message in the WebSocket stream
type ExecMessage struct {
	Type      string `json:"type"`                // "stdin", "stdout", "stderr", "resize", "exit"
	Data      string `json:"data"`                // message data (base64 encoded for binary data)
	Timestamp string `json:"timestamp"`           // message timestamp
	ExitCode  int    `json:"exit_code,omitempty"` // exit code (only for "exit" type)
}

// ResizeMessage represents a terminal resize message
type ResizeMessage struct {
	Type string `json:"type"` // "resize"
	Cols int    `json:"cols"` // terminal columns
	Rows int    `json:"rows"` // terminal rows
}

// ExecStatus represents the status of an exec session
type ExecStatus struct {
	ExecID    string `json:"exec_id"`
	Status    string `json:"status"` // "running", "completed", "error"
	StartedAt string `json:"started_at"`
	ExitCode  int    `json:"exit_code,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ListExecResponse represents the response when listing exec sessions
type ListExecResponse struct {
	Sessions []ExecStatus `json:"sessions"`
	Total    int          `json:"total"`
}

// ShellRequest represents a request to start an interactive shell
type ShellRequest struct {
	ProjectID   string            `json:"project_id"`
	AddonID     string            `json:"addon_id,omitempty"`    // optional, for addon containers
	ServiceName string            `json:"service_name"`          // service to connect to
	Container   string            `json:"container,omitempty"`   // specific container name
	Shell       string            `json:"shell,omitempty"`       // shell to use (bash, sh, zsh, etc.)
	Environment map[string]string `json:"environment,omitempty"` // environment variables
	WorkingDir  string            `json:"working_dir,omitempty"` // working directory
	User        string            `json:"user,omitempty"`        // user to run shell as
	Cols        int               `json:"cols,omitempty"`        // terminal columns
	Rows        int               `json:"rows,omitempty"`        // terminal rows
}

// ShellResponse represents the response when starting a shell session
type ShellResponse struct {
	SessionID    string `json:"session_id"`    // unique identifier for this shell session
	WebSocketURL string `json:"websocket_url"` // WebSocket URL for interactive session
	Status       string `json:"status"`        // "starting", "running", "completed", "error"
	StartedAt    string `json:"started_at"`    // when the shell session was started
}

// ContainerInfo represents information about a container that can be accessed
type ContainerInfo struct {
	Name         string            `json:"name"`
	ServiceName  string            `json:"service_name"`  // parent service name
	Image        string            `json:"image"`         // container image
	Status       string            `json:"status"`        // "running", "stopped", "restarting", etc.
	RestartCount int               `json:"restart_count"` // number of restarts
	Labels       map[string]string `json:"labels,omitempty"`
	CreatedAt    string            `json:"created_at"`
	StartedAt    string            `json:"started_at,omitempty"`
}

// ListContainersResponse represents available containers for a project/addon
type ListContainersResponse struct {
	Containers []ContainerInfo `json:"containers"`
	Total      int             `json:"total"`
}

// LogsExecRequest represents a request to get logs from an exec session
type LogsExecRequest struct {
	ExecID string `json:"exec_id"`
	Lines  int    `json:"lines,omitempty"` // number of lines to get from the end
}

// LogsExecResponse represents logs from an exec session
type LogsExecResponse struct {
	ExecID string   `json:"exec_id"`
	Logs   []string `json:"logs"`
}
