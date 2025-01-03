package k3s

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install k3s and connect to PipeOps",
	Long:  `Installs the k3s server and connects it to the PipeOps control plane using your service account token.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := installK3s(); err != nil {
			log.Fatalf("Error installing k3s: %v", err)
		}
	},
	Args: cobra.NoArgs,
}

func installK3s() error {
    log.Println("Detecting system configuration...")
    
    // Check if running as root
    if output, err := utils.RunCommand("id", "-u"); err != nil || strings.TrimSpace(output) != "0" {
        return fmt.Errorf("k3s installation requires root privileges. Please run with sudo")
    }

    // Get system information
    os := runtime.GOOS
    arch := runtime.GOARCH

    // Map architecture names to k3s nomenclature
    archMap := map[string]string{
        "amd64": "amd64",
        "arm64": "arm64",
        "arm":   "armhf",
    }

    // Check if system is supported
    if os != "linux" {
        return fmt.Errorf("k3s is only supported on Linux systems. Current OS: %s", os)
    }

    k3sArch, ok := archMap[arch]
    if !ok {
        return fmt.Errorf("unsupported architecture: %s", arch)
    }

    log.Printf("System detected: OS=%s, Architecture=%s", os, k3sArch)

    // Detect init system
    initSystem, err := detectInitSystem()
    if err != nil {
        return fmt.Errorf("failed to detect init system: %v", err)
    }

    log.Printf("Detected init system: %s", initSystem)

    // Choose installation method based on init system
    var installCmd string
    switch initSystem {
    case "systemd":
        installCmd = fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_ARCH=%s sh -", k3sArch)
    case "openrc":
        installCmd = fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_ARCH=%s INSTALL_K3S_EXEC='server --disable-agent' sh -", k3sArch)
    default:
        // For systems without systemd/openrc, use manual installation
        return installK3sManually(k3sArch)
    }

    log.Println("Installing k3s...")
    output, err := utils.RunCommand("sh", "-c", installCmd)
    if err != nil {
        return fmt.Errorf("installation failed: %v\nOutput: %s", err, output)
    }

    log.Println("k3s installed successfully")
    return nil
}

func detectInitSystem() (string, error) {
    // Check for systemd
    if _, err := utils.RunCommand("which", "systemctl"); err == nil {
        return "systemd", nil
    }

    // Check for openrc
    if _, err := utils.RunCommand("which", "rc-service"); err == nil {
        return "openrc", nil
    }

    return "other", nil
}

func installK3sManually(arch string) error {
    log.Println("Performing manual installation for system without systemd/openrc...")
    
    // Download k3s binary
    binaryURL := fmt.Sprintf("https://github.com/k3s-io/k3s/releases/latest/download/k3s-%s-%s", "linux", arch)
    log.Printf("Downloading k3s binary from: %s", binaryURL)
    
    if _, err := utils.RunCommand("curl", "-fLo", "/usr/local/bin/k3s", binaryURL); err != nil {
        return fmt.Errorf("failed to download k3s binary: %v", err)
    }

    // Make binary executable
    if _, err := utils.RunCommand("chmod", "+x", "/usr/local/bin/k3s"); err != nil {
        return fmt.Errorf("failed to make k3s binary executable: %v", err)
    }

    // Create necessary directories
    dirs := []string{"/etc/rancher/k3s", "/var/lib/rancher/k3s/data"}
    for _, dir := range dirs {
        if _, err := utils.RunCommand("mkdir", "-p", dir); err != nil {
            return fmt.Errorf("failed to create directory %s: %v", dir, err)
        }
    }

    // Start k3s with minimal configuration
    log.Println("Starting k3s server...")
    cmd := exec.Command("/usr/local/bin/k3s", "server",
        "--disable-agent",
        "--data-dir", "/var/lib/rancher/k3s/data",
        "--write-kubeconfig", "/etc/rancher/k3s/k3s.yaml")
    
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start k3s server: %v", err)
    }

    log.Println("K3s started successfully in manual mode")
    return nil
}


func (k *k3sModel) install() {
	k.rootCmd.AddCommand(installCmd)
}
