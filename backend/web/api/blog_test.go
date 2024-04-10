package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/todo"
	"github.com/theandrew168/bloggulus/backend/testutil"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

type jsonBlog struct {
	ID      uuid.UUID `json:"id"`
	FeedURL string    `json:"feedURL"`
	SiteURL string    `json:"siteURL"`
	Title   string    `json:"title"`
}

func TestHandleBlogRead(t *testing.T) {
	t.Parallel()

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		blog := todo.CreateMockBlog(t, store)

		url := fmt.Sprintf("/blogs/%s", blog.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string]jsonBlog
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["blog"]
		if !ok {
			t.Fatalf("response missing key: %v", "blog")
		}

		testutil.AssertEqual(t, got.ID, blog.ID())

		return storage.ErrRollback
	})
}

func TestHandleBlogReadNotFound(t *testing.T) {
	t.Parallel()

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	app := api.NewApplication(store)

	path := fmt.Sprintf("/blogs/%s", uuid.New())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", path, nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	rr := w.Result()
	testutil.AssertEqual(t, rr.StatusCode, 404)
}

func TestHandleBlogList(t *testing.T) {
	t.Parallel()

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		todo.CreateMockBlog(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/blogs", nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string][]jsonBlog
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["blogs"]
		if !ok {
			t.Fatalf("response missing key: %v", "blogs")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one blog")
		}

		return storage.ErrRollback
	})
}

func TestHandleBlogListPagination(t *testing.T) {
	t.Parallel()

	store, closer := testutil.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		// create 5 blogs to test with
		todo.CreateMockBlog(t, store)
		todo.CreateMockBlog(t, store)
		todo.CreateMockBlog(t, store)
		todo.CreateMockBlog(t, store)
		todo.CreateMockBlog(t, store)

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
			url := fmt.Sprintf("/blogs?limit=%d", tt.limit)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)

			router := app.Router()
			router.ServeHTTP(w, r)

			rr := w.Result()
			body, err := io.ReadAll(rr.Body)
			testutil.AssertNilError(t, err)

			testutil.AssertEqual(t, rr.StatusCode, 200)

			var resp map[string][]jsonBlog
			err = json.Unmarshal(body, &resp)
			testutil.AssertNilError(t, err)

			got, ok := resp["blogs"]
			if !ok {
				t.Fatalf("response missing key: %v", "blogs")
			}

			testutil.AssertEqual(t, len(got), tt.want)
		}
		return storage.ErrRollback
	})
}
