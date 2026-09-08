package value

import (
	"errors"
	"net/url"
	"strings"
)

type URL struct {
	value string
}

var ErrEmptyURL = errors.New("url: value cannot be empty or whitespace only")

type ErrInvalidURL struct {
	Value string
}

func (e *ErrInvalidURL) Error() string {
	return "url: value is not a valid URL: " + e.Value
}

func NewURL(value string) (URL, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return URL{}, ErrEmptyURL
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return URL{}, &ErrInvalidURL{Value: trimmed}
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
