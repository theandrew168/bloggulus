package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type dbPage struct {
	ID        uuid.UUID `db:"id"`
	URL       string    `db:"url"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func marshalPage(page *model.Page) (dbPage, error) {
	p := dbPage{
		ID:        page.ID(),
		URL:       page.URL(),
		Title:     page.Title(),
		Content:   page.Content(),
		CreatedAt: page.CreatedAt(),
		UpdatedAt: page.UpdatedAt(),
	}
	return p, nil
}

func (p dbPage) unmarshal() (*model.Page, error) {
	page := model.LoadPage(
		p.ID,
		p.URL,
		p.Title,
		p.Content,
		p.CreatedAt,
		p.UpdatedAt,
	)
	return page, nil
}

type PageRepository struct {
	conn postgres.Conn
}

func NewPageRepository(conn postgres.Conn) *PageRepository {
	r := PageRepository{
		conn: conn,
	}
	return &r
}

func (r *PageRepository) Create(page *model.Page) error {
	stmt := `
		INSERT INTO page
			(id, url, title, content, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)`

	row, err := marshalPage(page)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.URL,
		row.Title,
		row.Content,
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

func (r *PageRepository) Read(id uuid.UUID) (*model.Page, error) {
	stmt := `
		SELECT
			page.id,
			page.url,
			page.title,
			page.content,
			page.created_at,
			page.updated_at
		FROM page
		WHERE page.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPage])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *PageRepository) ReadByURL(url string) (*model.Page, error) {
	stmt := `
		SELECT
			page.id,
			page.url,
			page.title,
			page.content,
			page.created_at,
			page.updated_at
		FROM page
		WHERE page.url = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, url)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbPage])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *PageRepository) ListByAccount(account *model.Account, limit, offset int) ([]*model.Page, error) {
	stmt := `
		SELECT
			page.id,
			page.url,
			page.title,
			page.content,
			page.created_at,
			page.updated_at
		FROM page
		INNER JOIN account_page
			ON account_page.page_id = page.id
		WHERE account_page.account_id = $1
		ORDER BY page.created_at DESC
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, account.ID(), limit, offset)
	if err != nil {
		return nil, err
	}

	pageRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbPage])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var pages []*model.Page
	for _, row := range pageRows {
		page, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func (r *PageRepository) CountByAccount(account *model.Account) (int, error) {
	stmt := `
		SELECT count(*)
		FROM page
		INNER JOIN account_page
			ON account_page.page_id = page.id
		WHERE account_page.account_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, account.ID())
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (r *PageRepository) Delete(page *model.Page) error {
	stmt := `
		DELETE FROM page
		WHERE id = $1
		RETURNING id`

	err := page.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt, page.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
