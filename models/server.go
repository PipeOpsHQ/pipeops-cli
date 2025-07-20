package models

import "time"

// Server represents a PipeOps server
type Server struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Status      string            `json:"status"` // "running", "stopped", "starting", "stopping", "error"
	Type        string            `json:"type"`   // "k3s", "docker", "kubernetes", etc.
	Region      string            `json:"region"`
	IP          string            `json:"ip,omitempty"`
	Port        int               `json:"port,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	UserID      string            `json:"user_id"`
}

// ServersResponse represents the response from the servers API
type ServersResponse struct {
	Servers []Server `json:"servers"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
}

// ServerCreateRequest represents the request to create a server
type ServerCreateRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type"`
	Region      string            `json:"region"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// ServerUpdateRequest represents the request to update a server
type ServerUpdateRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Status      string            `json:"status,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// ServerStatus represents the current status of a server
type ServerStatus struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	CPUUsage     float64   `json:"cpu_usage,omitempty"`
	MemoryUsage  float64   `json:"memory_usage,omitempty"`
	DiskUsage    float64   `json:"disk_usage,omitempty"`
	Uptime       string    `json:"uptime,omitempty"`
	LastSeen     time.Time `json:"last_seen,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
}
