package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// ensure PostStorage interface is satisfied
var _ storage.PostStorage = (*PostgresPostStorage)(nil)

type dbPost struct {
	ID          uuid.UUID `db:"id"`
	BlogID      uuid.UUID `db:"blog_id"`
	URL         string    `db:"url"`
	Title       string    `db:"title"`
	Content     string    `db:"content"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func marshalPost(post *admin.Post) (dbPost, error) {
	p := dbPost{
		ID:          post.ID(),
		BlogID:      post.BlogID(),
		URL:         post.URL(),
		Title:       post.Title(),
		Content:     post.Content(),
		PublishedAt: post.PublishedAt(),
		CreatedAt:   post.CreatedAt(),
		UpdatedAt:   post.UpdatedAt(),
	}
	return p, nil
}

func (p dbPost) unmarshal() (*admin.Post, error) {
	post := admin.LoadPost(
		p.ID,
		p.BlogID,
		p.URL,
		p.Title,
		p.Content,
		p.PublishedAt,
		p.CreatedAt,
		p.UpdatedAt,
	)
	return post, nil
}

type PostgresPostStorage struct {
	conn postgres.Conn
}

func NewPostgresPostStorage(conn postgres.Conn) *PostgresPostStorage {
	s := PostgresPostStorage{
		conn: conn,
	}
	return &s
}

func (s *PostgresPostStorage) Create(post *admin.Post) error {
	stmt := `
		INSERT INTO post
			(id, blog_id, url, title, content, published_at, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)`

	row, err := marshalPost(post)
	if err != nil {
		return err
	}

	args := []interface{}{
		row.ID,
		row.BlogID,
		row.URL,
		row.Title,
		row.Content,
		row.PublishedAt,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err = s.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (s *PostgresPostStorage) Read(id uuid.UUID) (*admin.Post, error) {
	stmt := `
		SELECT
			id,
			blog_id,
			url,
			title,
			content,
			published_at,
			created_at,
			updated_at
		FROM post
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (s *PostgresPostStorage) ReadByURL(url string) (*admin.Post, error) {
	stmt := `
		SELECT
			id,
			blog_id,
			url,
			title,
			content,
			published_at,
			created_at,
			updated_at
		FROM post
		WHERE url = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, url)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (s *PostgresPostStorage) List(limit, offset int) ([]*admin.Post, error) {
	stmt := `
		SELECT
			id,
			blog_id,
			url,
			title,
			content,
			published_at,
			created_at,
			updated_at
		FROM post
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var posts []*admin.Post
	for _, row := range postRows {
		post, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresPostStorage) ListByBlog(blog *admin.Blog, limit, offset int) ([]*admin.Post, error) {
	stmt := `
		SELECT
			id,
			blog_id,
			url,
			title,
			content,
			published_at,
			created_at,
			updated_at
		FROM post
		WHERE post.blog_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, blog.ID(), limit, offset)
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var posts []*admin.Post
	for _, row := range postRows {
		post, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresPostStorage) Update(post *admin.Post) error {
	now := time.Now().UTC()
	stmt := `
		UPDATE post
		SET
			url = $1,
			title = $2,
			content = $3,
			published_at = $4,
			updated_at = $5
		WHERE id = $6
		  AND updated_at = $7
		RETURNING updated_at`

	row, err := marshalPost(post)
	if err != nil {
		return err
	}

	args := []any{
		row.URL,
		row.Title,
		row.Content,
		row.PublishedAt,
		now,
		row.ID,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[time.Time])
	if err != nil {
		return postgres.CheckUpdateError(err)
	}

	post.SetUpdatedAt(now)
	return nil
}

func (repo *PostgresPostStorage) Delete(post *admin.Post) error {
	stmt := `
		DELETE FROM post
		WHERE id = $1
		RETURNING id`

	err := post.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, post.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
