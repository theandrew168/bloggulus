package test

import "errors"

// TODO: move to domain/storage?

var (
	// sentinel error used to rollback transactions
	ErrRollback = errors.New("test: rollback")
)
