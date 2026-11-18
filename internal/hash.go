package internal

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	MEMORY      = 64 * 1024 // 64 KiB
	ITERATIONS  = 3
	PARALLELISM = 2
	SALT_LENGTH = 16
	KEY_LENGTH  = 32
)

func HashPassword(pass string) (string, error) {
	salt := make([]byte, SALT_LENGTH)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(pass), salt, ITERATIONS, MEMORY, PARALLELISM, KEY_LENGTH)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s$%s", encodedSalt, encodedHash), nil
}

func VerifyPassword(pass, storedHash string) (bool, error) {
	parts := strings.Split(storedHash, "$")
	if len(parts) != 2 {
		return false, ErrInvalidHashFormat
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	passHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	finalhash := argon2.IDKey([]byte(pass), salt, ITERATIONS, MEMORY, PARALLELISM, uint32(len(passHash)))

	// avoiding being prone to timing attacks.
	return subtle.ConstantTimeCompare(finalhash, passHash) == 1, nil
}
