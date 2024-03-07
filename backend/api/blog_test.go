package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/api"
	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestHandleReadBlog(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		blog := test.CreateMockBlog(t, store)

		url := fmt.Sprintf("/blog/%d", blog.ID)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Fatalf("want %v, got %v", 200, resp.StatusCode)
		}

		var env map[string]domain.Blog
		err = json.Unmarshal(body, &env)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := env["blog"]
		if !ok {
			t.Fatalf("response missing key: %v", "blog")
		}

		if got.ID != blog.ID {
			t.Fatalf("want %v, got %v", blog.ID, got.ID)
		}

		return test.ErrSkipCommit
	})
}

func TestHandleReadBlogNotFound(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		return test.ErrSkipCommit
	})

	app := api.NewApplication(logger, store)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/blog/999999999", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != 404 {
		t.Fatalf("want %v, got %v", 404, resp.StatusCode)
	}
}

func TestHandleReadBlogs(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		test.CreateMockBlog(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/blog", nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Fatalf("want %v, got %v", 200, resp.StatusCode)
		}

		var env map[string][]domain.Blog
		err = json.Unmarshal(body, &env)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := env["blogs"]
		if !ok {
			t.Fatalf("response missing key: %v", "blogs")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one blog")
		}

		return test.ErrSkipCommit
	})
}

func TestHandleReadBlogsPagination(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

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

		for _, test := range tests {
			url := fmt.Sprintf("/blog?limit=%d", test.limit)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)

			router := app.Router()
			router.ServeHTTP(w, r)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != 200 {
				t.Fatalf("want %v, got %v", 200, resp.StatusCode)
			}

			var env map[string][]domain.Blog
			err = json.Unmarshal(body, &env)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := env["blogs"]
			if !ok {
				t.Fatalf("response missing key: %v", "blogs")
			}

			if len(got) != test.want {
				t.Fatalf("want %v, got %v", test.want, len(got))
			}
		}
		return test.ErrSkipCommit
	})
}
