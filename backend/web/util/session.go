package util

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate a random, crypto-safe session ID.
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
