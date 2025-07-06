package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// PKCEChallenge represents a PKCE code challenge and verifier pair
type PKCEChallenge struct {
	CodeVerifier  string
	CodeChallenge string
	Method        string
}

// GeneratePKCEChallenge creates a new PKCE challenge according to RFC 7636
func GeneratePKCEChallenge() (*PKCEChallenge, error) {
	// Generate code verifier: 43-128 characters, URL-safe
	verifierBytes := make([]byte, 32) // 32 bytes = 43 characters when base64url encoded
	if _, err := rand.Read(verifierBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)

	// Generate code challenge: SHA256 hash of verifier, base64url encoded
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return &PKCEChallenge{
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
		Method:        "S256", // SHA256
	}, nil
}

// Verify checks if the given verifier matches this challenge
func (p *PKCEChallenge) Verify(verifier string) bool {
	hash := sha256.Sum256([]byte(verifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return expectedChallenge == p.CodeChallenge
}

// GenerateRandomState generates a random state parameter for CSRF protection
func GenerateRandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
