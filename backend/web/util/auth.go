package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Hash a user ID with a given key (via HMAC).
func HashUserID(userID, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(userID))
	userIDHash := mac.Sum(nil)
	return hex.EncodeToString(userIDHash)
}
