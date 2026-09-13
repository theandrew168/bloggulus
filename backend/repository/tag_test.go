package repository_test

import (
	"context"
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTagCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	tag := test.NewTag()
	err := repo.Tag().Create(context.Background(), tag)
	test.AssertNilError(t, err)
}

func TestTagCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	tag := test.CreateTag(t, repo)

	// attempt to create the same tag again
	err := repo.Tag().Create(context.Background(), tag)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestTagRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	tag := test.CreateTag(t, repo)
	got, err := repo.Tag().Read(context.Background(), tag.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), tag.ID())
}

func TestTagDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	tag := test.CreateTag(t, repo)

	err := repo.Tag().Delete(context.Background(), tag)
	test.AssertNilError(t, err)

	_, err = repo.Tag().Read(context.Background(), tag.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
