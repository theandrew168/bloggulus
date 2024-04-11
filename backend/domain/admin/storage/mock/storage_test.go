package mock_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/suite"
)

func TestMockStorage(t *testing.T) {
	t.Parallel()

	store := mock.NewStorage()
	suite.RunStorageTests(t, store)
}
