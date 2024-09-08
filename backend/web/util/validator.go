package util

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Validator map[string]string

func NewValidator() Validator {
	return make(Validator)
}

func (v Validator) IsValid() bool {
	return len(v) == 0
}

func (v Validator) Add(field, message string) {
	v[field] = message
}

func (v Validator) Check(field, message string, ok bool) {
	if !ok {
		v[field] = message
	}
}

func (v Validator) CheckRequired(field, value string) {
	message := "This field is required"
	v.Check(field, message, strings.TrimSpace(value) != "")
}

func (v Validator) CheckMaxCharacters(field, value string, n int) {
	message := fmt.Sprintf("This field cannot me more than %d characters long", n)
	v.Check(field, message, utf8.RuneCountInString(value) <= n)
}
