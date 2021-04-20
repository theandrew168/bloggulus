package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"
)

type postStorage struct {
	db *pgxpool.Pool
}

func NewPostStorage(db *pgxpool.Pool) storage.Post {
	return &postStorage{
		db: db,
	}
}

func (s *postStorage) Create(ctx context.Context, post *models.Post) (*models.Post, error) {
	command := `
		INSERT INTO post
			(blog_id, url, title, updated)
		VALUES
			($1, $2, $3, $4)
		RETURNING post_id`
	row := s.db.QueryRow(ctx, command, post.BlogID, post.URL, post.Title, post.Updated)

	err := row.Scan(&post.PostID)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		// https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, storage.ErrDuplicateModel
			}
		}
		return nil, err
	}

	return post, nil
}

func (s *postStorage) Read(ctx context.Context, postID int) (*models.Post, error) {
	query := `
		SELECT
			post.*,
			blog.*
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		WHERE post.post_id = $1`
	row := s.db.QueryRow(ctx, query, postID)

	var post models.Post
	err := row.Scan(
		&post.PostID,
		&post.BlogID,
		&post.URL,
		&post.Title,
		&post.Updated,
		&post.Blog.BlogID,
		&post.Blog.FeedURL,
		&post.Blog.SiteURL,
		&post.Blog.Title,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *postStorage) ReadRecent(ctx context.Context, n int) ([]*models.Post, error) {
	query := `
		SELECT
			post.*,
			blog.*
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		ORDER BY post.updated DESC
		LIMIT $1`
	rows, err := s.db.Query(ctx, query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.PostID,
			&post.BlogID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *postStorage) ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*models.Post, error) {
	query := `
		SELECT
			post.*,
			blog.*
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN account_blog
			ON account_blog.blog_id = blog.blog_id
		WHERE account_blog.account_id = $1
		ORDER BY post.updated DESC
		LIMIT $2`
	rows, err := s.db.Query(ctx, query, accountID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.PostID,
			&post.BlogID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *postStorage) Delete(ctx context.Context, postID int) error {
	command := `
		DELETE
		FROM post
		WHERE post_id = $1`
	_, err := s.db.Exec(ctx, command, postID)
	return err
}
