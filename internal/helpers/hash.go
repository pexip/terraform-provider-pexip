package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"math"
	"math/big"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
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
	const saltLength = 12
	const rounds = 36000

	// Generate random salt
	salt, err := GenerateRandomAlphanumeric(saltLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return DjangoPasswordWithSalt(password, salt, rounds)
}

func DjangoPasswordWithSalt(password, salt string, rounds int) (string, error) {
	const prefix = "pbkdf2_sha256"

	// Generate PBKDF2 hash
	hash := pbkdf2.Key([]byte(password), []byte(salt), rounds, sha256.Size, sha256.New)

	// Encode to base64
	hashB64 := base64.StdEncoding.EncodeToString(hash)

	// Format as Django password hash
	return fmt.Sprintf("%s$%d$%s$%s", prefix, rounds, salt, hashB64), nil
}

// GenerateRandomAlphanumeric generates a random alphanumeric string of the specified length
func GenerateRandomAlphanumeric(length int) (string, error) {
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

// Sha512Crypt generates a SHA512-crypt password hash compatible with Python's sha512_crypt.hash()
// This implements the Unix crypt(3) SHA-512 algorithm with default rounds (5000)
func Sha512Crypt(password string) (string, error) {
	const saltLength = 16
	const rounds = 5000

	salt, err := GenerateRandomAlphanumeric(saltLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return Sha512CryptWithSalt(password, salt, rounds)
}

// Sha512CryptWithSalt generates a SHA512-crypt password hash with a specific salt and rounds
func Sha512CryptWithSalt(password, salt string, rounds int) (string, error) {
	c := crypt.New(crypt.SHA512)
	if len(salt) != 16 {
		return "", fmt.Errorf("salt must be exactly 16 characters long, got %d", len(salt))
	}
	if rounds < 5000 {
		return "", fmt.Errorf("rounds must be at least 5000, got %d", rounds)
	}

	// Format salt for sha512_crypt
	saltConfig := fmt.Sprintf("$6$rounds=%d$%s$", rounds, salt)
	return c.Generate([]byte(password), []byte(saltConfig))
}

// Sha512CryptVerify verifies a SHA512-crypt password hash against a plaintext password
func Sha512CryptVerify(hash, password string) error {
	c := crypt.New(crypt.SHA512)
	return c.Verify(hash, []byte(password))
}

// SafeInt32 safely converts an int to int32, checking for overflow
func SafeInt32(value int) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("integer overflow: value %d cannot be safely converted to int32", value)
	}
	return int32(value), nil
}
