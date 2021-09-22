package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type postStorage struct {
	conn *pgxpool.Pool
}

func NewPostStorage(conn *pgxpool.Pool) core.PostStorage {
	s := postStorage{
		conn: conn,
	}
	return &s
}

func (s *postStorage) Create(ctx context.Context, post *core.Post) error {
	stmt := `
		INSERT INTO post
			(url, title, author, body, updated, blog_id)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING post_id`
	row := s.conn.QueryRow(ctx, stmt,
		post.URL,
		post.Title,
		post.Author,
		post.Body,
		post.Updated,
		post.Blog.BlogID)

	err := row.Scan(&post.PostID)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		// https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return core.ErrExist
			}
		}
		return err
	}

	return nil
}

func (s *postStorage) Read(ctx context.Context, postID int) (core.Post, error) {
	stmt := `
		SELECT
			post.post_id,
			post.url,
			post.title,
			post.author,
			post.body,
			post.updated,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		WHERE post.post_id = $1`
	row := s.conn.QueryRow(ctx, stmt, postID)

	var post core.Post
	err := row.Scan(
		&post.PostID,
		&post.URL,
		&post.Title,
		&post.Author,
		&post.Body,
		&post.Updated,
		&post.Blog.BlogID,
		&post.Blog.FeedURL,
		&post.Blog.SiteURL,
		&post.Blog.Title,
	)
	if err != nil {
		return core.Post{}, err
	}

	return post, nil
}

func (s *postStorage) ReadAllByBlog(ctx context.Context, blogID int) ([]core.Post, error) {
	stmt := `
		SELECT
			post.post_id,
			post.url,
			post.title,
			post.author,
			post.body,
			post.updated,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		WHERE blog.blog_id = $1`
	rows, err := s.conn.Query(ctx, stmt, blogID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post
	for rows.Next() {
		var post core.Post
		err := rows.Scan(
			&post.PostID,
			&post.URL,
			&post.Title,
			&post.Author,
			&post.Body,
			&post.Updated,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postStorage) ReadRecent(ctx context.Context, n int) ([]core.Post, error) {
	stmt := `
		SELECT
			post.post_id,
			post.url,
			post.title,
			post.author,
			post.body,
			post.updated,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		ORDER BY post.updated DESC
		LIMIT $1`
	rows, err := s.conn.Query(ctx, stmt, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post
	for rows.Next() {
		var post core.Post
		err := rows.Scan(
			&post.PostID,
			&post.URL,
			&post.Title,
			&post.Author,
			&post.Body,
			&post.Updated,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postStorage) ReadRecentByAccount(ctx context.Context, accountID int, n int) ([]core.Post, error) {
	stmt := `
		SELECT
			post.post_id,
			post.url,
			post.title,
			post.body,
			post.updated,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN follow
			ON follow.blog_id = blog.blog_id
		WHERE follow.account_id = $1
		ORDER BY post.updated DESC
		LIMIT $2`
	rows, err := s.conn.Query(ctx, stmt, accountID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post
	for rows.Next() {
		var post core.Post
		err := rows.Scan(
			&post.PostID,
			&post.URL,
			&post.Title,
			&post.Body,
			&post.Updated,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postStorage) Delete(ctx context.Context, postID int) error {
	stmt := `
		DELETE
		FROM post
		WHERE post_id = $1`
	_, err := s.conn.Exec(ctx, stmt, postID)
	return err
}
