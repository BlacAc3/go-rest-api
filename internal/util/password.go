package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	hashIteration uint32 = 2
	memoryUsage   uint32 = 64 * 1024 // in kilobytes
	threads       uint8   = 4
	outputLength  uint32    = 32
)

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)

	return salt
}

func HashPassword(password string) string {
	salt := generateSalt()
	hash := argon2.IDKey([]byte(password), salt, hashIteration, memoryUsage, threads, outputLength)
	return fmt.Sprintf(
		"%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

func VerifyPassword(password string, hashed string) bool {
	parts := strings.Split(hashed, "$")
	if len(parts) != 2 {
		return false
	}

	inputPassword := password
    storedSalt, _ := base64.RawStdEncoding.DecodeString(parts[0])
    storedHash, _ := base64.RawStdEncoding.DecodeString(parts[1])

	inputPasswordHash := argon2.IDKey([]byte(inputPassword), storedSalt, hashIteration, memoryUsage, threads, outputLength)
	return string(inputPasswordHash) == string(storedHash)

}
