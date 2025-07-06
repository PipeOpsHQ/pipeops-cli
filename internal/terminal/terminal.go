package terminal

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/models"
	"github.com/gorilla/websocket"
	"golang.org/x/term"
)

// Session represents a terminal session
type Session struct {
	ID            string
	WebSocketURL  string
	conn          *websocket.Conn
	cancel        context.CancelFunc
	isInteractive bool
	originalState *term.State
}

// Manager handles terminal sessions
type Manager struct {
	sessions map[string]*Session
}

// NewManager creates a new terminal manager
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

// StartExecSession starts a new exec session
func (m *Manager) StartExecSession(execID string, websocketURL string, interactive bool) (*Session, error) {
	// Parse WebSocket URL
	u, err := url.Parse(websocketURL)
	if err != nil {
		return nil, fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Create session
	session := &Session{
		ID:            execID,
		WebSocketURL:  websocketURL,
		conn:          conn,
		cancel:        cancel,
		isInteractive: interactive,
	}

	// Store session
	m.sessions[execID] = session

	if interactive {
		// Set up terminal for interactive mode
		if err := session.setupInteractiveTerminal(); err != nil {
			session.Close()
			return nil, fmt.Errorf("failed to setup interactive terminal: %w", err)
		}
	}

	// Start handling WebSocket messages
	go session.handleWebSocket(ctx)

	return session, nil
}

// StartShellSession starts a new shell session
func (m *Manager) StartShellSession(sessionID string, websocketURL string) (*Session, error) {
	return m.StartExecSession(sessionID, websocketURL, true)
}

// GetSession returns a session by ID
func (m *Manager) GetSession(sessionID string) (*Session, bool) {
	session, exists := m.sessions[sessionID]
	return session, exists
}

// CloseSession closes a session
func (m *Manager) CloseSession(sessionID string) error {
	session, exists := m.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	session.Close()
	delete(m.sessions, sessionID)
	return nil
}

// CloseAllSessions closes all active sessions
func (m *Manager) CloseAllSessions() {
	for sessionID, session := range m.sessions {
		session.Close()
		delete(m.sessions, sessionID)
	}
}

// setupInteractiveTerminal sets up the terminal for interactive mode
func (s *Session) setupInteractiveTerminal() error {
	// Check if we're in a terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("not running in a terminal")
	}

	// Get original terminal state
	originalState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to make terminal raw: %w", err)
	}
	s.originalState = originalState

	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		// Default size if we can't get it
		width, height = 80, 24
	}

	// Send initial resize message
	resizeMsg := models.ResizeMessage{
		Type: "resize",
		Cols: width,
		Rows: height,
	}

	if err := s.conn.WriteJSON(resizeMsg); err != nil {
		return fmt.Errorf("failed to send resize message: %w", err)
	}

	return nil
}

// handleWebSocket handles WebSocket messages
func (s *Session) handleWebSocket(ctx context.Context) {
	defer s.Close()

	// Start reading from stdin in a separate goroutine for interactive sessions
	if s.isInteractive {
		go s.handleStdin(ctx)
	}

	// Set up signal handling for terminal resize
	if s.isInteractive {
		go s.handleSignals(ctx)
	}

	// Read messages from WebSocket
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg models.ExecMessage
			if err := s.conn.ReadJSON(&msg); err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return
				}
				fmt.Printf("\n❌ WebSocket error: %v\n", err)
				return
			}

			switch msg.Type {
			case "stdout":
				s.handleStdout(msg.Data)
			case "stderr":
				s.handleStderr(msg.Data)
			case "exit":
				s.handleExit(msg.ExitCode)
				return
			}
		}
	}
}

// handleStdin reads from stdin and sends to WebSocket
func (s *Session) handleStdin(ctx context.Context) {
	buffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := os.Stdin.Read(buffer)
			if err != nil {
				return
			}

			if n > 0 {
				// Encode data as base64 for WebSocket transmission
				data := base64.StdEncoding.EncodeToString(buffer[:n])

				msg := models.ExecMessage{
					Type:      "stdin",
					Data:      data,
					Timestamp: time.Now().Format(time.RFC3339),
				}

				if err := s.conn.WriteJSON(msg); err != nil {
					return
				}
			}
		}
	}
}

// handleStdout handles stdout messages from WebSocket
func (s *Session) handleStdout(data string) {
	// Decode base64 data
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	// Write to stdout
	os.Stdout.Write(decoded)
}

// handleStderr handles stderr messages from WebSocket
func (s *Session) handleStderr(data string) {
	// Decode base64 data
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	// Write to stderr
	os.Stderr.Write(decoded)
}

// handleExit handles exit messages from WebSocket
func (s *Session) handleExit(exitCode int) {
	if s.isInteractive {
		fmt.Printf("\n✅ Session ended with exit code: %d\n", exitCode)
	}
}

// handleSignals handles terminal resize signals
func (s *Session) handleSignals(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			// Get new terminal size
			width, height, err := term.GetSize(int(os.Stdin.Fd()))
			if err != nil {
				continue
			}

			// Send resize message
			resizeMsg := models.ResizeMessage{
				Type: "resize",
				Cols: width,
				Rows: height,
			}

			s.conn.WriteJSON(resizeMsg)
		}
	}
}

// SendCommand sends a command to the session
func (s *Session) SendCommand(command string) error {
	data := base64.StdEncoding.EncodeToString([]byte(command))

	msg := models.ExecMessage{
		Type:      "stdin",
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return s.conn.WriteJSON(msg)
}

// Close closes the session
func (s *Session) Close() {
	// Cancel context
	if s.cancel != nil {
		s.cancel()
	}

	// Restore terminal state
	if s.originalState != nil {
		term.Restore(int(os.Stdin.Fd()), s.originalState)
	}

	// Close WebSocket connection
	if s.conn != nil {
		s.conn.Close()
	}
}

// WaitForCompletion waits for the session to complete
func (s *Session) WaitForCompletion() {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\n🛑 Session interrupted by user.")
		s.Close()
	}
}

// ExecCommand executes a single command (non-interactive)
func (m *Manager) ExecCommand(execID string, websocketURL string, command []string) error {
	session, err := m.StartExecSession(execID, websocketURL, false)
	if err != nil {
		return err
	}
	defer session.Close()

	// Send command
	commandStr := ""
	for i, cmd := range command {
		if i > 0 {
			commandStr += " "
		}
		commandStr += cmd
	}
	commandStr += "\n"

	if err := session.SendCommand(commandStr); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	// Wait for completion
	session.WaitForCompletion()

	return nil
}

// GetTerminalSize returns the current terminal size
func GetTerminalSize() (int, int, error) {
	return term.GetSize(int(os.Stdin.Fd()))
}
