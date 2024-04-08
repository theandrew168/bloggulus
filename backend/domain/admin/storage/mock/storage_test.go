package mock_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/test"
)

func TestMockStorage(t *testing.T) {
	store := mock.NewStorage()
	test.RunStorageTests(t, store)
}
