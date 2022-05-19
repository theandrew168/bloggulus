package database

import (
	"errors"
)

var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("database: already exists")
	ErrNotExist = errors.New("database: does not exist")

	// storage errors
	ErrRetry    = errors.New("database: retry storage operation")
	ErrConflict = errors.New("database: conflict in storage operation")
)
