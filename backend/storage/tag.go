package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/domain"
)

// ensure TagStorage interface is satisfied
var _ TagStorage = (*PostgresTagStorage)(nil)

type TagStorage interface {
	Create(tag domain.Tag) error
	List(limit, offset int) ([]domain.Tag, error)
	Delete(tag domain.Tag) error
}

type dbTag struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostgresTagStorage struct {
	conn database.Conn
}

func NewPostgresTagStorage(conn database.Conn) *PostgresTagStorage {
	s := PostgresTagStorage{
		conn: conn,
	}
	return &s
}

func (s *PostgresTagStorage) marshal(tag domain.Tag) (dbTag, error) {
	row := dbTag{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
	return row, nil
}

func (s *PostgresTagStorage) unmarshal(row dbTag) (domain.Tag, error) {
	post := domain.Tag{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return post, nil
}

func (s *PostgresTagStorage) Create(tag domain.Tag) error {
	stmt := `
		INSERT INTO tag
			(id, name, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)`

	row, err := s.marshal(tag)
	if err != nil {
		return err
	}

	args := []interface{}{
		row.ID,
		row.Name,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	_, err = s.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return checkCreateError(err)
	}

	return nil
}

func (s *PostgresTagStorage) List(limit, offset int) ([]domain.Tag, error) {
	stmt := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM tag
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}

	tagRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbTag])
	if err != nil {
		return nil, checkListError(err)
	}

	var tags []domain.Tag
	for _, row := range tagRows {
		tag, err := s.unmarshal(row)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (repo *PostgresTagStorage) Delete(tag domain.Tag) error {
	stmt := `
		DELETE FROM tag
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, tag.ID)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return checkDeleteError(err)
	}

	return nil
}
