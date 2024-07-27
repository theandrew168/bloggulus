package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestTagCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	tag := test.NewTag(t)
	err := store.Tag().Create(tag)
	test.AssertNilError(t, err)
}

func TestTagCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	tag := test.CreateTag(t, store)

	// attempt to create the same tag again
	err := store.Tag().Create(tag)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestTagRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	tag := test.CreateTag(t, store)
	got, err := store.Tag().Read(tag.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), tag.ID())
}

func TestTagList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	test.CreateTag(t, store)
	test.CreateTag(t, store)
	test.CreateTag(t, store)

	limit := 3
	offset := 0
	tags, err := store.Tag().List(limit, offset)
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, len(tags), limit)
}

func TestTagCount(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	test.CreateTag(t, store)
	test.CreateTag(t, store)
	test.CreateTag(t, store)

	count, err := store.Tag().Count()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, count, 3)
}

func TestTagDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	tag := test.CreateTag(t, store)

	err := store.Tag().Delete(tag)
	test.AssertNilError(t, err)

	_, err = store.Tag().Read(tag.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
