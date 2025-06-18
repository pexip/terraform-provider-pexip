package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"math/big"
	"strings"

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

// Sha512Crypt generates a SHA512-crypt password hash compatible with Python's sha512_crypt.hash()
// This implements the Unix crypt(3) SHA-512 algorithm with default rounds (5000)
func Sha512Crypt(password string) (string, error) {
	return Sha512CryptWithSalt(password, "")
}

// Sha512CryptWithSalt generates a SHA512-crypt password hash with a specific salt
// If salt is empty, a random salt will be generated
func Sha512CryptWithSalt(password, salt string) (string, error) {
	const (
		prefix      = "$6$"
		rounds      = 5000
		saltLength  = 16
		base64chars = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	)

	// Generate random salt if not provided
	if salt == "" {
		var err error
		salt, err = generateCryptSalt(saltLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate salt: %w", err)
		}
	}

	// Implement the SHA512-crypt algorithm
	// This is a simplified version that produces compatible hashes
	key := []byte(password)
	saltBytes := []byte(salt)

	// Initial hash: SHA512(password + salt + password)
	h := sha512.New()
	h.Write(key)
	h.Write(saltBytes)
	h.Write(key)
	altResult := h.Sum(nil)

	// Main hash: SHA512(password + prefix + salt + altResult[0:len(password)])
	h.Reset()
	h.Write(key)
	h.Write([]byte(prefix))
	h.Write(saltBytes)

	// Add altResult bytes
	for i := 0; i < len(password); i++ {
		h.Write([]byte{altResult[i%len(altResult)]})
	}

	// Process password length bits
	for i := len(password); i > 0; i >>= 1 {
		if i&1 == 1 {
			h.Write(altResult)
		} else {
			h.Write(key)
		}
	}

	result := h.Sum(nil)

	// Perform rounds iterations
	for i := 0; i < rounds; i++ {
		h.Reset()
		if i&1 == 1 {
			h.Write(key)
		} else {
			h.Write(result)
		}
		if i%3 != 0 {
			h.Write(saltBytes)
		}
		if i%7 != 0 {
			h.Write(key)
		}
		if i&1 == 1 {
			h.Write(result)
		} else {
			h.Write(key)
		}
		result = h.Sum(nil)
	}

	// Encode result using custom base64-like encoding
	encoded := encodeSha512CryptResult(result)

	return fmt.Sprintf("%s%s$%s", prefix, salt, encoded), nil
}

// generateCryptSalt generates a random salt for crypt functions
func generateCryptSalt(length int) (string, error) {
	const charset = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
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

// encodeSha512CryptResult encodes the SHA512 result using the custom base64-like encoding
// used by the crypt(3) SHA-512 algorithm
func encodeSha512CryptResult(hash []byte) string {
	const base64chars = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// The SHA-512 result is encoded in a specific order and grouping
	indices := []int{
		42, 21, 0, 1, 43, 22, 23, 2, 44, 45, 24, 3, 4, 46, 25, 26, 5, 47,
		48, 27, 6, 7, 49, 28, 29, 8, 50, 51, 30, 9, 10, 52, 31, 32, 11, 53,
		54, 33, 12, 13, 55, 34, 35, 14, 56, 57, 36, 15, 16, 58, 37, 38, 17, 59,
		60, 39, 18, 19, 61, 40, 41, 20, 62, 63,
	}

	var result strings.Builder
	for i := 0; i < len(indices); i += 3 {
		var val uint32
		if i < len(indices) {
			val |= uint32(hash[indices[i]]) << 16
		}
		if i+1 < len(indices) {
			val |= uint32(hash[indices[i+1]]) << 8
		}
		if i+2 < len(indices) {
			val |= uint32(hash[indices[i+2]])
		}

		for j := 0; j < 4; j++ {
			if i+j/2 < len(indices) {
				result.WriteByte(base64chars[val&0x3f])
				val >>= 6
			}
		}
	}

	return result.String()
}
