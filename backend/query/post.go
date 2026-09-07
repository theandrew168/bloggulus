package query

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type PostDetails struct {
	ID          uuid.UUID `db:"id"`
	BlogID      uuid.UUID `db:"blog_id"`
	URL         string    `db:"url"`
	Title       string    `db:"title"`
	PublishedAt time.Time `db:"published_at"`
}

type PostQuery struct {
	conn postgres.Conn
}

func NewPost(conn postgres.Conn) *PostQuery {
	qry := PostQuery{
		conn: conn,
	}
	return &qry
}

// Powers the post details page (admin only).
func (qry *PostQuery) ReadDetailsByID(postID uuid.UUID) (PostDetails, error) {
	stmt := `
		SELECT
			post.id,
			post.blog_id,
			post.url,
			post.title,
			post.published_at
		FROM post
		WHERE post.id = $1;
	`

	rows, err := qry.conn.Query(context.Background(), stmt, postID)
	if err != nil {
		return PostDetails{}, err
	}

	details, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PostDetails])
	if err != nil {
		return PostDetails{}, postgres.CheckReadError(err)
	}

	return details, nil
}

// Powers the blog details page (admin only).
func (qry *PostQuery) ListDetailsByBlogID(blogID uuid.UUID) ([]PostDetails, error) {
	stmt := `
		SELECT
			post.id,
			post.blog_id,
			post.url,
			post.title,
			post.published_at
		FROM post
		WHERE post.blog_id = $1
		ORDER BY post.published_at DESC;
	`

	rows, err := qry.conn.Query(context.Background(), stmt, blogID)
	if err != nil {
		return nil, err
	}

	details, err := pgx.CollectRows(rows, pgx.RowToStructByName[PostDetails])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return details, nil
}
