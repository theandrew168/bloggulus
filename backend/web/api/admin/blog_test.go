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

type jsonBlog struct {
	ID      uuid.UUID `json:"id"`
	FeedURL string    `json:"feedURL"`
	SiteURL string    `json:"siteURL"`
	Title   string    `json:"title"`
}

func TestHandleBlogRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)
		router := app.Router()

		blog := test.CreateBlog(t, store)

		url := fmt.Sprintf("/blogs/%s", blog.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)
		router.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string]jsonBlog
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["blog"]
		if !ok {
			t.Fatalf("response missing key: %v", "blog")
		}

		test.AssertEqual(t, got.ID, blog.ID())

		return postgres.ErrRollback
	})
}

func TestHandleBlogReadNotFound(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	app := admin.NewApplication(store)
	router := app.Router()

	path := fmt.Sprintf("/blogs/%s", uuid.New())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", path, nil)
	router.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, 404)
}

func TestHandleBlogList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)
		router := app.Router()

		test.CreateBlog(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/blogs", nil)
		router.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string][]jsonBlog
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["blogs"]
		if !ok {
			t.Fatalf("response missing key: %v", "blogs")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one blog")
		}

		return postgres.ErrRollback
	})
}

func TestHandleBlogListPagination(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := admin.NewApplication(store)
		router := app.Router()

		// create 5 blogs to test with
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)

		tests := []struct {
			size int
			want int
		}{
			{1, 1},
			{3, 3},
			{5, 5},
		}

		for _, tt := range tests {
			url := fmt.Sprintf("/blogs?size=%d", tt.size)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)
			router.ServeHTTP(w, r)

			rr := w.Result()
			respBody, err := io.ReadAll(rr.Body)
			test.AssertNilError(t, err)

			test.AssertEqual(t, rr.StatusCode, 200)

			var resp map[string][]jsonBlog
			err = json.Unmarshal(respBody, &resp)
			test.AssertNilError(t, err)

			got, ok := resp["blogs"]
			if !ok {
				t.Fatalf("response missing key: %v", "blogs")
			}

			test.AssertEqual(t, len(got), tt.want)
		}
		return postgres.ErrRollback
	})
}
