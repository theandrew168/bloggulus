package query

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// List: recent, recent by account, search, search by account
// Count: all, all by account, search, search by account

type Article struct {
	Title       string    `db:"title" json:"title"`
	URL         string    `db:"url" json:"url"`
	BlogTitle   string    `db:"blog_title" json:"blogTitle"`
	BlogURL     string    `db:"blog_url" json:"blogURL"`
	PublishedAt time.Time `db:"published_at" json:"publishedAt"`
	Tags        []string  `db:"tags" json:"tags"`
}

func (qry *Query) ListRecentArticles(limit, offset int) ([]Article, error) {
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
			MAX(blog.title) as blog_title,
			MAX(blog.site_url) as blog_url,
			post.published_at,
			(array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL))[1:3] as tags
		FROM latest
		INNER JOIN post
			ON post.id = latest.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		GROUP BY post.id
		ORDER BY post.published_at DESC`

	rows, err := qry.conn.Query(context.Background(), stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (qry *Query) ListRecentArticlesByAccount(account *model.Account, limit, offset int) ([]Article, error) {
	stmt := `
		WITH latest AS (
			SELECT
				post.id
			FROM post
			INNER JOIN blog
				ON blog.id = post.blog_id
			INNER JOIN account_blog
				ON account_blog.blog_id = blog.id
				AND account_blog.account_id = $1
			ORDER BY post.published_at DESC
			LIMIT $2 OFFSET $3
		)
		SELECT
			post.title,
			post.url,
			MAX(blog.title) as blog_title,
			MAX(blog.site_url) as blog_url,
			post.published_at,
			(array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL))[1:3] as tags
		FROM latest
		INNER JOIN post
			ON post.id = latest.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		GROUP BY post.id
		ORDER BY post.published_at DESC`

	rows, err := qry.conn.Query(context.Background(), stmt, account.ID(), limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (qry *Query) ListRelevantArticles(search string, limit, offset int) ([]Article, error) {
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
			MAX(blog.title) as blog_title,
			MAX(blog.site_url) as blog_url,
			post.published_at,
			(array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL))[1:3] as tags
		FROM relevant
		INNER JOIN post
			ON post.id = relevant.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $1)
		GROUP BY post.id
		ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $1)) DESC`

	rows, err := qry.conn.Query(context.Background(), stmt, search, limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (qry *Query) ListRelevantArticlesByAccount(account *model.Account, search string, limit, offset int) ([]Article, error) {
	stmt := `
		WITH relevant AS (
			SELECT
				post.id
			FROM post
			INNER JOIN blog
				ON blog.id = post.blog_id
			INNER JOIN account_blog
				ON account_blog.blog_id = blog.id
				AND account_blog.account_id = $1
			ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $2)) DESC
			LIMIT $3 OFFSET $4
		)
		SELECT
			post.title,
			post.url,
			MAX(blog.title) as blog_title,
			MAX(blog.site_url) as blog_url,
			post.published_at,
			(array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL))[1:3] as tags
		FROM relevant
		INNER JOIN post
			ON post.id = relevant.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $2)
		GROUP BY post.id
		ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $2)) DESC`

	rows, err := qry.conn.Query(context.Background(), stmt, account.ID(), search, limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (qry *Query) CountRecentArticles() (int, error) {
	stmt := `
		SELECT count(*)
		FROM post`

	rows, err := qry.conn.Query(context.Background(), stmt)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (qry *Query) CountRecentArticlesByAccount(account *model.Account) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		INNER JOIN account_blog
			ON account_blog.blog_id = blog.id
			AND account_blog.account_id = $1`

	rows, err := qry.conn.Query(context.Background(), stmt, account.ID())
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (qry *Query) CountRelevantArticles(search string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $1)`

	rows, err := qry.conn.Query(context.Background(), stmt, search)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (qry *Query) CountRelevantArticlesByAccount(account *model.Account, search string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		INNER JOIN account_blog
			ON account_blog.blog_id = blog.id
			AND account_blog.account_id = $1
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $2)`

	rows, err := qry.conn.Query(context.Background(), stmt, account.ID(), search)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}
