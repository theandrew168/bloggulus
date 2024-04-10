package suite

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	storageTest "github.com/theandrew168/bloggulus/backend/domain/admin/storage/test"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTagCreate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := mock.NewTag()
		err := store.Tag().Create(tag)
		test.AssertNilError(t, err)

		return storage.ErrRollback
	})
}

func TestTagCreateAlreadyExists(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := storageTest.CreateMockTag(t, store)

		// attempt to create the same tag again
		err := store.Tag().Create(tag)
		test.AssertErrorIs(t, err, storage.ErrConflict)

		return storage.ErrRollback
	})
}

func TestTagRead(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := storageTest.CreateMockTag(t, store)
		got, err := store.Tag().Read(tag.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), tag.ID())

		return storage.ErrRollback
	})
}

func TestTagList(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		storageTest.CreateMockTag(t, store)
		storageTest.CreateMockTag(t, store)
		storageTest.CreateMockTag(t, store)
		storageTest.CreateMockTag(t, store)
		storageTest.CreateMockTag(t, store)

		limit := 3
		offset := 0
		tags, err := store.Tag().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(tags), limit)

		return storage.ErrRollback
	})
}

func TestTagDelete(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		tag := storageTest.CreateMockTag(t, store)

		err := store.Tag().Delete(tag)
		test.AssertNilError(t, err)

		_, err = store.Tag().Read(tag.ID())
		test.AssertErrorIs(t, err, storage.ErrNotFound)

		return storage.ErrRollback
	})
}
