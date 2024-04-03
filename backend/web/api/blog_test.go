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

type jsonBlog struct {
	ID      uuid.UUID `json:"id"`
	FeedURL string    `json:"feedURL"`
	SiteURL string    `json:"siteURL"`
	Title   string    `json:"title"`
}

func TestHandleBlogRead(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		blog := test.CreateMockBlog(t, store)

		url := fmt.Sprintf("/blogs/%s", blog.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string]jsonBlog
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["blog"]
		if !ok {
			t.Fatalf("response missing key: %v", "blog")
		}

		test.AssertEqual(t, got.ID, blog.ID())

		return test.ErrRollback
	})
}

func TestHandleBlogReadNotFound(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	app := api.NewApplication(store)

	path := fmt.Sprintf("/blogs/%s", uuid.New())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", path, nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, 404)
}

func TestHandleBlogList(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		test.CreateMockBlog(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/blogs", nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

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

		return test.ErrRollback
	})
}

func TestHandleBlogListPagination(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		app := api.NewApplication(store)

		// create 5 blogs to test with
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)

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
			test.AssertNilError(t, err)

			test.AssertEqual(t, rr.StatusCode, 200)

			var resp map[string][]jsonBlog
			err = json.Unmarshal(body, &resp)
			test.AssertNilError(t, err)

			got, ok := resp["blogs"]
			if !ok {
				t.Fatalf("response missing key: %v", "blogs")
			}

			test.AssertEqual(t, len(got), tt.want)
		}
		return test.ErrRollback
	})
}