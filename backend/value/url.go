package value

import (
	"errors"
	"net/url"
	"strings"
)

type URL struct {
	value string
}

// TODO: Custom error type that includes the invalid URL string in the error message

var ErrEmptyURL = errors.New("url: value cannot be empty or whitespace only")
var ErrInvalidURL = errors.New("url: value is not a valid URL")

func NewURL(value string) (URL, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return URL{}, ErrEmptyURL
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return URL{}, ErrInvalidURL
	}

	u := URL{
		value: parsed.String(),
	}
	return u, nil
}

func MustNewURL(value string) URL {
	u, err := NewURL(value)
	if err != nil {
		panic(err)
	}

	return u
}

func (u URL) Value() string {
	return u.value
}
