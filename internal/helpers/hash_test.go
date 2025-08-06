//go:build unit

/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "simple string",
			input:    "hello",
			expected: 907060870,
		},
		{
			name:     "another string",
			input:    "world",
			expected: 980881731,
		},
		{
			name:     "numeric string",
			input:    "12345",
			expected: 3421846044,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := String(tt.input)
			if result != tt.expected {
				t.Errorf("String(%q) = %d, want %d", tt.input, result, tt.expected)
			}
			// Ensure result is always non-negative
			if result < 0 {
				t.Errorf("String(%q) = %d, want non-negative value", tt.input, result)
			}
		})
	}
}

func TestString_Consistency(t *testing.T) {
	// Test that the same input always produces the same output
	input := "test string"
	first := String(input)
	for i := 0; i < 10; i++ {
		result := String(input)
		if result != first {
			t.Errorf("String(%q) inconsistent: first=%d, iteration %d=%d", input, first, i, result)
		}
	}
}

func TestDjangoPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!@#$%^&*()",
		},
		{
			name:     "unicode password",
			password: "пароль123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := DjangoPassword(tt.password)
			if err != nil {
				t.Fatalf("DjangoPassword(%q) error = %v", tt.password, err)
			}

			// Check format: pbkdf2_sha256$36000$salt$hash
			parts := strings.Split(hash, "$")
			if len(parts) != 4 {
				t.Errorf("DjangoPassword(%q) = %q, want 4 parts separated by $", tt.password, hash)
			}

			if parts[0] != "pbkdf2_sha256" {
				t.Errorf("DjangoPassword(%q) prefix = %q, want pbkdf2_sha256", tt.password, parts[0])
			}

			if parts[1] != "36000" {
				t.Errorf("DjangoPassword(%q) rounds = %q, want 36000", tt.password, parts[1])
			}

			if len(parts[2]) != 12 {
				t.Errorf("DjangoPassword(%q) salt length = %d, want 12", tt.password, len(parts[2]))
			}

			// Verify salt is alphanumeric
			saltPattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
			if !saltPattern.MatchString(parts[2]) {
				t.Errorf("DjangoPassword(%q) salt = %q, want alphanumeric", tt.password, parts[2])
			}

			// Verify hash is valid base64
			_, err = base64.StdEncoding.DecodeString(parts[3])
			if err != nil {
				t.Errorf("DjangoPassword(%q) hash part is not valid base64: %v", tt.password, err)
			}
		})
	}
}

func TestDjangoPassword_Uniqueness(t *testing.T) {
	// Test that the same password produces different hashes (due to different salts)
	password := "testpassword"
	hash1, err1 := DjangoPassword(password)
	hash2, err2 := DjangoPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("DjangoPassword errors: %v, %v", err1, err2)
	}

	if hash1 == hash2 {
		t.Errorf("DjangoPassword should produce different hashes for same password due to salt")
	}
}

func TestDjangoPassword_Verification(t *testing.T) {
	// Test that we can manually verify the Django password hash
	password := "testpassword"
	hash, err := DjangoPassword(password)
	if err != nil {
		t.Fatalf("DjangoPassword error: %v", err)
	}

	// Parse the hash
	parts := strings.Split(hash, "$")
	salt := parts[2]
	expectedHashB64 := parts[3]

	// Recreate the hash manually
	recreatedHash := pbkdf2.Key([]byte(password), []byte(salt), 36000, sha256.Size, sha256.New)
	recreatedHashB64 := base64.StdEncoding.EncodeToString(recreatedHash)

	if recreatedHashB64 != expectedHashB64 {
		t.Errorf("Manual verification failed: expected %q, got %q", expectedHashB64, recreatedHashB64)
	}
}

func TestGenerateRandomAlphanumeric(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "zero length",
			length: 0,
		},
		{
			name:   "small length",
			length: 5,
		},
		{
			name:   "medium length",
			length: 12,
		},
		{
			name:   "large length",
			length: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateRandomAlphanumeric(tt.length)
			if err != nil {
				t.Fatalf("GenerateRandomAlphanumeric(%d) error = %v", tt.length, err)
			}

			if len(result) != tt.length {
				t.Errorf("GenerateRandomAlphanumeric(%d) length = %d, want %d", tt.length, len(result), tt.length)
			}

			// Verify all characters are alphanumeric
			alphanumericPattern := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
			if !alphanumericPattern.MatchString(result) {
				t.Errorf("GenerateRandomAlphanumeric(%d) = %q, contains non-alphanumeric characters", tt.length, result)
			}
		})
	}
}

func TestGenerateRandomAlphanumeric_Uniqueness(t *testing.T) {
	// Test that multiple calls produce different results
	length := 16
	results := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		result, err := GenerateRandomAlphanumeric(length)
		if err != nil {
			t.Fatalf("GenerateRandomAlphanumeric(%d) error = %v", length, err)
		}

		if results[result] {
			t.Errorf("GenerateRandomAlphanumeric(%d) produced duplicate result: %q", length, result)
		}
		results[result] = true
	}
}

func TestSha512Crypt(t *testing.T) {
	tests := []string{
		"password123",
		"",
		"P@ssw0rd!@#$%^&*()",
		"пароль123",
	}

	for _, password := range tests {
		t.Run(fmt.Sprintf("password_%q", password), func(t *testing.T) {
			hash, err := Sha512Crypt(password)
			if err != nil {
				t.Fatalf("Sha512Crypt(%q) error = %v", password, err)
			}

			// Verify hash format: $6$rounds=5000$salt$hash
			if !strings.HasPrefix(hash, "$6$rounds=5000$") {
				t.Errorf("Sha512Crypt(%q) = %q, want prefix $6$rounds=5000$", password, hash)
			}

			// Verify we can verify the password
			err = Sha512CryptVerify(hash, password)
			if err != nil {
				t.Errorf("Sha512CryptVerify failed for password %q: %v", password, err)
			}
		})
	}
}

func TestSha512Crypt_Uniqueness(t *testing.T) {
	// Test that the same password produces different hashes (due to different salts)
	password := "testpassword"
	hash1, err1 := Sha512Crypt(password)
	hash2, err2 := Sha512Crypt(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("Sha512Crypt errors: %v, %v", err1, err2)
	}

	if hash1 == hash2 {
		t.Errorf("Sha512Crypt should produce different hashes for same password due to salt")
	}
}

func TestSha512CryptWithSalt(t *testing.T) {
	tests := []struct {
		name     string
		password string
		salt     string
		rounds   int
		wantErr  bool
	}{
		{
			name:     "valid input",
			password: "password123",
			salt:     "1234567890abcdef",
			rounds:   5000,
			wantErr:  false,
		},
		{
			name:     "valid high rounds",
			password: "password123",
			salt:     "abcdefghijklmnop",
			rounds:   10000,
			wantErr:  false,
		},
		{
			name:     "invalid salt length - too short",
			password: "password123",
			salt:     "short",
			rounds:   5000,
			wantErr:  true,
		},
		{
			name:     "invalid salt length - too long",
			password: "password123",
			salt:     "toolongsaltstring",
			rounds:   5000,
			wantErr:  true,
		},
		{
			name:     "invalid rounds - too low",
			password: "password123",
			salt:     "1234567890abcdef",
			rounds:   4999,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := Sha512CryptWithSalt(tt.password, tt.salt, tt.rounds)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Sha512CryptWithSalt() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Sha512CryptWithSalt() error = %v", err)
			}

			// Verify hash format contains the expected rounds and salt
			expectedPrefix := fmt.Sprintf("$6$rounds=%d$%s$", tt.rounds, tt.salt)
			if !strings.HasPrefix(hash, expectedPrefix) {
				t.Errorf("Sha512CryptWithSalt() = %q, want prefix %q", hash, expectedPrefix)
			}

			// Verify we can verify the password
			err = Sha512CryptVerify(hash, tt.password)
			if err != nil {
				t.Errorf("Sha512CryptVerify failed: %v", err)
			}
		})
	}
}

func TestSha512CryptVerify(t *testing.T) {
	password := "testpassword"
	hash, err := Sha512Crypt(password)
	if err != nil {
		t.Fatalf("Sha512Crypt error: %v", err)
	}

	tests := []struct {
		name     string
		hash     string
		password string
		wantErr  bool
	}{
		{
			name:     "correct password",
			hash:     hash,
			password: password,
			wantErr:  false,
		},
		{
			name:     "wrong password",
			hash:     hash,
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "invalid hash format",
			hash:     "invalid-hash",
			password: password,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Sha512CryptVerify(tt.hash, tt.password)
			if tt.wantErr && err == nil {
				t.Errorf("Sha512CryptVerify() expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Sha512CryptVerify() error = %v", err)
			}
		})
	}
}
