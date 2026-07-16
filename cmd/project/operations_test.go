package project

import (
	"testing"

	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func TestMaskSecretValue(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"ab", "****"},
		{"abcd", "****"},
		{"secret-value", "****ue"},
		{"super-long-password-123", "****23"},
	}
	for _, tc := range tests {
		if got := maskSecretValue(tc.in); got != tc.want {
			t.Errorf("maskSecretValue(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMaskEnvVariables(t *testing.T) {
	in := []sdk.EnvVariable{
		{Key: "API_KEY", Value: "sk-live-abcdef"},
		{Key: "EMPTY", Value: ""},
	}

	revealed := maskEnvVariables(in, true)
	if revealed[0].Value != "sk-live-abcdef" {
		t.Fatalf("reveal=true should keep plaintext, got %q", revealed[0].Value)
	}

	masked := maskEnvVariables(in, false)
	if masked[0].Value == in[0].Value {
		t.Fatal("reveal=false should not return plaintext values")
	}
	if masked[0].Key != "API_KEY" {
		t.Fatalf("keys must be preserved, got %q", masked[0].Key)
	}
	if masked[1].Value != "" {
		t.Fatalf("empty values stay empty, got %q", masked[1].Value)
	}
}
