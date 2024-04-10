package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/todo"
	"github.com/theandrew168/bloggulus/backend/testutil"
	"github.com/theandrew168/bloggulus/backend/web/api"
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

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		post := todo.CreateMockPost(t, store)

		url := fmt.Sprintf("/posts/%s", post.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string]jsonPost
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["post"]
		if !ok {
			t.Fatalf("response missing key: %v", "post")
		}

		testutil.AssertEqual(t, got.ID, post.ID())

		return storage.ErrRollback
	})
}

func TestHandlePostReadNotFound(t *testing.T) {
	t.Parallel()

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		path := fmt.Sprintf("/posts/%s", uuid.New())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", path, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		testutil.AssertEqual(t, rr.StatusCode, 404)

		return storage.ErrRollback
	})
}

func TestHandlePostList(t *testing.T) {
	t.Parallel()

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		todo.CreateMockPost(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/posts", nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string][]jsonPost
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["posts"]
		if !ok {
			t.Fatalf("response missing key: %v", "posts")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one blog")
		}

		return storage.ErrRollback
	})
}

func TestHandlePostListPagination(t *testing.T) {
	t.Parallel()

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		// create 5 posts to test with
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)

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
			testutil.AssertNilError(t, err)

			testutil.AssertEqual(t, rr.StatusCode, 200)

			var resp map[string][]jsonPost
			err = json.Unmarshal(body, &resp)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := resp["posts"]
			if !ok {
				t.Fatalf("response missing key: %v", "posts")
			}

			testutil.AssertEqual(t, len(got), tt.want)
		}

		return storage.ErrRollback
	})
}
