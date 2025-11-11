package auth

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateSecret(t *testing.T) {
	t.Run("generates secret of correct length", func(t *testing.T) {
		secret, err := GenerateSecret()
		if err != nil {
			t.Fatalf("GenerateSecret() returned error: %v", err)
		}

		expectedLength := 64
		if len(secret) != expectedLength {
			t.Errorf("GenerateSecret() length = %d, want %d", len(secret), expectedLength)
		}
	})

	t.Run("generates hex-encoded secret", func(t *testing.T) {
		secret, err := GenerateSecret()
		if err != nil {
			t.Fatalf("GenerateSecret() returned error: %v", err)
		}

		// Try to decode as hex - should not error
		_, err = hex.DecodeString(secret)
		if err != nil {
			t.Errorf("GenerateSecret() did not produce valid hex: %v", err)
		}

		// Check that secret only contains hex characters (0-9, a-f)
		for i, c := range secret {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("GenerateSecret() contains non-hex character '%c' at position %d", c, i)
			}
		}
	})

	t.Run("generates unique secrets", func(t *testing.T) {
		// Generate multiple secrets and ensure they're different
		secrets := make(map[string]bool)
		iterations := 100

		for i := 0; i < iterations; i++ {
			secret, err := GenerateSecret()
			if err != nil {
				t.Fatalf("GenerateSecret() returned error on iteration %d: %v", i, err)
			}

			if secrets[secret] {
				t.Errorf("GenerateSecret() generated duplicate secret: %s", secret)
			}
			secrets[secret] = true
		}

		if len(secrets) != iterations {
			t.Errorf("GenerateSecret() generated %d unique secrets, want %d", len(secrets), iterations)
		}
	})

	t.Run("generates lowercase hex only", func(t *testing.T) {
		secret, err := GenerateSecret()
		if err != nil {
			t.Fatalf("GenerateSecret() returned error: %v", err)
		}

		if secret != strings.ToLower(secret) {
			t.Errorf("GenerateSecret() contains uppercase characters: %s", secret)
		}
	})
}

func TestValidateSecret(t *testing.T) {
	t.Run("accepts valid secret", func(t *testing.T) {
		// Generate a valid secret
		secret, err := GenerateSecret()
		if err != nil {
			t.Fatalf("GenerateSecret() returned error: %v", err)
		}

		if !ValidateSecret(secret) {
			t.Errorf("ValidateSecret() rejected valid secret: %s", secret)
		}
	})

	t.Run("accepts valid manually created secret", func(t *testing.T) {
		// Create a valid 64-character hex string
		validSecret := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

		if !ValidateSecret(validSecret) {
			t.Errorf("ValidateSecret() rejected valid secret: %s", validSecret)
		}
	})

	t.Run("rejects secret that is too short", func(t *testing.T) {
		shortSecret := "0123456789abcdef" // Only 16 characters

		if ValidateSecret(shortSecret) {
			t.Errorf("ValidateSecret() accepted short secret: %s", shortSecret)
		}
	})

	t.Run("rejects secret that is too long", func(t *testing.T) {
		longSecret := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef00" // 66 characters

		if ValidateSecret(longSecret) {
			t.Errorf("ValidateSecret() accepted long secret: %s", longSecret)
		}
	})

	t.Run("rejects secret with invalid characters", func(t *testing.T) {
		invalidSecrets := []string{
			"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeg", // 'g' is not hex
			"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeG", // uppercase not standard
			"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcd-f", // dash
			"0123456789abcdef 0123456789abcdef0123456789abcdef0123456789abcde", // space
			"0123456789abcdef!123456789abcdef0123456789abcdef0123456789abcdef", // special char
		}

		for _, secret := range invalidSecrets {
			if ValidateSecret(secret) {
				t.Errorf("ValidateSecret() accepted invalid secret: %s", secret)
			}
		}
	})

	t.Run("rejects empty string", func(t *testing.T) {
		if ValidateSecret("") {
			t.Error("ValidateSecret() accepted empty string")
		}
	})

	t.Run("rejects non-hex string of correct length", func(t *testing.T) {
		// 64 characters but not hex
		invalidSecret := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"

		if ValidateSecret(invalidSecret) {
			t.Errorf("ValidateSecret() accepted non-hex secret: %s", invalidSecret)
		}
	})

	t.Run("accepts secrets with all valid hex characters", func(t *testing.T) {
		// Test with all lowercase hex digits
		allDigits := "0000111122223333444455556666777788889999aaaabbbbccccddddeeeeffff"
		if !ValidateSecret(allDigits) {
			t.Errorf("ValidateSecret() rejected valid secret with all hex digits: %s", allDigits)
		}
	})
}

func BenchmarkGenerateSecret(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateSecret()
		if err != nil {
			b.Fatalf("GenerateSecret() error: %v", err)
		}
	}
}

func BenchmarkValidateSecret(b *testing.B) {
	secret, err := GenerateSecret()
	if err != nil {
		b.Fatalf("GenerateSecret() error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateSecret(secret)
	}
}
