package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/suite"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

func TestPostgresStorage(t *testing.T) {
	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	suite.RunStorageTests(t, store)
}
