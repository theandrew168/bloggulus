package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/domain/reader"
	"github.com/theandrew168/bloggulus/backend/domain/reader/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// ensure PostStorage interface is satisfied
var _ storage.PostStorage = (*PostgresPostStorage)(nil)

type dbPost struct {
	Title       string    `db:"title"`
	URL         string    `db:"url"`
	BlogTitle   string    `db:"blog_title"`
	BlogURL     string    `db:"blog_url"`
	PublishedAt time.Time `db:"published_at"`
	Tags        []string  `db:"tags"`
}

func (p dbPost) unmarshal() (*reader.Post, error) {
	post := reader.LoadPost(
		p.Title,
		p.URL,
		p.BlogTitle,
		p.BlogURL,
		p.PublishedAt,
		p.Tags,
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

// TODO: use templates to condense these into a single query?

func (s *PostgresPostStorage) List(limit, offset int) ([]*reader.Post, error) {
	stmt := `
		SELECT
			post.title,
			post.url,
			blog.title as blog_title,
			blog.site_url as blog_url,
			post.published_at,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, to_tsquery(tag.name)) DESC), NULL) as tags
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ post.fts_data
		GROUP BY 1,2,3,4,5
		ORDER BY post.published_at DESC
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

	var posts []*reader.Post
	for _, row := range postRows {
		post, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresPostStorage) Search(query string, limit, offset int) ([]*reader.Post, error) {
	stmt := `
		SELECT
			post.title,
			post.url,
			blog.title as blog_title,
			blog.site_url as blog_url,
			post.published_at,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, to_tsquery(tag.name)) DESC), NULL) as tags
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ post.fts_data
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $1)
		GROUP BY 1,2,3,4,5,post.fts_data
		ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $1)) DESC
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, query, limit, offset)
	if err != nil {
		return nil, err
	}

	postRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPost])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var posts []*reader.Post
	for _, row := range postRows {
		post, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
