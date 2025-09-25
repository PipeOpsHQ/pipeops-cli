package auth

import "errors"

// Authentication error types
var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenRevoked     = errors.New("token has been revoked")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrRefreshFailed    = errors.New("token refresh failed")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrAuthExpired      = errors.New("authentication expired")
)

// AuthError represents an authentication error with additional context
type AuthError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AuthError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "authentication error"
}

// Unwrap returns the underlying error
func (e *AuthError) Unwrap() error {
	return e.Err
}

// NewAuthError creates a new authentication error
func NewAuthError(errType string, message string, code int, err error) *AuthError {
	return &AuthError{
		Type:    errType,
		Message: message,
		Code:    code,
		Err:     err,
	}
}

// IsTokenExpired checks if the error indicates token expiration
func IsTokenExpired(err error) bool {
	if err == nil {
		return false
	}

	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type == "token_expired"
	}

	return errors.Is(err, ErrTokenExpired) ||
		errors.Is(err, ErrAuthExpired)
}

// IsTokenRevoked checks if the error indicates token revocation
func IsTokenRevoked(err error) bool {
	if err == nil {
		return false
	}

	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type == "token_revoked"
	}

	return errors.Is(err, ErrTokenRevoked)
}

// IsTokenInvalid checks if the error indicates an invalid token
func IsTokenInvalid(err error) bool {
	if err == nil {
		return false
	}

	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type == "token_invalid"
	}

	return errors.Is(err, ErrTokenInvalid) ||
		errors.Is(err, ErrTokenMalformed)
}

// IsRefreshFailed checks if the error indicates refresh failure
func IsRefreshFailed(err error) bool {
	if err == nil {
		return false
	}

	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type == "refresh_failed"
	}

	return errors.Is(err, ErrRefreshFailed)
}

// GetAuthErrorType returns the type of authentication error
func GetAuthErrorType(err error) string {
	if err == nil {
		return ""
	}

	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type
	}

	// Check for specific error types
	if IsTokenExpired(err) {
		return "token_expired"
	}
	if IsTokenRevoked(err) {
		return "token_revoked"
	}
	if IsTokenInvalid(err) {
		return "token_invalid"
	}
	if IsRefreshFailed(err) {
		return "refresh_failed"
	}

	return "unknown"
}

// GetUserFriendlyMessage returns a user-friendly error message
func GetUserFriendlyMessage(err error) string {
	if err == nil {
		return ""
	}

	var authErr *AuthError
	if errors.As(err, &authErr) && authErr.Message != "" {
		return authErr.Message
	}

	// Provide user-friendly messages based on error type
	if IsTokenExpired(err) {
		return "Your session has expired. Please run 'pipeops auth login' to authenticate again."
	}
	if IsTokenRevoked(err) {
		return "Your session has been revoked. Please run 'pipeops auth login' to authenticate again."
	}
	if IsTokenInvalid(err) {
		return "Your authentication token is invalid. Please run 'pipeops auth login' to authenticate again."
	}
	if IsRefreshFailed(err) {
		return "Failed to refresh your session. Please run 'pipeops auth login' to authenticate again."
	}

	// Default message
	return "Authentication failed. Please run 'pipeops auth login' to authenticate again."
}
