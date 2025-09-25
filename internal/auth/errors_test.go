package auth

import (
	"testing"
)

func TestAuthErrorDetection(t *testing.T) {
	tests := []struct {
		name                string
		err                 error
		expectExpired       bool
		expectRevoked       bool
		expectInvalid       bool
		expectRefreshFailed bool
		expectType          string
	}{
		{
			name:                "Token Expired",
			err:                 NewAuthError("token_expired", "Your session has expired", 401, nil),
			expectExpired:       true,
			expectRevoked:       false,
			expectInvalid:       false,
			expectRefreshFailed: false,
			expectType:          "token_expired",
		},
		{
			name:                "Token Revoked",
			err:                 NewAuthError("token_revoked", "Your session has been revoked", 401, nil),
			expectExpired:       false,
			expectRevoked:       true,
			expectInvalid:       false,
			expectRefreshFailed: false,
			expectType:          "token_revoked",
		},
		{
			name:                "Token Invalid",
			err:                 NewAuthError("token_invalid", "Your token is invalid", 401, nil),
			expectExpired:       false,
			expectRevoked:       false,
			expectInvalid:       true,
			expectRefreshFailed: false,
			expectType:          "token_invalid",
		},
		{
			name:                "Refresh Failed",
			err:                 NewAuthError("refresh_failed", "Failed to refresh token", 401, nil),
			expectExpired:       false,
			expectRevoked:       false,
			expectInvalid:       false,
			expectRefreshFailed: true,
			expectType:          "refresh_failed",
		},
		{
			name:                "Generic Auth Error",
			err:                 NewAuthError("authentication_failed", "Authentication failed", 401, nil),
			expectExpired:       false,
			expectRevoked:       false,
			expectInvalid:       false,
			expectRefreshFailed: false,
			expectType:          "authentication_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test error type detection
			if got := IsTokenExpired(tt.err); got != tt.expectExpired {
				t.Errorf("IsTokenExpired() = %v, want %v", got, tt.expectExpired)
			}

			if got := IsTokenRevoked(tt.err); got != tt.expectRevoked {
				t.Errorf("IsTokenRevoked() = %v, want %v", got, tt.expectRevoked)
			}

			if got := IsTokenInvalid(tt.err); got != tt.expectInvalid {
				t.Errorf("IsTokenInvalid() = %v, want %v", got, tt.expectInvalid)
			}

			if got := IsRefreshFailed(tt.err); got != tt.expectRefreshFailed {
				t.Errorf("IsRefreshFailed() = %v, want %v", got, tt.expectRefreshFailed)
			}

			if got := GetAuthErrorType(tt.err); got != tt.expectType {
				t.Errorf("GetAuthErrorType() = %v, want %v", got, tt.expectType)
			}

			// Test that user-friendly message is not empty
			userMsg := GetUserFriendlyMessage(tt.err)
			if userMsg == "" {
				t.Errorf("GetUserFriendlyMessage() returned empty string")
			}
		})
	}
}

func TestAuthErrorCreation(t *testing.T) {
	err := NewAuthError("test_type", "test message", 401, nil)

	if err.Type != "test_type" {
		t.Errorf("Expected Type to be 'test_type', got %s", err.Type)
	}

	if err.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %s", err.Message)
	}

	if err.Code != 401 {
		t.Errorf("Expected Code to be 401, got %d", err.Code)
	}

	if err.Error() != "test message" {
		t.Errorf("Expected Error() to return 'test message', got %s", err.Error())
	}
}

func TestAuthErrorWithUnderlyingError(t *testing.T) {
	underlyingErr := ErrTokenExpired
	err := NewAuthError("test_type", "test message", 401, underlyingErr)

	if err.Unwrap() != underlyingErr {
		t.Errorf("Expected Unwrap() to return underlying error")
	}
}
