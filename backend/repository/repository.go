package repository

import (
	"context"

	"github.com/theandrew168/bloggulus/backend/postgres"
)

type Repository struct {
	conn postgres.Conn

	blog        *BlogRepository
	post        *PostRepository
	tag         *TagRepository
	account     *AccountRepository
	session     *SessionRepository
	accountBlog *AccountBlogRepository
	page        *PageRepository
	accountPage *AccountPageRepository
}

func New(conn postgres.Conn) *Repository {
	r := Repository{
		conn: conn,

		blog:        NewBlogRepository(conn),
		post:        NewPostRepository(conn),
		tag:         NewTagRepository(conn),
		account:     NewAccountRepository(conn),
		session:     NewSessionRepository(conn),
		accountBlog: NewAccountBlogRepository(conn),
		page:        NewPageRepository(conn),
		accountPage: NewAccountPageRepository(conn),
	}
	return &r
}

func (r *Repository) Blog() *BlogRepository {
	return r.blog
}

func (r *Repository) Post() *PostRepository {
	return r.post
}

func (r *Repository) Tag() *TagRepository {
	return r.tag
}

func (r *Repository) Account() *AccountRepository {
	return r.account
}

func (r *Repository) Session() *SessionRepository {
	return r.session
}

func (r *Repository) AccountBlog() *AccountBlogRepository {
	return r.accountBlog
}

func (r *Repository) Page() *PageRepository {
	return r.page
}

func (r *Repository) AccountPage() *AccountPageRepository {
	return r.accountPage
}

func (r *Repository) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := r.conn.Exec(ctx, sql, args...)
	return err
}

// Based on:
// https://pkg.go.dev/github.com/jackc/pgx#hdr-Transactions
func (r *Repository) WithTransaction(operation func(repo *Repository) error) error {
	// Calling the Begin() method on the connection creates a new pgx.Tx
	// object, which represents the in-progress database transaction.
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	// Defer a call to tx.Rollback() to ensure it is always called before the
	// function returns. If the transaction succeeds it will be already be
	// committed by the time tx.Rollback() is called, making tx.Rollback() a
	// no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
	// the changes before the function returns.
	defer tx.Rollback(context.Background())

	// Create a new Repository struct using the pgx.Tx as its Conn. Note
	// that this new repo will be backed by single connection and not
	// a pool (therefore only one query can be executed at a time). When
	// inside of a transaction, the connection is NOT concurrency safe.
	repo := New(tx)

	// Use the pgx.Tx-based Repository struct for this atomic operation.
	// If an error occurs within this operation, the transction will
	// be rolled back.
	err = operation(repo)
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
