package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTagCreate(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		tag := test.NewMockTag()
		err := store.Tag.Create(tag)
		if err != nil {
			t.Fatal(err)
		}

		return test.ErrSkipCommit
	})
}

func TestTagCreateAlreadyExists(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		tag := test.CreateMockTag(t, store)

		// attempt to create the same tag again
		err := store.Tag.Create(tag)
		if !errors.Is(err, storage.ErrConflict) {
			t.Fatal("duplicate tag should return an error")
		}

		return test.ErrSkipCommit
	})
}

func TestTagList(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)

		limit := 3
		offset := 0
		tags, err := store.Tag.List(limit, offset)
		if err != nil {
			t.Fatal(err)
		}

		if len(tags) != limit {
			t.Fatalf("want %v, got %v", limit, len(tags))
		}

		return test.ErrSkipCommit
	})
}

func TestTagDelete(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		tag := test.NewMockTag()
		err := store.Tag.Create(tag)
		if err != nil {
			t.Fatal(err)
		}

		err = store.Tag.Delete(tag)
		if err != nil {
			t.Fatal(err)
		}

		return test.ErrSkipCommit
	})
}
