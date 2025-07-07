package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// ProjectNameValidator validates project names
type ProjectNameValidator struct {
	MinLength int
	MaxLength int
}

// NewProjectNameValidator creates a new project name validator
func NewProjectNameValidator() *ProjectNameValidator {
	return &ProjectNameValidator{
		MinLength: 1,
		MaxLength: 100,
	}
}

// Validate validates a project name
func (v *ProjectNameValidator) Validate(name string) error {
	name = strings.TrimSpace(name)

	if len(name) < v.MinLength {
		return fmt.Errorf("project name is too short (minimum %d characters)", v.MinLength)
	}

	if len(name) > v.MaxLength {
		return fmt.Errorf("project name is too long (maximum %d characters)", v.MaxLength)
	}

	// Check for valid characters (alphanumeric, spaces, hyphens, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("project name contains invalid characters (only letters, numbers, spaces, hyphens, and underscores are allowed)")
	}

	return nil
}

// ProjectDescriptionValidator validates project descriptions
type ProjectDescriptionValidator struct {
	MaxLength int
}

// NewProjectDescriptionValidator creates a new project description validator
func NewProjectDescriptionValidator() *ProjectDescriptionValidator {
	return &ProjectDescriptionValidator{
		MaxLength: 500,
	}
}

// Validate validates a project description
func (v *ProjectDescriptionValidator) Validate(description string) error {
	description = strings.TrimSpace(description)

	if len(description) > v.MaxLength {
		return fmt.Errorf("project description is too long (maximum %d characters)", v.MaxLength)
	}

	return nil
}

// TokenValidator validates authentication tokens
type TokenValidator struct {
	MinLength int
	MaxLength int
}

// NewTokenValidator creates a new token validator
func NewTokenValidator() *TokenValidator {
	return &TokenValidator{
		MinLength: 10,
		MaxLength: 1000,
	}
}

// Validate validates an authentication token
func (v *TokenValidator) Validate(token string) error {
	token = strings.TrimSpace(token)

	if len(token) < v.MinLength {
		return fmt.Errorf("token is too short (minimum %d characters)", v.MinLength)
	}

	if len(token) > v.MaxLength {
		return fmt.Errorf("token is too long (maximum %d characters)", v.MaxLength)
	}

	// Check for suspicious characters or patterns
	if strings.Contains(token, " ") {
		return fmt.Errorf("token should not contain spaces")
	}

	return nil
}

// ProjectIDValidator validates project IDs
type ProjectIDValidator struct {
	MinLength int
	MaxLength int
}

// NewProjectIDValidator creates a new project ID validator
func NewProjectIDValidator() *ProjectIDValidator {
	return &ProjectIDValidator{
		MinLength: 1,
		MaxLength: 50,
	}
}

// Validate validates a project ID
func (v *ProjectIDValidator) Validate(id string) error {
	id = strings.TrimSpace(id)

	if len(id) < v.MinLength {
		return fmt.Errorf("project ID is too short (minimum %d characters)", v.MinLength)
	}

	if len(id) > v.MaxLength {
		return fmt.Errorf("project ID is too long (maximum %d characters)", v.MaxLength)
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	validID := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	if !validID.MatchString(id) {
		return fmt.Errorf("project ID contains invalid characters (only letters, numbers, hyphens, and underscores are allowed)")
	}

	return nil
}

// ValidateProjectName validates a project name using the default validator
func ValidateProjectName(name string) error {
	validator := NewProjectNameValidator()
	return validator.Validate(name)
}

// ValidateProjectDescription validates a project description using the default validator
func ValidateProjectDescription(description string) error {
	validator := NewProjectDescriptionValidator()
	return validator.Validate(description)
}

// ValidateToken validates an authentication token using the default validator
func ValidateToken(token string) error {
	validator := NewTokenValidator()
	return validator.Validate(token)
}

// ValidateProjectID validates a project ID using the default validator
func ValidateProjectID(id string) error {
	validator := NewProjectIDValidator()
	return validator.Validate(id)
}
