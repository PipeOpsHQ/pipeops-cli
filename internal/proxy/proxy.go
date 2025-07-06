package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/models"
)

// Manager handles proxy sessions
type Manager struct {
	proxies map[string]*ProxySession
	mutex   sync.RWMutex
}

// ProxySession represents an active proxy session
type ProxySession struct {
	ID          string
	Target      models.ProxyTarget
	LocalPort   int
	RemoteHost  string
	RemotePort  int
	Status      string
	StartedAt   time.Time
	BytesIn     int64
	BytesOut    int64
	Connections int
	listener    net.Listener
	cancel      context.CancelFunc
	mutex       sync.RWMutex
}

// NewManager creates a new proxy manager
func NewManager() *Manager {
	return &Manager{
		proxies: make(map[string]*ProxySession),
	}
}

// StartProxy starts a new proxy session
func (m *Manager) StartProxy(target models.ProxyTarget, localPort int, remoteHost string, remotePort int) (*models.ProxyResponse, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Generate proxy ID
	proxyID := fmt.Sprintf("proxy-%d", time.Now().Unix())

	// Find available local port if not specified
	if localPort == 0 {
		var err error
		localPort, err = findAvailablePort()
		if err != nil {
			return nil, fmt.Errorf("failed to find available port: %w", err)
		}
	}

	// Create listener on local port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", localPort, err)
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Create proxy session
	session := &ProxySession{
		ID:         proxyID,
		Target:     target,
		LocalPort:  localPort,
		RemoteHost: remoteHost,
		RemotePort: remotePort,
		Status:     "active",
		StartedAt:  time.Now(),
		listener:   listener,
		cancel:     cancel,
	}

	// Store the session
	m.proxies[proxyID] = session

	// Start handling connections in a goroutine
	go session.handleConnections(ctx)

	return &models.ProxyResponse{
		ProxyID:    proxyID,
		Target:     target,
		LocalPort:  localPort,
		RemoteHost: remoteHost,
		RemotePort: remotePort,
		Status:     "active",
		StartedAt:  session.StartedAt.Format(time.RFC3339),
	}, nil
}

// StopProxy stops a proxy session
func (m *Manager) StopProxy(proxyID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.proxies[proxyID]
	if !exists {
		return fmt.Errorf("proxy %s not found", proxyID)
	}

	// Cancel the context
	session.cancel()

	// Close the listener
	if session.listener != nil {
		session.listener.Close()
	}

	// Update status
	session.mutex.Lock()
	session.Status = "stopped"
	session.mutex.Unlock()

	// Remove from active proxies
	delete(m.proxies, proxyID)

	return nil
}

// GetProxyStatus returns the status of a proxy
func (m *Manager) GetProxyStatus(proxyID string) (*models.ProxyStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.proxies[proxyID]
	if !exists {
		return nil, fmt.Errorf("proxy %s not found", proxyID)
	}

	session.mutex.RLock()
	defer session.mutex.RUnlock()

	return &models.ProxyStatus{
		ProxyID:       session.ID,
		Status:        session.Status,
		LocalPort:     session.LocalPort,
		RemoteHost:    session.RemoteHost,
		RemotePort:    session.RemotePort,
		BytesIn:       session.BytesIn,
		BytesOut:      session.BytesOut,
		ConnectionsIn: session.Connections,
		StartedAt:     session.StartedAt.Format(time.RFC3339),
	}, nil
}

// ListProxies returns all active proxy sessions
func (m *Manager) ListProxies() *models.ListProxiesResponse {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var proxies []models.ProxyStatus
	for _, session := range m.proxies {
		session.mutex.RLock()
		status := models.ProxyStatus{
			ProxyID:       session.ID,
			Status:        session.Status,
			LocalPort:     session.LocalPort,
			RemoteHost:    session.RemoteHost,
			RemotePort:    session.RemotePort,
			BytesIn:       session.BytesIn,
			BytesOut:      session.BytesOut,
			ConnectionsIn: session.Connections,
			StartedAt:     session.StartedAt.Format(time.RFC3339),
		}
		session.mutex.RUnlock()
		proxies = append(proxies, status)
	}

	return &models.ListProxiesResponse{
		Proxies: proxies,
		Total:   len(proxies),
	}
}

// StopAllProxies stops all active proxy sessions
func (m *Manager) StopAllProxies() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for proxyID := range m.proxies {
		session := m.proxies[proxyID]
		session.cancel()
		if session.listener != nil {
			session.listener.Close()
		}
	}

	// Clear all proxies
	m.proxies = make(map[string]*ProxySession)

	return nil
}

// handleConnections handles incoming connections for a proxy session
func (s *ProxySession) handleConnections(ctx context.Context) {
	defer s.listener.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Accept connection with timeout
			conn, err := s.listener.Accept()
			if err != nil {
				s.mutex.Lock()
				s.Status = "error"
				s.mutex.Unlock()
				return
			}

			// Handle connection in a separate goroutine
			go s.handleConnection(conn, ctx)
		}
	}
}

// handleConnection handles a single connection
func (s *ProxySession) handleConnection(localConn net.Conn, ctx context.Context) {
	defer localConn.Close()

	// Increment connection count
	s.mutex.Lock()
	s.Connections++
	s.mutex.Unlock()

	// Defer decrement
	defer func() {
		s.mutex.Lock()
		s.Connections--
		s.mutex.Unlock()
	}()

	// Connect to remote host
	remoteConn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", s.RemoteHost, s.RemotePort),
		10*time.Second)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	// Copy data bidirectionally
	go func() {
		written, _ := io.Copy(remoteConn, localConn)
		s.mutex.Lock()
		s.BytesOut += written
		s.mutex.Unlock()
	}()

	written, _ := io.Copy(localConn, remoteConn)
	s.mutex.Lock()
	s.BytesIn += written
	s.mutex.Unlock()
}

// findAvailablePort finds an available local port
func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// IsPortAvailable checks if a port is available
func IsPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

// GetPortFromString parses a port from string and validates it
func GetPortFromString(portStr string) (int, error) {
	if portStr == "" {
		return 0, nil // Auto-assign
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port number: %s", portStr)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port number must be between 1 and 65535, got: %d", port)
	}

	return port, nil
}
