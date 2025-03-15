package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

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

func marshalPost(post *model.Post) (dbPost, error) {
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

func (p dbPost) unmarshal() (*model.Post, error) {
	post := model.LoadPost(
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

type PostRepository struct {
	conn postgres.Conn
}

func NewPostRepository(conn postgres.Conn) *PostRepository {
	r := PostRepository{
		conn: conn,
	}
	return &r
}

func (r *PostRepository) Create(post *model.Post) error {
	stmt := `
		INSERT INTO post
			(id, blog_id, url, title, content, published_at, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)`

	row, err := marshalPost(post)
	if err != nil {
		return err
	}

	args := []any{
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

	_, err = r.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *PostRepository) Read(id uuid.UUID) (*model.Post, error) {
	stmt := `
		SELECT
			post.id,
			post.blog_id,
			post.url,
			post.title,
			post.content,
			post.published_at,
			post.created_at,
			post.updated_at
		FROM post
		WHERE post.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *PostRepository) ReadByURL(url string) (*model.Post, error) {
	stmt := `
		SELECT
			post.id,
			post.blog_id,
			post.url,
			post.title,
			post.content,
			post.published_at,
			post.created_at,
			post.updated_at
		FROM post
		WHERE post.url = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, url)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *PostRepository) ListByBlog(blog *model.Blog) ([]*model.Post, error) {
	stmt := `
		SELECT
			post.id,
			post.blog_id,
			post.url,
			post.title,
			post.content,
			post.published_at,
			post.created_at,
			post.updated_at
		FROM post
		WHERE post.blog_id = $1
		ORDER BY post.published_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, blog.ID())
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var posts []*model.Post
	for _, row := range postRows {
		post, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) CountByBlog(blog *model.Blog) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		WHERE post.blog_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, blog.ID())
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (r *PostRepository) Update(post *model.Post) error {
	now := timeutil.Now()
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

	rows, err := r.conn.Query(ctx, stmt, args...)
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

func (r *PostRepository) Delete(post *model.Post) error {
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

	rows, err := r.conn.Query(ctx, stmt, post.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
