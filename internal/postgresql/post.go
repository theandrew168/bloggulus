package postgresql

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/feed"
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
	// attempt to read post body, log and ignore any errors
	body, err := feed.ReadPostBody(*post)
	if err != nil {
		body = ""
		log.Println(err)
	}

	stmt := `
		INSERT INTO post
			(url, title, updated, body, blog_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING post_id`
	row := s.conn.QueryRow(ctx, stmt,
		post.URL,
		post.Title,
		post.Updated,
		body,
		post.Blog.BlogID)

	err = row.Scan(&post.PostID)
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
			post.updated,
			(array_agg(tag.name ORDER BY ts_rank_cd(post.content_index, to_tsquery(tag.name)) DESC)) as tags,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN tag
			ON to_tsquery(tag.name) @@ post.content_index
		WHERE post.post_id = $1
		GROUP BY 1,2,3,4,6,7,8,9`
	row := s.conn.QueryRow(ctx, stmt, postID)

	var post core.Post
	err := row.Scan(
		&post.PostID,
		&post.URL,
		&post.Title,
		&post.Updated,
		&post.Tags,
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
			post.updated,
			array_agg(tag.name ORDER BY ts_rank_cd(post.content_index, to_tsquery(tag.name)) DESC) as tags,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN tag
			ON to_tsquery(tag.name) @@ post.content_index
		WHERE blog.blog_id = $1
		GROUP BY 1,2,3,4,6,7,8,9
		ORDER BY post.updated DESC`
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
			&post.Updated,
			&post.Tags,
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

func (s *postStorage) ReadRecent(ctx context.Context, limit, offset int) ([]core.Post, error) {
	stmt := `
		SELECT
			post.post_id,
			post.url,
			post.title,
			post.updated,
			array_agg(tag.name ORDER BY ts_rank_cd(post.content_index, to_tsquery(tag.name)) DESC) as tags,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN tag
			ON to_tsquery(tag.name) @@ post.content_index
		GROUP BY 1,2,3,4,6,7,8,9
		ORDER BY post.updated DESC
		LIMIT $1
		OFFSET $2`
	rows, err := s.conn.Query(ctx, stmt, limit, offset)
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
			&post.Updated,
			&post.Tags,
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

func (s *postStorage) ReadSearch(ctx context.Context, query string, limit, offset int) ([]core.Post, error) {
	stmt := `
		SELECT
			post.post_id,
			post.url,
			post.title,
			post.updated,
			array_agg(tag.name ORDER BY ts_rank(post.content_index, to_tsquery(tag.name)) DESC) as tags,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		INNER JOIN tag
			ON to_tsquery(tag.name) @@ post.content_index
		WHERE post.content_index @@ websearch_to_tsquery('english',  $1)
		GROUP BY 1,2,3,4,6,7,8,9
		ORDER BY
			ts_rank(post.content_index, websearch_to_tsquery('english',  $1)) DESC
		LIMIT $2
		OFFSET $3`
	rows, err := s.conn.Query(ctx, stmt, query, limit, offset)
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
			&post.Updated,
			&post.Tags,
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

func (s *postStorage) CountRecent(ctx context.Context) (int, error) {
	stmt := `
		SELECT
			count(*)
		FROM post`
	row := s.conn.QueryRow(ctx, stmt)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *postStorage) CountSearch(ctx context.Context, query string) (int, error) {
	stmt := `
		SELECT
			count(*)
		FROM post
		WHERE content_index @@ websearch_to_tsquery('english',  $1)`
	row := s.conn.QueryRow(ctx, stmt, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
