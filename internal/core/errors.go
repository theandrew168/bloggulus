package core

import (
	"errors"
)

var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("core: already exists")
	ErrNotExist = errors.New("core: does not exist")

	ErrRetry    = errors.New("core: retry storage operation")
	ErrConflict = errors.New("core: conflict in storage operation")
)
