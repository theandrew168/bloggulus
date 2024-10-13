package random

import (
	"crypto/rand"
	"encoding/base64"
)

func BytesBase64(n int) (string, error) {
	buf := make([]byte, n)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
