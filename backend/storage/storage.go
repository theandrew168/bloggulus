package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/backend/postgres"
)

type Storage struct {
	conn postgres.Conn

	article *ArticleStorage
	blog    *BlogStorage
	post    *PostStorage
	tag     *TagStorage
	account *AccountStorage
	token   *TokenStorage
}

func New(conn postgres.Conn) *Storage {
	s := Storage{
		conn: conn,

		article: NewArticleStorage(conn),
		blog:    NewBlogStorage(conn),
		post:    NewPostStorage(conn),
		tag:     NewTagStorage(conn),
		account: NewAccountStoragee(conn),
		token:   NewTokenStorage(conn),
	}
	return &s
}

func (s *Storage) Article() *ArticleStorage {
	return s.article
}

func (s *Storage) Blog() *BlogStorage {
	return s.blog
}

func (s *Storage) Post() *PostStorage {
	return s.post
}

func (s *Storage) Tag() *TagStorage {
	return s.tag
}

func (s *Storage) Account() *AccountStorage {
	return s.account
}

func (s *Storage) Token() *TokenStorage {
	return s.token
}

func (s *Storage) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := s.conn.Exec(ctx, sql, args...)
	return err
}

// Based on:
// https://pkg.go.dev/github.com/jackc/pgx#hdr-Transactions
func (s *Storage) WithTransaction(operation func(store *Storage) error) error {
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
