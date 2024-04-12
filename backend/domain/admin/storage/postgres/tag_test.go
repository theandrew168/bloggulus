package postgres_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	storageMock "github.com/theandrew168/bloggulus/backend/domain/admin/storage/mock"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTagCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		tag := mock.NewTag()
		err := store.Tag().Create(tag)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestTagCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		tag := storageMock.CreateTag(t, store)

		// attempt to create the same tag again
		err := store.Tag().Create(tag)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestTagRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		tag := storageMock.CreateTag(t, store)
		got, err := store.Tag().Read(tag.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), tag.ID())

		return postgres.ErrRollback
	})
}

func TestTagList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		storageMock.CreateTag(t, store)
		storageMock.CreateTag(t, store)
		storageMock.CreateTag(t, store)
		storageMock.CreateTag(t, store)
		storageMock.CreateTag(t, store)

		limit := 5
		offset := 0
		tags, err := store.Tag().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(tags), limit)

		return postgres.ErrRollback
	})
}

func TestTagDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		tag := storageMock.CreateTag(t, store)

		err := store.Tag().Delete(tag)
		test.AssertNilError(t, err)

		_, err = store.Tag().Read(tag.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
