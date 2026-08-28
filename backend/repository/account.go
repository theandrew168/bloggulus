package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"slices"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/value"
)

type dbAccount struct {
	ID              uuid.UUID   `db:"id"`
	Username        string      `db:"username"`
	IsAdmin         bool        `db:"is_admin"`
	FollowedBlogIDs []uuid.UUID `db:"followed_blog_ids"`
	MetaCreatedAt   time.Time   `db:"meta_created_at"`
	MetaUpdatedAt   time.Time   `db:"meta_updated_at"`
	MetaVersion     int         `db:"meta_version"`
}

func marshalAccount(account *model.Account) (dbAccount, error) {
	a := dbAccount{
		ID:              account.ID(),
		Username:        account.Username().Value(),
		IsAdmin:         account.IsAdmin(),
		FollowedBlogIDs: account.FollowedBlogIDs(),
		MetaCreatedAt:   account.Meta().CreatedAt(),
		MetaUpdatedAt:   account.Meta().UpdatedAt(),
		MetaVersion:     account.Meta().Version().Value(),
	}
	return a, nil
}

func (a dbAccount) unmarshal() (*model.Account, error) {
	username, err := value.NewName(a.Username)
	if err != nil {
		return nil, err
	}

	version, err := value.NewCount(a.MetaVersion)
	if err != nil {
		return nil, err
	}
	meta := model.LoadMeta(a.MetaCreatedAt, a.MetaUpdatedAt, version)

	account := model.LoadAccount(
		a.ID,
		username,
		a.IsAdmin,
		a.FollowedBlogIDs,
		meta,
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
			(id, username, meta_created_at, meta_updated_at, meta_version)
		VALUES
			($1, $2, $3, $4, $5)`

	row, err := marshalAccount(account)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.Username,
		row.MetaCreatedAt,
		row.MetaUpdatedAt,
		row.MetaVersion,
	}

	_, err = r.conn.Exec(context.Background(), stmt, args...)
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
			account.meta_created_at,
			account.meta_updated_at,
			account.meta_version
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		WHERE account.id = $1
		GROUP BY account.id`

	rows, err := r.conn.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbAccount])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *AccountRepository) ReadByUsername(username value.Name) (*model.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.is_admin,
			ARRAY_AGG(account_blog.blog_id) AS followed_blog_ids,
			account.meta_created_at,
			account.meta_updated_at,
			account.meta_version
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		WHERE account.username = $1
		GROUP BY account.id`

	rows, err := r.conn.Query(context.Background(), stmt, username.Value())
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
			account.meta_created_at,
			account.meta_updated_at,
			account.meta_version
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		INNER JOIN session
			ON session.account_id = account.id
		WHERE session.hash = $1
		GROUP BY account.id`

	hashBytes := sha256.Sum256([]byte(sessionID))
	hash := hex.EncodeToString(hashBytes[:])

	rows, err := r.conn.Query(context.Background(), stmt, hash)
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
			account.meta_created_at,
			account.meta_updated_at,
			account.meta_version
		FROM account
		LEFT JOIN account_blog
			ON account_blog.account_id = account.id
		GROUP BY account.id
		ORDER BY account.meta_created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.conn.Query(context.Background(), stmt, limit, offset)
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

	rows, err := r.conn.Query(context.Background(), stmt, account.ID())
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

	// TODO: Optim Oppty: Batch these additions and removals.

	// Add and remove blogs as necessary.
	stmtFollow := `
		INSERT INTO account_blog
			(account_id, blog_id)
		VALUES ($1, $2)`
	for _, blogID := range blogsToFollow {
		_, err = r.conn.Exec(context.Background(), stmtFollow, account.ID(), blogID)
		if err != nil {
			return postgres.CheckCreateError(err)
		}
	}

	stmtUnfollow := `
		DELETE FROM account_blog
		WHERE account_id = $1 AND blog_id = $2`
	for _, blogID := range blogsToUnfollow {
		_, err = r.conn.Exec(context.Background(), stmtUnfollow, account.ID(), blogID)
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

	rows, err := r.conn.Query(context.Background(), stmt, account.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
