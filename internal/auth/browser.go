package auth

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenBrowser opens the default browser to the specified URL
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// GetPlatformName returns a human-readable platform name
func GetPlatformName() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return runtime.GOOS
	}
}

// IsBrowserAvailable checks if a browser is available on the system
func IsBrowserAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		return isCommandAvailable("open")
	case "linux":
		return isCommandAvailable("xdg-open")
	case "windows":
		return isCommandAvailable("rundll32")
	default:
		return false
	}
}

// isCommandAvailable checks if a command is available in the system PATH
func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
