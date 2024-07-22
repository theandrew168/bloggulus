package util

import (
	"net/url"
	"strconv"
)

func ReadInt(qs url.Values, key string, defaultValue int, e *Errors) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		e.AddField("must be an integer", key)
		return defaultValue
	}

	return i
}
