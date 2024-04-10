package storage_test

import (
	"testing"

	storageTest "github.com/theandrew168/bloggulus/backend/domain/admin/storage/test"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

func TestPostgresStorage(t *testing.T) {
	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	storageTest.RunStorageTests(t, store)
}
