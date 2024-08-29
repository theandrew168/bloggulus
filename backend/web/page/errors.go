package page

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Errors map[string]string

func NewErrors() Errors {
	return make(Errors)
}

func (e Errors) OK() bool {
	return len(e) == 0
}

func (e Errors) Add(field, message string) {
	e[field] = message
}

func (e Errors) Check(field, message string, ok bool) {
	if !ok {
		e[field] = message
	}
}

func (e Errors) CheckRequired(field, value string) {
	message := "This field is required"
	e.Check(field, message, strings.TrimSpace(value) != "")
}

func (e Errors) CheckMaxCharacters(field, value string, n int) {
	message := fmt.Sprintf("This field cannot me more than %d characters long", n)
	e.Check(field, message, utf8.RuneCountInString(value) <= n)
}
