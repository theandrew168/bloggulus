package core

import (
	"errors"
)

var (
	ErrExist    = errors.New("core: already exists")
	ErrNotExist = errors.New("core: does not exist")
)
