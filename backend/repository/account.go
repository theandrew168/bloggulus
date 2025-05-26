package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type dbAccount struct {
	ID              uuid.UUID   `db:"id"`
	Username        string      `db:"username"`
	IsAdmin         bool        `db:"is_admin"`
	FollowedBlogIDs []uuid.UUID `db:"followed_blog_ids"`
	CreatedAt       time.Time   `db:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at"`
}

func marshalAccount(account *model.Account) (dbAccount, error) {
	a := dbAccount{
		ID:              account.ID(),
		Username:        account.Username(),
		IsAdmin:         account.IsAdmin(),
		FollowedBlogIDs: account.FollowedBlogIDs(),
		CreatedAt:       account.CreatedAt(),
		UpdatedAt:       account.UpdatedAt(),
	}
	return a, nil
}

func (a dbAccount) unmarshal() (*model.Account, error) {
	account := model.LoadAccount(
		a.ID,
		a.Username,
		a.IsAdmin,
		a.FollowedBlogIDs,
		a.CreatedAt,
		a.UpdatedAt,
	)
	return account, nil
}

type AccountRepository struct {
	conn postgres.Conn
}

func NewAccountRepository(conn postgres.Conn) *AccountRepository {
	r := AccountRepository{
		conn: conn,
	}
	return &r
}

func (r *AccountRepository) Create(account *model.Account) error {
	stmt := `
		INSERT INTO account
			(id, username, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)`

	row, err := marshalAccount(account)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.Username,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err = r.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *AccountRepository) Read(id uuid.UUID) (*model.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.is_admin,
			ARRAY_AGG(account_blog.blog_id) AS followed_blog_ids,
			account.created_at,
			account.updated_at
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		WHERE account.id = $1
		GROUP BY account.id`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *AccountRepository) ReadByUsername(username string) (*model.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.is_admin,
			ARRAY_AGG(account_blog.blog_id) AS followed_blog_ids,
			account.created_at,
			account.updated_at
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		WHERE account.username = $1
		GROUP BY account.id`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, username)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *AccountRepository) ReadBySessionID(sessionID string) (*model.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.is_admin,
			ARRAY_AGG(account_blog.blog_id) AS followed_blog_ids,
			account.created_at,
			account.updated_at
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		INNER JOIN session
			ON session.account_id = account.id
		WHERE session.hash = $1
		GROUP BY account.id`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	hashBytes := sha256.Sum256([]byte(sessionID))
	hash := hex.EncodeToString(hashBytes[:])

	rows, err := r.conn.Query(ctx, stmt, hash)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *AccountRepository) List(limit, offset int) ([]*model.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.is_admin,
			ARRAY_AGG(account_blog.blog_id) AS followed_blog_ids,
			account.created_at,
			account.updated_at
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		GROUP BY account.id
		ORDER BY account.created_at DESC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	accountRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var accounts []*model.Account
	for _, row := range accountRows {
		account, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *AccountRepository) Update(account *model.Account) error {
	// List blogs currently being followed in the database.
	stmt := `
		SELECT
			account_blog.blog_id
		FROM account_blog
		WHERE account_blog.account_id = $1`

	rows, err := QueryWithTimeout(r.conn, stmt, account.ID())
	if err != nil {
		return err
	}

	followedBlogIDs, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckListError(err)
	}

	// Set diff to find which blogs to add or remove.
	var blogsToFollow []uuid.UUID
	for _, blogID := range account.FollowedBlogIDs() {
		if !slices.Contains(followedBlogIDs, blogID) {
			blogsToFollow = append(blogsToFollow, blogID)
		}
	}

	var blogsToUnfollow []uuid.UUID
	for _, blogID := range followedBlogIDs {
		if !slices.Contains(account.FollowedBlogIDs(), blogID) {
			blogsToUnfollow = append(blogsToUnfollow, blogID)
		}
	}

	// Add and remove blogs as necessary.
	stmtFollow := `
		INSERT INTO account_blog
			(account_id, blog_id)
		VALUES ($1, $2)`
	for _, blogID := range blogsToFollow {
		err = ExecWithTimeout(r.conn, stmtFollow, account.ID(), blogID)
		if err != nil {
			return postgres.CheckCreateError(err)
		}
	}

	stmtUnfollow := `
		DELETE FROM account_blog
		WHERE account_id = $1 AND blog_id = $2`
	for _, blogID := range blogsToUnfollow {
		err = ExecWithTimeout(r.conn, stmtUnfollow, account.ID(), blogID)
		if err != nil {
			return postgres.CheckDeleteError(err)
		}
	}

	return nil
}

func (r *AccountRepository) Delete(account *model.Account) error {
	stmt := `
		DELETE FROM account
		WHERE id = $1
		RETURNING id`

	err := account.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, account.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
