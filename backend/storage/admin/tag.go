package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type dbTag struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func marshalTag(tag *admin.Tag) (dbTag, error) {
	t := dbTag{
		ID:        tag.ID(),
		Name:      tag.Name(),
		CreatedAt: tag.CreatedAt(),
		UpdatedAt: tag.UpdatedAt(),
	}
	return t, nil
}

func (t dbTag) unmarshal() (*admin.Tag, error) {
	tag := admin.LoadTag(
		t.ID,
		t.Name,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return tag, nil
}

type TagStorage struct {
	conn postgres.Conn
}

func NewTagStorage(conn postgres.Conn) *TagStorage {
	s := TagStorage{
		conn: conn,
	}
	return &s
}

func (s *TagStorage) Create(tag *admin.Tag) error {
	stmt := `
		INSERT INTO tag
			(id, name, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)`

	row, err := marshalTag(tag)
	if err != nil {
		return err
	}

	args := []interface{}{
		row.ID,
		row.Name,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	_, err = s.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (s *TagStorage) Read(id uuid.UUID) (*admin.Tag, error) {
	stmt := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM tag
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbTag])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (s *TagStorage) List(limit, offset int) ([]*admin.Tag, error) {
	stmt := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM tag
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	tagRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbTag])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var tags []*admin.Tag
	for _, row := range tagRows {
		tag, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (repo *TagStorage) Delete(tag *admin.Tag) error {
	stmt := `
		DELETE FROM tag
		WHERE id = $1
		RETURNING id`

	err := tag.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), postgres.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, tag.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
