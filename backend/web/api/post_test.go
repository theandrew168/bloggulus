package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

func TestHandleReadPost(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		post := test.CreateMockPost(t, store)

		url := fmt.Sprintf("/posts/%d", post.ID)
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

		var env map[string]domain.Post
		err = json.Unmarshal(body, &env)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := env["post"]
		if !ok {
			t.Fatalf("response missing key: %v", "post")
		}

		if got.ID != post.ID {
			t.Fatalf("want %v, got %v", post.ID, got.ID)
		}

		return test.ErrSkipCommit
	})
}

func TestHandleReadPostNotFound(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/posts/999999999", nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != 404 {
			t.Fatalf("want %v, got %v", 404, resp.StatusCode)
		}

		return test.ErrSkipCommit
	})
}

func TestHandleReadPosts(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		test.CreateMockPost(t, store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/posts", nil)

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

		var env map[string][]domain.Post
		err = json.Unmarshal(body, &env)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := env["posts"]
		if !ok {
			t.Fatalf("response missing key: %v", "posts")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one blog")
		}

		return test.ErrSkipCommit
	})
}

func TestHandleReadPostsPagination(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		// create 5 posts to test with
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)

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
			url := fmt.Sprintf("/posts?limit=%d", test.limit)
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

			var env map[string][]domain.Post
			err = json.Unmarshal(body, &env)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := env["posts"]
			if !ok {
				t.Fatalf("response missing key: %v", "posts")
			}

			if len(got) != test.want {
				t.Fatalf("want %v, got %v", test.want, len(got))
			}
		}

		return test.ErrSkipCommit
	})

}

func TestHandleReadPostsSearch(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		blog := test.CreateMockBlog(t, store)
		q := "python rust"

		// create searchable post
		post := domain.NewPost(test.RandomURL(32), q, test.RandomTime(), test.RandomString(32), blog)
		err := store.Post.Create(&post)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/posts?q=%s", url.QueryEscape(q))
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

		var env map[string][]domain.Post
		err = json.Unmarshal(body, &env)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := env["posts"]
		if !ok {
			t.Fatalf("response missing key: %v", "posts")
		}

		if len(got) < 1 {
			t.Fatalf("expected at least one matching post")
		}

		return test.ErrSkipCommit
	})
}
