package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

type jsonTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func TestHandleTagList(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(logger, store)

		test.CreateMockTag(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/tags", nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatal(err)
		}

		if rr.StatusCode != 200 {
			t.Fatalf("want %v, got %v", 200, rr.StatusCode)
		}

		var resp map[string][]jsonTag
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["tags"]
		if !ok {
			t.Fatalf("response missing key: %v", "tags")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one tag")
		}

		return test.ErrSkipCommit
	})
}

func TestHandleTagListPagination(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(logger, store)

		// create 5 tags to test with
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)
		test.CreateMockTag(t, store)

		tests := []struct {
			limit int
			want  int
		}{
			{0, 0},
			{1, 1},
			{3, 3},
			{5, 5},
		}

		for _, test := range tests {
			url := fmt.Sprintf("/tags?limit=%d", test.limit)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)

			router := app.Router()
			router.ServeHTTP(w, r)

			rr := w.Result()
			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rr.StatusCode != 200 {
				t.Fatalf("want %v, got %v", 200, rr.StatusCode)
			}

			var resp map[string][]jsonTag
			err = json.Unmarshal(body, &resp)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := resp["tags"]
			if !ok {
				t.Fatalf("response missing key: %v", "tags")
			}

			if len(got) != test.want {
				t.Fatalf("want %v, got %v", test.want, len(got))
			}
		}
		return test.ErrSkipCommit
	})
}
