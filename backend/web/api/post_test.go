package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
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
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		post := test.CreateMockPost(t, store)

		url := fmt.Sprintf("/posts/%s", post.ID)
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

		var resp map[string]jsonPost
		err = json.Unmarshal(body, &resp)
		if err != nil {
			t.Fatal(err)
		}

		got, ok := resp["post"]
		if !ok {
			t.Fatalf("response missing key: %v", "post")
		}

		if got.ID != post.ID {
			t.Fatalf("want %v, got %v", post.ID, got.ID)
		}

		return test.ErrSkipCommit
	})
}

func TestHandlePostReadNotFound(t *testing.T) {
	logger := test.NewLogger(t)
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		app := api.NewApplication(logger, store)

		path := fmt.Sprintf("/posts/%s", uuid.New())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", path, nil)

		router := app.Router()
		router.ServeHTTP(w, r)

		rr := w.Result()
		if rr.StatusCode != 404 {
			t.Fatalf("want %v, got %v", 404, rr.StatusCode)
		}

		return test.ErrSkipCommit
	})
}

func TestHandlePostList(t *testing.T) {
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

		rr := w.Result()
		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatal(err)
		}

		if rr.StatusCode != 200 {
			t.Fatalf("want %v, got %v", 200, rr.StatusCode)
		}

		var resp map[string][]domain.Post
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

		return test.ErrSkipCommit
	})
}

func TestHandlePostListPagination(t *testing.T) {
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

			rr := w.Result()
			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rr.StatusCode != 200 {
				t.Fatalf("want %v, got %v", 200, rr.StatusCode)
			}

			var resp map[string][]domain.Post
			err = json.Unmarshal(body, &resp)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := resp["posts"]
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
