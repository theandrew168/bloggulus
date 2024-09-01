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
	Title       string    `db:"title"`
	URL         string    `db:"url"`
	BlogTitle   string    `db:"blog_title"`
	BlogURL     string    `db:"blog_url"`
	PublishedAt time.Time `db:"published_at"`
	Tags        []string  `db:"tags"`
}

func (q *Query) ListArticles(limit, offset int) ([]Article, error) {
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
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL) as tags
		FROM latest
		INNER JOIN post
			ON post.id = latest.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		GROUP BY 1,2,3,4,5
		ORDER BY post.published_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (q *Query) ListArticlesByAccount(account *model.Account, limit, offset int) ([]Article, error) {
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
			blog.title as blog_title,
			blog.site_url as blog_url,
			post.published_at,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL) as tags
		FROM latest
		INNER JOIN post
			ON post.id = latest.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		GROUP BY 1,2,3,4,5
		ORDER BY post.published_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt, account.ID(), limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (q *Query) SearchArticles(search string, limit, offset int) ([]Article, error) {
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
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL) as tags
		FROM relevant
		INNER JOIN post
			ON post.id = relevant.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $1)
		GROUP BY 1,2,3,4,5,post.fts_data
		ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $1)) DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt, search, limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (q *Query) SearchArticlesByAccount(account *model.Account, search string, limit, offset int) ([]Article, error) {
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
			blog.title as blog_title,
			blog.site_url as blog_url,
			post.published_at,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.fts_data, plainto_tsquery('english', tag.name)) DESC), NULL) as tags
		FROM relevant
		INNER JOIN post
			ON post.id = relevant.id
		INNER JOIN blog
			ON blog.id = post.blog_id
		LEFT JOIN tag
			ON plainto_tsquery('english', tag.name) @@ post.fts_data
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $2)
		GROUP BY 1,2,3,4,5,post.fts_data
		ORDER BY ts_rank_cd(post.fts_data, websearch_to_tsquery('english',  $2)) DESC`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt, account.ID(), search, limit, offset)
	if err != nil {
		return nil, err
	}

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return articles, nil
}

func (q *Query) CountArticles() (int, error) {
	stmt := `
		SELECT count(*)
		FROM post`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (q *Query) CountArticlesByAccount(account *model.Account) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		INNER JOIN account_blog
			ON account_blog.blog_id = blog.id
			AND account_blog.account_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt, account.ID())
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (q *Query) CountSearchArticles(search string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $1)`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt, search)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (q *Query) CountSearchArticlesByAccount(account *model.Account, search string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		INNER JOIN blog
			ON blog.id = post.blog_id
		INNER JOIN account_blog
			ON account_blog.blog_id = blog.id
			AND account_blog.account_id = $1
		WHERE post.fts_data @@ websearch_to_tsquery('english',  $2)`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := q.conn.Query(ctx, stmt, account.ID(), search)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}
