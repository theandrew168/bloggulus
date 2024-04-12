package admin_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	storageMock "github.com/theandrew168/bloggulus/backend/domain/admin/storage/mock"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
	api "github.com/theandrew168/bloggulus/backend/web/api/admin"
)

type jsonPost struct {
	ID          uuid.UUID `json:"id"`
	BlogID      uuid.UUID `json:"blogID"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"publishedAt"`
}

func TestHandlePostRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		post := storageMock.CreateMockPost(t, store)

		url := fmt.Sprintf("/posts/%s", post.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string]jsonPost
		err = json.Unmarshal(body, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["post"]
		if !ok {
			t.Fatalf("response missing key: %v", "post")
		}

		test.AssertEqual(t, got.ID, post.ID())

		return postgres.ErrRollback
	})
}

func TestHandlePostReadNotFound(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		path := fmt.Sprintf("/posts/%s", uuid.New())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", path, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		test.AssertEqual(t, rr.StatusCode, 404)

		return postgres.ErrRollback
	})
}

func TestHandlePostList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		storageMock.CreateMockPost(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/posts", nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string][]jsonPost
		err = json.Unmarshal(body, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["posts"]
		if !ok {
			t.Fatalf("response missing key: %v", "posts")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one blog")
		}

		return postgres.ErrRollback
	})
}

func TestHandlePostListPagination(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		// create 5 posts to test with
		storageMock.CreateMockPost(t, store)
		storageMock.CreateMockPost(t, store)
		storageMock.CreateMockPost(t, store)
		storageMock.CreateMockPost(t, store)
		storageMock.CreateMockPost(t, store)

		tests := []struct {
			limit int
			want  int
		}{
			{0, 0},
			{1, 1},
			{3, 3},
			{5, 5},
		}

		for _, tt := range tests {
			url := fmt.Sprintf("/posts?limit=%d", tt.limit)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)

			router := app.Router()
			router.ServeHTTP(w, r)

			rr := w.Result()
			body, err := io.ReadAll(rr.Body)
			test.AssertNilError(t, err)

			test.AssertEqual(t, rr.StatusCode, 200)

			var resp map[string][]jsonPost
			err = json.Unmarshal(body, &resp)
			test.AssertNilError(t, err)

			got, ok := resp["posts"]
			if !ok {
				t.Fatalf("response missing key: %v", "posts")
			}

			test.AssertEqual(t, len(got), tt.want)
		}

		return postgres.ErrRollback
	})
}
