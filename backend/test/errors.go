package test

import "errors"

var (
	// sentinel error used to rollback transactions
	ErrRollback = errors.New("test: rollback")
)
