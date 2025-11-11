package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateSecret generates a cryptographically secure random secret
// Returns a 64-character hex-encoded string (32 bytes)
func GenerateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secret: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateSecret validates that a secret is in the correct format
// A valid secret must be exactly 64 characters and hex-encoded
func ValidateSecret(secret string) bool {
	if len(secret) != 64 {
		return false
	}
	_, err := hex.DecodeString(secret)
	return err == nil
}
