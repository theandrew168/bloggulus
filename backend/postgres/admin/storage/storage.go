package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// ensure Storage interface is satisfied
var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	conn postgres.Conn

	blog *PostgresBlogStorage
	post *PostgresPostStorage
	tag  *PostgresTagStorage
}

func New(conn postgres.Conn) *Storage {
	s := Storage{
		conn: conn,

		blog: NewPostgresBlogStorage(conn),
		post: NewPostgresPostStorage(conn),
		tag:  NewPostgresTagStorage(conn),
	}
	return &s
}

func (s *Storage) Blog() storage.BlogStorage {
	return s.blog
}

func (s *Storage) Post() storage.PostStorage {
	return s.post
}

func (s *Storage) Tag() storage.TagStorage {
	return s.tag
}

// Based on:
// https://pkg.go.dev/github.com/jackc/pgx#hdr-Transactions
func (s *Storage) WithTransaction(operation func(store storage.Storage) error) error {
	// Calling the Begin() method on the connection creates a new pgx.Tx
	// object, which represents the in-progress database transaction.
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	// Defer a call to tx.Rollback() to ensure it is always called before the
	// function returns. If the transaction succeeds it will be already be
	// committed by the time tx.Rollback() is called, making tx.Rollback() a
	// no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
	// the changes before the function returns.
	defer tx.Rollback(context.Background())

	// Create a new Storage struct using the pgx.Tx as its Conn.
	store := New(tx)

	// Use the pgx.Tx-based Storage struct for this atomic operation.
	// If an error occurs within this operation, the transction will
	// be rolled back.
	err = operation(store)
	if err != nil {
		return err
	}

	// If there are no errors, the operation can be committed
	// to the database with the tx.Commit() method.
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
