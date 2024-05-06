package admin_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api/admin"
)

type jsonTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func TestHandleTagList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)
		router := app.Router()

		test.CreateTag(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/tags", nil)
		router.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string][]jsonTag
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["tags"]
		if !ok {
			t.Fatalf("response missing key: %v", "tags")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one tag")
		}

		return postgres.ErrRollback
	})
}

func TestHandleTagListPagination(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)
		router := app.Router()

		// create 5 tags to test with
		test.CreateTag(t, store)
		test.CreateTag(t, store)
		test.CreateTag(t, store)
		test.CreateTag(t, store)
		test.CreateTag(t, store)

		tests := []struct {
			size int
			want int
		}{
			{1, 1},
			{3, 3},
			{5, 5},
		}

		for _, tt := range tests {
			url := fmt.Sprintf("/tags?size=%d", tt.size)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)
			router.ServeHTTP(w, r)

			rr := w.Result()
			respBody, err := io.ReadAll(rr.Body)
			test.AssertNilError(t, err)

			test.AssertEqual(t, rr.StatusCode, 200)

			var resp map[string][]jsonTag
			err = json.Unmarshal(respBody, &resp)
			test.AssertNilError(t, err)

			got, ok := resp["tags"]
			if !ok {
				t.Fatalf("response missing key: %v", "tags")
			}

			test.AssertEqual(t, len(got), tt.want)
		}
		return postgres.ErrRollback
	})
}
