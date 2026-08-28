package value

import (
	"errors"
	"strings"
)

type Name struct {
	value string
}

var ErrEmptyString = errors.New("name: value cannot be empty or whitespace only")

func NewName(value string) (Name, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Name{}, ErrEmptyString
	}

	n := Name{
		value: trimmed,
	}
	return n, nil
}

func (n Name) Value() string {
	return n.value
}
