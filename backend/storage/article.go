package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type dbArticle struct {
	Title       string    `db:"title"`
	URL         string    `db:"url"`
	BlogTitle   string    `db:"blog_title"`
	BlogURL     string    `db:"blog_url"`
	PublishedAt time.Time `db:"published_at"`
	Tags        []string  `db:"tags"`
}

func (a dbArticle) unmarshal() (*model.Article, error) {
	article := model.LoadArticle(
		a.Title,
		a.URL,
		a.BlogTitle,
		a.BlogURL,
		a.PublishedAt,
		a.Tags,
	)
	return article, nil
}

type ArticleStorage struct {
	conn postgres.Conn
}

func NewArticleStorage(conn postgres.Conn) *ArticleStorage {
	s := ArticleStorage{
		conn: conn,
	}
	return &s
}

func (s *ArticleStorage) List(limit, offset int) ([]*model.Article, error) {
	stmt := `
		WITH latest AS (
			SELECT
				post.id
			FROM post
			ORDER BY post.published_at DESC
			LIMIT $1 OFFSET $2
		)
		SELECT
			post.title,
			post.url,
			blog.title as blog_title,
			blog.site_url as blog_url,
			post.published_at,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, to_tsquery(tag.name)) DESC), NULL) as tags
		FROM latest
		INNER JOIN post
			ON post.id = latest.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ post.fts_data
		GROUP BY 1,2,3,4,5
		ORDER BY post.published_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	articleRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbArticle])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var articles []*model.Article
	for _, row := range articleRows {
		article, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (s *ArticleStorage) ListSearch(search string, limit, offset int) ([]*model.Article, error) {
	stmt := `
		WITH relevant AS (
			SELECT
				post.id
			FROM post
			ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $1)) DESC
			LIMIT $2 OFFSET $3
		)
		SELECT
			post.title,
			post.url,
			blog.title as blog_title,
			blog.site_url as blog_url,
			post.published_at,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, to_tsquery(tag.name)) DESC), NULL) as tags
		FROM relevant
		INNER JOIN post
			ON post.id = relevant.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ post.fts_data
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $1)
		GROUP BY 1,2,3,4,5,post.fts_data
		ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $1)) DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, search, limit, offset)
	if err != nil {
		return nil, err
	}

	articleRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbArticle])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var articles []*model.Article
	for _, row := range articleRows {
		article, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (s *ArticleStorage) Count() (int, error) {
	stmt := `
		SELECT count(*)
		FROM post`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (s *ArticleStorage) CountSearch(search string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		WHERE fts_data @@ websearch_to_tsquery('english',  $1)`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, search)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}
