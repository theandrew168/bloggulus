package test

import "errors"

var (
	// sentinal error used to skip committing transactions
	ErrSkipCommit = errors.New("test: skip commit")
)
