package tailscale

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Client represents a Tailscale client for managing VPN connections and Funnel exposure
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

// EnableFunnel enables Tailscale Funnel for port 80 exposure
func (c *Client) EnableFunnel(port int) error {
	if !c.IsInstalled() {
		return errors.New("tailscale is not installed")
	}

	if port == 0 {
		port = 80 // Default to port 80
	}

	cmd := exec.Command("tailscale", "serve", "funnel", fmt.Sprintf("%d", port))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable funnel on port %d: %w, output: %s", port, err, string(output))
	}

	return nil
}

// DisableFunnel disables Tailscale Funnel
func (c *Client) DisableFunnel() error {
	if !c.IsInstalled() {
		return errors.New("tailscale is not installed")
	}

	cmd := exec.Command("tailscale", "serve", "funnel", "off")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable funnel: %w, output: %s", err, string(output))
	}

	return nil
}

// GetFunnelStatus returns the current Funnel status
func (c *Client) GetFunnelStatus() (string, error) {
	if !c.IsInstalled() {
		return "", errors.New("tailscale is not installed")
	}

	cmd := exec.Command("tailscale", "serve", "status")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get funnel status: %w", err)
	}

	return string(output), nil
}

// GetFunnelURL returns the public URL for the Funnel service
func (c *Client) GetFunnelURL() (string, error) {
	if !c.IsInstalled() {
		return "", errors.New("tailscale is not installed")
	}

	cmd := exec.Command("tailscale", "serve", "status", "--json")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get funnel URL: %w", err)
	}

	// Parse JSON output to extract the public URL
	// This is a simplified implementation - in production you'd want proper JSON parsing
	outputStr := string(output)
	if strings.Contains(outputStr, "funnel") {
		// Extract URL from the JSON response
		// This is a basic implementation - you might want to use proper JSON parsing
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			if strings.Contains(line, "https://") {
				// Extract URL from the line
				start := strings.Index(line, "https://")
				if start != -1 {
					end := strings.Index(line[start:], "\"")
					if end != -1 {
						return line[start : start+end], nil
					}
				}
			}
		}
	}

	return "", errors.New("no funnel URL found")
}

// InstallTailscale installs Tailscale on the system
func (c *Client) InstallTailscale() error {
	// Check if already installed
	if c.IsInstalled() {
		return nil
	}

	// Detect OS and install accordingly
	cmd := exec.Command("sh", "-c", "curl -fsSL https://tailscale.com/install.sh | sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install tailscale: %w, output: %s", err, string(output))
	}

	return nil
}

// SetupKubernetesOperator installs and configures the Tailscale Kubernetes operator
func (c *Client) SetupKubernetesOperator() error {
	if !c.IsInstalled() {
		return errors.New("tailscale is not installed")
	}

	// Install the Tailscale Kubernetes operator
	operatorCmd := `kubectl apply -f https://raw.githubusercontent.com/tailscale/tailscale/main/cmd/k8s-operator/deploy.yaml`
	cmd := exec.Command("sh", "-c", operatorCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install tailscale kubernetes operator: %w, output: %s", err, string(output))
	}

	return nil
}

// CreateIngressWithFunnel creates a Kubernetes ingress with Tailscale Funnel enabled
func (c *Client) CreateIngressWithFunnel(serviceName, hostname string, port int) error {
	if !c.IsInstalled() {
		return errors.New("tailscale is not installed")
	}

	if port == 0 {
		port = 80
	}

	// Create ingress manifest with Tailscale Funnel annotation
	ingressManifest := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: %s-funnel-ingress
  annotations:
    tailscale.com/funnel: "true"
spec:
  ingressClassName: tailscale
  rules:
  - host: %s
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: %s
            port:
              number: %d
`, serviceName, hostname, serviceName, port)

	// Apply the ingress manifest
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(ingressManifest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create ingress with funnel: %w, output: %s", err, string(output))
	}

	return nil
}
