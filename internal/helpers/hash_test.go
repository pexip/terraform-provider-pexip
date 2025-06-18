package helpers

import (
	"strings"
	"testing"
)

func TestDjangoPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "admin",
		},
		{
			name:     "complex password",
			password: "super_secure_password123!",
		},
		{
			name:     "empty password",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := DjangoPassword(tt.password)
			if err != nil {
				t.Fatalf("DjangoPassword() error = %v", err)
			}

			// Verify hash format: pbkdf2_sha256$36000$salt$hash
			parts := strings.Split(hash, "$")
			if len(parts) != 4 {
				t.Errorf("DjangoPassword() hash format incorrect, got %d parts, want 4", len(parts))
			}

			if parts[0] != "pbkdf2_sha256" {
				t.Errorf("DjangoPassword() algorithm = %v, want pbkdf2_sha256", parts[0])
			}

			if parts[1] != "36000" {
				t.Errorf("DjangoPassword() rounds = %v, want 36000", parts[1])
			}

			if len(parts[2]) != 12 {
				t.Errorf("DjangoPassword() salt length = %v, want 12", len(parts[2]))
			}

			if len(parts[3]) == 0 {
				t.Error("DjangoPassword() hash part is empty")
			}

			// Verify salt is alphanumeric
			salt := parts[2]
			for _, char := range salt {
				if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
					t.Errorf("DjangoPassword() salt contains non-alphanumeric character: %c", char)
				}
			}
		})
	}
}

func TestDjangoPasswordUniqueness(t *testing.T) {
	password := "test_password"

	// Generate multiple hashes for the same password
	hash1, err1 := DjangoPassword(password)
	hash2, err2 := DjangoPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("DjangoPassword() errors: %v, %v", err1, err2)
	}

	// Hashes should be different due to random salt
	if hash1 == hash2 {
		t.Error("DjangoPassword() generated identical hashes for same password (should be different due to random salt)")
	}
}

func TestGenerateRandomAlphanumeric(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length 12",
			length: 12,
		},
		{
			name:   "length 1",
			length: 1,
		},
		{
			name:   "length 50",
			length: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generateRandomAlphanumeric(tt.length)
			if err != nil {
				t.Fatalf("generateRandomAlphanumeric() error = %v", err)
			}

			if len(result) != tt.length {
				t.Errorf("generateRandomAlphanumeric() length = %v, want %v", len(result), tt.length)
			}

			// Verify all characters are alphanumeric
			for _, char := range result {
				if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
					t.Errorf("generateRandomAlphanumeric() contains non-alphanumeric character: %c", char)
				}
			}
		})
	}
}

func TestGenerateRandomAlphanumericUniqueness(t *testing.T) {
	length := 12

	// Generate multiple random strings
	str1, err1 := generateRandomAlphanumeric(length)
	str2, err2 := generateRandomAlphanumeric(length)

	if err1 != nil || err2 != nil {
		t.Fatalf("generateRandomAlphanumeric() errors: %v, %v", err1, err2)
	}

	// Strings should be different (extremely unlikely to be the same)
	if str1 == str2 {
		t.Error("generateRandomAlphanumeric() generated identical strings (should be random)")
	}
}
