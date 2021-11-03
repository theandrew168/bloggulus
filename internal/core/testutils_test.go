package core_test

import (
	"math/rand"
	"time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func randomString(n int) string {
	validRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	b := make([]byte, n)
	for i := range b {
		b[i] = validRunes[rand.Intn(len(validRunes))]
	}

	return string(b)
}
