package test

import (
	"math/rand"
	"time"

	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/value"
)

func MustNewCount(count int) value.Count {
	c, err := value.NewCount(count)
	if err != nil {
		panic(err)
	}

	return c
}

func MustNewName(name string) value.Name {
	n, err := value.NewName(name)
	if err != nil {
		panic(err)
	}

	return n
}

func MustNewURL(url string) value.URL {
	u, err := value.NewURL(url)
	if err != nil {
		panic(err)
	}

	return u
}

func RandomString(n int) string {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return string(buf)
}

func RandomURL(n int) value.URL {
	return MustNewURL("https://" + RandomString(n))
}

func RandomTime() time.Time {
	return timeutil.Now()
}

func RandomName(n int) value.Name {
	return MustNewName(RandomString(n))
}
