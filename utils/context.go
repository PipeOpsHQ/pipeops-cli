package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ProjectContext represents the context of a linked project
type ProjectContext struct {
	ProjectID   string    `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Directory   string    `json:"directory"`
	LinkedAt    time.Time `json:"linked_at"`
}

// SaveProjectContext saves project context to .pipeops/project.json
func SaveProjectContext(context *ProjectContext) error {
	context.LinkedAt = time.Now()

	// Create .pipeops directory
	pipeopsDir := filepath.Join(context.Directory, ".pipeops")
	if err := os.MkdirAll(pipeopsDir, 0755); err != nil {
		return fmt.Errorf("failed to create .pipeops directory: %w", err)
	}

	// Save JSON context file
	contextFile := filepath.Join(pipeopsDir, "project.json")
	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project context: %w", err)
	}

	if err := os.WriteFile(contextFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write project context: %w", err)
	}

	// Also create/update legacy .pipeops file for backward compatibility
	legacyFile := filepath.Join(context.Directory, ".pipeops")
	legacyContent := fmt.Sprintf("project_id=%s\n", context.ProjectID)
	if err := os.WriteFile(legacyFile, []byte(legacyContent), 0644); err != nil {
		// Don't fail if legacy file can't be written
		fmt.Printf("Warning: Could not write legacy .pipeops file: %v\n", err)
	}

	return nil
}

// LoadProjectContext loads project context from .pipeops/project.json
func LoadProjectContext() (*ProjectContext, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Look for .pipeops/project.json in current directory and parent directories
	for {
		contextFile := filepath.Join(currentDir, ".pipeops", "project.json")

		if _, err := os.Stat(contextFile); err == nil {
			// File exists, read context
			data, err := os.ReadFile(contextFile)
			if err != nil {
				return nil, fmt.Errorf("error reading project context: %w", err)
			}

			var context ProjectContext
			if err := json.Unmarshal(data, &context); err != nil {
				return nil, fmt.Errorf("error parsing project context: %w", err)
			}

			return &context, nil
		}

		// Move to parent directory
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached root directory
			break
		}
		currentDir = parent
	}

	return nil, fmt.Errorf("no project context found")
}

// GetLinkedProject returns the project ID linked to the current directory
func GetLinkedProject() (string, error) {
	// Try new context format first
	if context, err := LoadProjectContext(); err == nil {
		return context.ProjectID, nil
	}

	// Fall back to legacy format
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
