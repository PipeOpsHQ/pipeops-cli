package auth

import "testing"

func TestValidOAuthState(t *testing.T) {
	tests := []struct {
		name     string
		received string
		expected string
		want     bool
	}{
		{
			name:     "exact match",
			received: "LQ6mH8XoS6T2U0QY293Keg",
			expected: "LQ6mH8XoS6T2U0QY293Keg",
			want:     true,
		},
		{
			name:     "provider strips base64 padding",
			received: "LQ6mH8XoS6T2U0QY293Keg",
			expected: "LQ6mH8XoS6T2U0QY293Keg==",
			want:     true,
		},
		{
			name:     "provider preserves base64 padding",
			received: "LQ6mH8XoS6T2U0QY293Keg==",
			expected: "LQ6mH8XoS6T2U0QY293Keg",
			want:     true,
		},
		{
			name:     "different state rejected",
			received: "different",
			expected: "LQ6mH8XoS6T2U0QY293Keg",
			want:     false,
		},
		{
			name:     "empty received rejected",
			received: "",
			expected: "LQ6mH8XoS6T2U0QY293Keg",
			want:     false,
		},
		{
			name:     "empty expected rejected",
			received: "LQ6mH8XoS6T2U0QY293Keg",
			expected: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validOAuthState(tt.received, tt.expected); got != tt.want {
				t.Fatalf("validOAuthState() = %v, want %v", got, tt.want)
			}
		})
	}
}
