package repository

import (
	"context"
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

	metaParams := model.LoadMetaParams{
		CreatedAt: t.MetaCreatedAt,
		UpdatedAt: t.MetaUpdatedAt,
	}

	tag := model.LoadTag(
		t.ID,
		name,
		model.LoadMeta(metaParams),
	)
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

func (r *TagRepository) Create(tag *model.Tag) error {
	stmt := `
		INSERT INTO tag
			(id, name, meta_created_at, meta_updated_at)
		VALUES
			($1, $2, $3, $4)`

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

	_, err = r.conn.Exec(context.Background(), stmt, args...)
	if err != nil {
		return postgres.CheckCreateError(err)
	}

	return nil
}

func (r *TagRepository) Read(id uuid.UUID) (*model.Tag, error) {
	stmt := `
		SELECT
			tag.id,
			tag.name,
			tag.meta_created_at,
			tag.meta_updated_at
		FROM tag
		WHERE tag.id = $1`

	rows, err := r.conn.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbTag])
	if err != nil {
		return nil, postgres.CheckReadError(err)
	}

	return row.unmarshal()
}

func (r *TagRepository) List(limit, offset int) ([]*model.Tag, error) {
	stmt := `
		SELECT
			tag.id,
			tag.name,
			tag.meta_created_at,
				tag.meta_updated_at
		FROM tag
		ORDER BY tag.meta_created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.conn.Query(context.Background(), stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	tagRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbTag])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	var tags []*model.Tag
	for _, row := range tagRows {
		tag, err := row.unmarshal()
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepository) Count() (int, error) {
	stmt := `
		SELECT count(*)
		FROM tag`

	rows, err := r.conn.Query(context.Background(), stmt)
	if err != nil {
		return 0, err
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, postgres.CheckReadError(err)
	}

	return count, nil
}

func (r *TagRepository) Delete(tag *model.Tag) error {
	stmt := `
		DELETE FROM tag
		WHERE id = $1
		RETURNING id`

	rows, err := r.conn.Query(context.Background(), stmt, tag.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return postgres.CheckDeleteError(err)
	}

	return nil
}
