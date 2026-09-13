package repository

import (
	"context"
	"strings"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/value"
)

type dbTag struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`

	MetaCreatedAt time.Time `db:"meta_created_at"`
	MetaUpdatedAt time.Time `db:"meta_updated_at"`
}

func marshalTag(tag *model.Tag) (dbTag, error) {
	t := dbTag{
		ID:            tag.ID(),
		Name:          tag.Name().Value(),
		MetaCreatedAt: tag.Meta().CreatedAt(),
		MetaUpdatedAt: tag.Meta().UpdatedAt(),
	}

	return t, nil
}

func (t dbTag) unmarshal() (*model.Tag, error) {
	name, err := value.NewName(t.Name)
	if err != nil {
		return nil, err
	}

	tag := model.LoadTag(model.LoadTagParams{
		ID:   t.ID,
		Name: name,
		Meta: model.LoadMeta(model.LoadMetaParams{
			CreatedAt: t.MetaCreatedAt,
			UpdatedAt: t.MetaUpdatedAt,
		}),
	})

	return tag, nil
}

type TagRepository struct {
	conn postgres.Conn
}

func NewTagRepository(conn postgres.Conn) *TagRepository {
	r := TagRepository{
		conn: conn,
	}
	return &r
}

func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	stmt := `
		-- name: Repository_Tag_Create
		INSERT INTO tag
			(id, name, meta_created_at, meta_updated_at)
		VALUES
			($1, $2, $3, $4);
	`

	row, err := marshalTag(tag)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.Name,
		row.MetaCreatedAt,
		row.MetaUpdatedAt,
	}

	_, err = r.conn.Exec(ctx, strings.TrimSpace(stmt), args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *TagRepository) Read(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	stmt := `
		-- name: Repository_Tag_Read
		SELECT
			tag.id,
			tag.name,
			tag.meta_created_at,
			tag.meta_updated_at
		FROM tag
		WHERE tag.id = $1;
	`

	rows, err := r.conn.Query(ctx, strings.TrimSpace(stmt), id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbTag])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *TagRepository) Delete(ctx context.Context, tag *model.Tag) error {
	stmt := `
		-- name: Repository_Tag_Delete
		DELETE FROM tag
		WHERE id = $1
		RETURNING id;
	`

	rows, err := r.conn.Query(ctx, strings.TrimSpace(stmt), tag.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
