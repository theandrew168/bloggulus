package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/storage"
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
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		blog := test.CreateMockBlog(t, store)

		url := fmt.Sprintf("/blogs/%s", blog.ID)
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

		var resp map[string]domain.Blog
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["blog"]
		if !ok {
			t.Fatalf("response missing key: %v", "blog")
		}

		if got.ID != blog.ID {
			t.Fatalf("want %v, got %v", blog.ID, got.ID)
		}

		return test.ErrSkipCommit
	})
}

func TestHandleBlogReadNotFound(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		return test.ErrSkipCommit
	})

	app := api.NewApplication(logger, store)

	path := fmt.Sprintf("/blogs/%s", uuid.New())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", path, nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	rr := w.Result()
	if rr.StatusCode != 404 {
		t.Fatalf("want %v, got %v", 404, rr.StatusCode)
	}
}

func TestHandleBlogList(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		test.CreateMockBlog(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/blogs", nil)

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

		var resp map[string][]domain.Blog
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

		return test.ErrSkipCommit
	})
}

func TestHandleBlogListPagination(t *testing.T) {
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
			url := fmt.Sprintf("/blogs?limit=%d", test.limit)
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

			var resp map[string][]domain.Blog
			err = json.Unmarshal(body, &resp)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := resp["blogs"]
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
