package validation

import (
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"valid name", "My Project", false},
		{"valid name with hyphens", "My-Project", false},
		{"valid name with underscores", "My_Project", false},
		{"invalid name with @", "invalid@name", true},
		{"invalid name with special chars", "invalid!name", true},
		{"empty name", "", true},
		{"too long name", "a" + string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input '%s', but got none", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for input '%s', but got: %v", tt.input, err)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"valid token", "validtoken123456", false},
		{"short token", "abc", true},
		{"token with space", "invalid token", true},
		{"empty token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input '%s', but got none", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for input '%s', but got: %v", tt.input, err)
			}
		})
	}
}
