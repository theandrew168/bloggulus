package admin_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTagCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		tag := test.NewTag()
		err := store.Admin().Tag().Create(tag)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestTagCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		tag := test.CreateTag(t, store)

		// attempt to create the same tag again
		err := store.Admin().Tag().Create(tag)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestTagRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		tag := test.CreateTag(t, store)
		got, err := store.Admin().Tag().Read(tag.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), tag.ID())

		return postgres.ErrRollback
	})
}

func TestTagList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		test.CreateTag(t, store)
		test.CreateTag(t, store)
		test.CreateTag(t, store)
		test.CreateTag(t, store)
		test.CreateTag(t, store)

		limit := 5
		offset := 0
		tags, err := store.Admin().Tag().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(tags), limit)

		return postgres.ErrRollback
	})
}

func TestTagDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		tag := test.CreateTag(t, store)

		err := store.Admin().Tag().Delete(tag)
		test.AssertNilError(t, err)

		_, err = store.Admin().Tag().Read(tag.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
