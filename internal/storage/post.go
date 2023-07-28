package storage

import (
	"context"

	"github.com/theandrew168/bloggulus"
	"github.com/theandrew168/bloggulus/internal/database"
)

type Post struct {
	db database.Conn
}

func NewPost(db database.Conn) *Post {
	s := Post{
		db: db,
	}
	return &s
}

func (s *Post) Create(post *bloggulus.Post) error {
	stmt := `
		INSERT INTO post
			(url, title, updated, body, blog_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`

	args := []interface{}{
		post.URL,
		post.Title,
		post.Updated,
		post.Body,
		post.Blog.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &post.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Post) Read(id int) (bloggulus.Post, error) {
	stmt := `
		SELECT
			post.id,
			post.url,
			post.title,
			post.updated,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.content_index, to_tsquery(tag.name)) DESC), NULL) as tags,
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ post.content_index
		WHERE post.id = $1
		GROUP BY 1,2,3,4,6,7,8,9`

	var post bloggulus.Post
	dest := []interface{}{
		&post.ID,
		&post.URL,
		&post.Title,
		&post.Updated,
		&post.Tags,
		&post.Blog.ID,
		&post.Blog.FeedURL,
		&post.Blog.SiteURL,
		&post.Blog.Title,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return bloggulus.Post{}, err
	}

	return post, nil
}

func (s *Post) ReadAll(limit, offset int) ([]bloggulus.Post, error) {
	stmt := `
		WITH posts AS (
			SELECT
				post.id,
				post.url,
				post.title,
				post.updated,
				post.content_index,
				blog.id AS blog_id,
				blog.feed_url AS blog_feed_url,
				blog.site_url AS blog_site_url,
				blog.title AS blog_title
			FROM post
			INNER JOIN blog
				ON blog.id = post.blog_id
			ORDER BY post.updated DESC
			LIMIT $1 OFFSET $2
		)
		SELECT
			posts.id,
			posts.url,
			posts.title,
			posts.updated,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(posts.content_index, to_tsquery(tag.name)) DESC), NULL) as tags,
			posts.blog_id,
			posts.blog_feed_url,
			posts.blog_site_url,
			posts.blog_title
		FROM posts
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ posts.content_index
		GROUP BY 1,2,3,4,6,7,8,9
		ORDER BY posts.updated DESC`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// use make here to encode JSON as an empty array instead of null
	posts := make([]bloggulus.Post, 0)
	for rows.Next() {
		var post bloggulus.Post
		dest := []interface{}{
			&post.ID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Tags,
			&post.Blog.ID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return posts, nil
}

func (s *Post) ReadAllByBlog(blog bloggulus.Blog, limit, offset int) ([]bloggulus.Post, error) {
	stmt := `
		SELECT
			post.id,
			post.url,
			post.title,
			post.updated,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.content_index, to_tsquery(tag.name)) DESC), NULL) as tags,
			blog.id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ post.content_index
		WHERE blog.id = $1
		GROUP BY 1,2,3,4,6,7,8,9
		ORDER BY post.updated DESC
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, blog.ID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]bloggulus.Post, 0)
	for rows.Next() {
		var post bloggulus.Post
		dest := []interface{}{
			&post.ID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Tags,
			&post.Blog.ID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return posts, nil
}

func (s *Post) Search(query string, limit, offset int) ([]bloggulus.Post, error) {
	stmt := `
		WITH posts AS (
			SELECT
				post.id,
				post.url,
				post.title,
				post.updated,
				post.content_index,
				blog.id AS blog_id,
				blog.feed_url AS blog_feed_url,
				blog.site_url AS blog_site_url,
				blog.title AS blog_title
			FROM post
			INNER JOIN blog
				ON blog.id = post.blog_id
			WHERE post.content_index @@ websearch_to_tsquery('english',  $1)
			ORDER BY ts_rank_cd(post.content_index, websearch_to_tsquery('english',  $1)) DESC
			LIMIT $2 OFFSET $3
		)
		SELECT
			posts.id,
			posts.url,
			posts.title,
			posts.updated,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(posts.content_index, to_tsquery(tag.name)) DESC), NULL) as tags,
			posts.blog_id,
			posts.blog_feed_url,
			posts.blog_site_url,
			posts.blog_title
		FROM posts
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ content_index
		GROUP BY 1,2,3,4,6,7,8,9,posts.content_index
		ORDER BY ts_rank_cd(posts.content_index, websearch_to_tsquery('english',  $1)) DESC`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]bloggulus.Post, 0)
	for rows.Next() {
		var post bloggulus.Post
		dest := []interface{}{
			&post.ID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Tags,
			&post.Blog.ID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Post) Count() (int, error) {
	stmt := `
		SELECT count(*)
		FROM post`

	var count int
	dest := []interface{}{
		&count,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt)
	err := database.Scan(row, dest...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Post) CountSearch(query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		WHERE content_index @@ websearch_to_tsquery('english',  $1)`

	var count int
	dest := []interface{}{
		&count,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, query)
	err := database.Scan(row, dest...)
	if err != nil {
		return 0, err
	}

	return count, nil
}
