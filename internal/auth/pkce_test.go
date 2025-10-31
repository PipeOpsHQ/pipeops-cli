package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGeneratePKCEChallenge(t *testing.T) {
	challenge, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge() error = %v", err)
	}

	if challenge == nil {
		t.Fatal("GeneratePKCEChallenge() returned nil")
	}

	if challenge.CodeVerifier == "" {
		t.Error("CodeVerifier is empty")
	}

	if challenge.CodeChallenge == "" {
		t.Error("CodeChallenge is empty")
	}

	if challenge.Method != "S256" {
		t.Errorf("Method = %v, want S256", challenge.Method)
	}

	// Verify code verifier length (should be 43 characters for base64url of 32 bytes)
	if len(challenge.CodeVerifier) != 43 {
		t.Errorf("CodeVerifier length = %d, want 43", len(challenge.CodeVerifier))
	}

	// Verify code challenge is correctly derived from verifier
	hash := sha256.Sum256([]byte(challenge.CodeVerifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	if challenge.CodeChallenge != expectedChallenge {
		t.Error("CodeChallenge does not match expected SHA256 hash of CodeVerifier")
	}
}

func TestPKCEChallengeVerify(t *testing.T) {
	challenge, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge() error = %v", err)
	}

	// Test valid verification
	if !challenge.Verify(challenge.CodeVerifier) {
		t.Error("Verify() should return true for correct verifier")
	}

	// Test invalid verification
	if challenge.Verify("wrong_verifier") {
		t.Error("Verify() should return false for incorrect verifier")
	}

	if challenge.Verify("") {
		t.Error("Verify() should return false for empty verifier")
	}
}

func TestGenerateRandomState(t *testing.T) {
	state, err := GenerateRandomState()
	if err != nil {
		t.Fatalf("GenerateRandomState() error = %v", err)
	}

	if state == "" {
		t.Error("GenerateRandomState() returned empty string")
	}

	// Generate multiple states to ensure they're different
	states := make(map[string]bool)
	for i := 0; i < 10; i++ {
		s, err := GenerateRandomState()
		if err != nil {
			t.Fatalf("GenerateRandomState() error = %v", err)
		}
		if states[s] {
			t.Error("GenerateRandomState() generated duplicate state")
		}
		states[s] = true
	}
}

func TestPKCEChallengeUniqueness(t *testing.T) {
	// Generate multiple challenges to ensure they're unique
	challenges := make(map[string]bool)
	verifiers := make(map[string]bool)

	for i := 0; i < 10; i++ {
		challenge, err := GeneratePKCEChallenge()
		if err != nil {
			t.Fatalf("GeneratePKCEChallenge() error = %v", err)
		}

		if challenges[challenge.CodeChallenge] {
			t.Error("Generated duplicate code challenge")
		}
		challenges[challenge.CodeChallenge] = true

		if verifiers[challenge.CodeVerifier] {
			t.Error("Generated duplicate code verifier")
		}
		verifiers[challenge.CodeVerifier] = true
	}
}

func TestPKCEChallengeBase64URLEncoding(t *testing.T) {
	challenge, err := GeneratePKCEChallenge()
	if err != nil {
		t.Fatalf("GeneratePKCEChallenge() error = %v", err)
	}

	// Verify that the code verifier can be decoded as base64url
	_, err = base64.RawURLEncoding.DecodeString(challenge.CodeVerifier)
	if err != nil {
		t.Errorf("CodeVerifier is not valid base64url: %v", err)
	}

	// Verify that the code challenge can be decoded as base64url
	_, err = base64.RawURLEncoding.DecodeString(challenge.CodeChallenge)
	if err != nil {
		t.Errorf("CodeChallenge is not valid base64url: %v", err)
	}

	// Verify no padding (RawURLEncoding)
	if len(challenge.CodeVerifier) > 0 && challenge.CodeVerifier[len(challenge.CodeVerifier)-1] == '=' {
		t.Error("CodeVerifier should not have padding")
	}

	if len(challenge.CodeChallenge) > 0 && challenge.CodeChallenge[len(challenge.CodeChallenge)-1] == '=' {
		t.Error("CodeChallenge should not have padding")
	}
}
