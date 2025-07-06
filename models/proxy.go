package models

// ProxyTarget represents a target service to proxy to
type ProxyTarget struct {
	ProjectID   string `json:"project_id"`
	AddonID     string `json:"addon_id,omitempty"` // optional, for addon services
	ServiceName string `json:"service_name"`       // name of the service to proxy to
	Port        int    `json:"port"`               // target port on the service
}

// ProxyRequest represents a request to start a proxy
type ProxyRequest struct {
	Target    ProxyTarget `json:"target"`
	LocalPort int         `json:"local_port,omitempty"` // desired local port (0 for auto-assign)
}

// ProxyResponse represents the response when starting a proxy
type ProxyResponse struct {
	ProxyID    string      `json:"proxy_id"`    // unique identifier for this proxy session
	Target     ProxyTarget `json:"target"`      // target information
	LocalPort  int         `json:"local_port"`  // actual local port assigned
	RemoteHost string      `json:"remote_host"` // remote host to connect to
	RemotePort int         `json:"remote_port"` // remote port to connect to
	Status     string      `json:"status"`      // proxy status
	StartedAt  string      `json:"started_at"`  // when the proxy was started
}

// ProxyStatus represents the current status of a proxy
type ProxyStatus struct {
	ProxyID       string `json:"proxy_id"`
	Status        string `json:"status"` // "active", "stopped", "error"
	LocalPort     int    `json:"local_port"`
	RemoteHost    string `json:"remote_host"`
	RemotePort    int    `json:"remote_port"`
	BytesIn       int64  `json:"bytes_in"`       // bytes received from remote
	BytesOut      int64  `json:"bytes_out"`      // bytes sent to remote
	ConnectionsIn int    `json:"connections_in"` // current inbound connections
	StartedAt     string `json:"started_at"`
	LastActivity  string `json:"last_activity,omitempty"`
	Error         string `json:"error,omitempty"`
}

// ListProxiesResponse represents the response when listing active proxies
type ListProxiesResponse struct {
	Proxies []ProxyStatus `json:"proxies"`
	Total   int           `json:"total"`
}

// ProxyStopRequest represents a request to stop a proxy
type ProxyStopRequest struct {
	ProxyID string `json:"proxy_id"`
}

// ServiceInfo represents information about a service that can be proxied
type ServiceInfo struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // "web", "api", "database", etc.
	Port        int               `json:"port"`
	Protocol    string            `json:"protocol"` // "http", "https", "tcp", "udp"
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels,omitempty"`
	Health      string            `json:"health"` // "healthy", "unhealthy", "unknown"
}

// ListServicesResponse represents available services for a project/addon
type ListServicesResponse struct {
	Services []ServiceInfo `json:"services"`
	Total    int           `json:"total"`
}
