package tailscale

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Client represents a Tailscale client for managing VPN connections
type Client struct {
	authKey string
}

// NewClient creates a new Tailscale client
func NewClient() *Client {
	return &Client{}
}

// IsInstalled checks if Tailscale is installed on the system
func (c *Client) IsInstalled() bool {
	_, err := exec.LookPath("tailscale")
	return err == nil
}

// IsConnected checks if Tailscale is connected and active
func (c *Client) IsConnected() (bool, error) {
	cmd := exec.Command("tailscale", "status", "--json")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check tailscale status: %w", err)
	}

	// Simple check - if output contains "BackendState", it's likely connected
	return strings.Contains(string(output), "BackendState"), nil
}

// GetStatus returns the current Tailscale status
func (c *Client) GetStatus() (string, error) {
	if !c.IsInstalled() {
		return "", errors.New("tailscale is not installed")
	}

	cmd := exec.Command("tailscale", "status")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get tailscale status: %w", err)
	}

	return string(output), nil
}

// Connect connects to Tailscale using the provided auth key
func (c *Client) Connect(authKey string) error {
	if !c.IsInstalled() {
		return errors.New("tailscale is not installed")
	}

	if authKey == "" {
		return errors.New("auth key is required")
	}

	c.authKey = authKey

	cmd := exec.Command("tailscale", "up", "--authkey", authKey)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to connect to tailscale: %w, output: %s", err, string(output))
	}

	return nil
}

// Disconnect disconnects from Tailscale
func (c *Client) Disconnect() error {
	if !c.IsInstalled() {
		return errors.New("tailscale is not installed")
	}

	cmd := exec.Command("tailscale", "down")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disconnect from tailscale: %w, output: %s", err, string(output))
	}

	return nil
}

// GetIP returns the Tailscale IP address
func (c *Client) GetIP() (string, error) {
	if !c.IsInstalled() {
		return "", errors.New("tailscale is not installed")
	}

	cmd := exec.Command("tailscale", "ip", "-4")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get tailscale IP: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// ListPeers returns a list of connected peers
func (c *Client) ListPeers() ([]string, error) {
	if !c.IsInstalled() {
		return nil, errors.New("tailscale is not installed")
	}

	cmd := exec.Command("tailscale", "status", "--peers")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list peers: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var peers []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			peers = append(peers, line)
		}
	}

	return peers, nil
}

// Ping pings a peer in the Tailscale network
func (c *Client) Ping(peer string) error {
	if !c.IsInstalled() {
		return errors.New("tailscale is not installed")
	}

	if peer == "" {
		return errors.New("peer address is required")
	}

	cmd := exec.Command("tailscale", "ping", peer)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to ping peer %s: %w, output: %s", peer, err, string(output))
	}

	return nil
}
