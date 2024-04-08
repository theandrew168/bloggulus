package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTagCreate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := test.NewMockTag()
		err := store.Tag().Create(tag)
		test.AssertNilError(t, err)

		return test.ErrRollback
	})
}

func TestTagCreateAlreadyExists(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := test.CreateMockTag(t, store)

		// attempt to create the same tag again
		err := store.Tag().Create(tag)
		test.AssertErrorIs(t, err, storage.ErrConflict)

		return test.ErrRollback
	})
}

func TestTagRead(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := test.CreateMockTag(t, store)
		got, err := store.Tag().Read(tag.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), tag.ID())

		return test.ErrRollback
	})
}

func TestTagList(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)

		limit := 3
		offset := 0
		tags, err := store.Tag().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(tags), limit)

		return test.ErrRollback
	})
}

func TestTagDelete(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := test.CreateMockTag(t, store)

		err := store.Tag().Delete(tag)
		test.AssertNilError(t, err)

		_, err = store.Tag().Read(tag.ID())
		test.AssertErrorIs(t, err, storage.ErrNotFound)

		return test.ErrRollback
	})
}
