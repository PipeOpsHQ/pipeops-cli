//go:build windows

package utils

// IsRoot checks if the current process is running with root privileges
func IsRoot() bool {
	// On Windows, checking for administrator privileges is more complex
	// and typically requires checking token elevation.
	// For the context of this CLI tool's usage of bash/sudo, returning false
	// is a safe default as we don't support sudo on Windows in the same way.
	return false
}
