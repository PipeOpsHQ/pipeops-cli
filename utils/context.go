package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetLinkedProject returns the project ID linked to the current directory
func GetLinkedProject() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Look for .pipeops file in current directory and parent directories
	for {
		pipeopsFile := filepath.Join(currentDir, ".pipeops")

		if _, err := os.Stat(pipeopsFile); err == nil {
			// File exists, read project ID
			projectID, err := readProjectIDFromFile(pipeopsFile)
			if err != nil {
				return "", fmt.Errorf("error reading .pipeops file: %w", err)
			}
			return projectID, nil
		}

		// Move to parent directory
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached root directory
			break
		}
		currentDir = parent
	}

	return "", fmt.Errorf("no linked project found in current directory or parent directories")
}

// readProjectIDFromFile reads the project ID from a .pipeops file
func readProjectIDFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "project_id=") {
			return strings.TrimPrefix(line, "project_id="), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("project_id not found in .pipeops file")
}

// IsLinkedProject checks if current directory has a linked project
func IsLinkedProject() bool {
	_, err := GetLinkedProject()
	return err == nil
}

// GetProjectIDOrLinked returns the provided project ID or the linked project ID
// If projectID is provided, it uses that. Otherwise, it tries to get the linked project.
func GetProjectIDOrLinked(projectID string) (string, error) {
	if projectID != "" {
		return projectID, nil
	}

	// Try to get linked project
	linkedID, err := GetLinkedProject()
	if err != nil {
		return "", fmt.Errorf("no project ID provided and no linked project found. Use 'pipeops link <project-id>' to link a project to this directory")
	}

	return linkedID, nil
}

// PrintProjectContext prints information about the current project context
func PrintProjectContext(projectID string) {
	if IsLinkedProject() {
		if linkedID, err := GetLinkedProject(); err == nil && linkedID == projectID {
			fmt.Printf("📂 Using linked project: %s\n", projectID)
		} else {
			fmt.Printf("🎯 Using project: %s (overriding linked project)\n", projectID)
		}
	} else {
		fmt.Printf("🎯 Using project: %s\n", projectID)
	}
}
