package pipeops

import (
	"errors"
	"os"
	"testing"
)

func TestShouldPromptForWorkspaceWithMode(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		stdinMode  os.FileMode
		statErr    error
		wantPrompt bool
	}{
		{
			name:       "interactive terminal allows prompt",
			stdinMode:  os.ModeCharDevice,
			wantPrompt: true,
		},
		{
			name:       "json output disables prompt",
			env:        map[string]string{"PIPEOPS_OUTPUT_JSON": "true"},
			stdinMode:  os.ModeCharDevice,
			wantPrompt: false,
		},
		{
			name:       "explicit non-interactive disables prompt",
			env:        map[string]string{"PIPEOPS_NON_INTERACTIVE": "1"},
			stdinMode:  os.ModeCharDevice,
			wantPrompt: false,
		},
		{
			name:       "ci disables prompt",
			env:        map[string]string{"CI": "true"},
			stdinMode:  os.ModeCharDevice,
			wantPrompt: false,
		},
		{
			name:       "piped stdin disables prompt",
			stdinMode:  0,
			wantPrompt: false,
		},
		{
			name:       "stdin stat error disables prompt",
			stdinMode:  os.ModeCharDevice,
			statErr:    errors.New("stat failed"),
			wantPrompt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PIPEOPS_OUTPUT_JSON", "")
			t.Setenv("PIPEOPS_NON_INTERACTIVE", "")
			t.Setenv("CI", "")
			t.Setenv("GITHUB_ACTIONS", "")
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			got := shouldPromptForWorkspaceWithMode(tt.stdinMode, tt.statErr)
			if got != tt.wantPrompt {
				t.Fatalf("shouldPromptForWorkspaceWithMode() = %v, want %v", got, tt.wantPrompt)
			}
		})
	}
}

func TestEnvEnabled(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{value: "true", want: true},
		{value: "1", want: true},
		{value: "yes", want: true},
		{value: "on", want: true},
		{value: "false", want: false},
		{value: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv("PIPEOPS_TEST_FLAG", tt.value)
			got := envEnabled("PIPEOPS_TEST_FLAG")
			if got != tt.want {
				t.Fatalf("envEnabled(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
