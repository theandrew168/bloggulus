package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/suite"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostgresStorage(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	suite.RunStorageTests(t, store)
}
