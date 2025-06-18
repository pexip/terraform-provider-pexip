package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	return 0
}

// DjangoPassword generates a Django-compatible PBKDF2 password hash
// Compatible with Django's default PBKDF2PasswordHasher
func DjangoPassword(password string) (string, error) {
	const (
		rounds     = 36000
		prefix     = "pbkdf2_sha256"
		saltLength = 12
	)

	// Generate random salt
	salt, err := generateRandomAlphanumeric(saltLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate PBKDF2 hash
	hash := pbkdf2.Key([]byte(password), []byte(salt), rounds, sha256.Size, sha256.New)

	// Encode to base64
	hashB64 := base64.StdEncoding.EncodeToString(hash)

	// Format as Django password hash
	return fmt.Sprintf("%s$%d$%s$%s", prefix, rounds, salt, hashB64), nil
}

// generateRandomAlphanumeric generates a random alphanumeric string of the specified length
func generateRandomAlphanumeric(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
