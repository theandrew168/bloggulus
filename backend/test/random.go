package test

import (
	"math/rand"
	"time"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func RandomString(n int) string {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return string(buf)
}

func RandomURL(n int) string {
	return "https://" + RandomString(n)
}

func RandomTime() time.Time {
	return timeutil.Now()
}
