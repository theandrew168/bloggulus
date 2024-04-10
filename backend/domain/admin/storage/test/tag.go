package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/todo"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

func TestTagCreate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := mock.NewTag()
		err := store.Tag().Create(tag)
		testutil.AssertNilError(t, err)

		return storage.ErrRollback
	})
}

func TestTagCreateAlreadyExists(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := todo.CreateMockTag(t, store)

		// attempt to create the same tag again
		err := store.Tag().Create(tag)
		testutil.AssertErrorIs(t, err, storage.ErrConflict)

		return storage.ErrRollback
	})
}

func TestTagRead(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := todo.CreateMockTag(t, store)
		got, err := store.Tag().Read(tag.ID())
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, got.ID(), tag.ID())

		return storage.ErrRollback
	})
}

func TestTagList(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		todo.CreateMockTag(t, store)
		todo.CreateMockTag(t, store)
		todo.CreateMockTag(t, store)
		todo.CreateMockTag(t, store)
		todo.CreateMockTag(t, store)

		limit := 3
		offset := 0
		tags, err := store.Tag().List(limit, offset)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, len(tags), limit)

		return storage.ErrRollback
	})
}

func TestTagDelete(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := todo.CreateMockTag(t, store)

		err := store.Tag().Delete(tag)
		testutil.AssertNilError(t, err)

		_, err = store.Tag().Read(tag.ID())
		testutil.AssertErrorIs(t, err, storage.ErrNotFound)

		return storage.ErrRollback
	})
}
