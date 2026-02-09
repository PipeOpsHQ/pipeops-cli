//go:build !windows

package utils

import "os"

// IsRoot checks if the current process is running with root privileges
func IsRoot() bool {
	return os.Geteuid() == 0
}
