package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
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
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		syncService := test.NewSyncService(t, store, nil, nil)

		app := api.NewApplication(store, syncService)
		handler := app.Handler()

		blog := test.CreateBlog(t, store)
		post := test.CreatePost(t, store, blog)

		url := fmt.Sprintf("/blogs/%s/posts/%s", blog.ID(), post.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string]jsonPost
		err = json.Unmarshal(respBody, &resp)
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

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		syncService := test.NewSyncService(t, store, nil, nil)

		app := api.NewApplication(store, syncService)
		handler := app.Handler()

		blog := test.CreateBlog(t, store)

		path := fmt.Sprintf("/blogs/%s/posts/%s", blog.ID(), uuid.New())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", path, nil)
		handler.ServeHTTP(w, r)

		rr := w.Result()
		test.AssertEqual(t, rr.StatusCode, 404)

		return postgres.ErrRollback
	})
}

func TestHandlePostList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		syncService := test.NewSyncService(t, store, nil, nil)

		app := api.NewApplication(store, syncService)
		handler := app.Handler()

		blog := test.CreateBlog(t, store)
		test.CreatePost(t, store, blog)

		path := fmt.Sprintf("/blogs/%s/posts", blog.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", path, nil)
		handler.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string][]jsonPost
		err = json.Unmarshal(respBody, &resp)
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

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		syncService := test.NewSyncService(t, store, nil, nil)

		app := api.NewApplication(store, syncService)
		handler := app.Handler()

		// create 5 posts to test with
		blog := test.CreateBlog(t, store)
		test.CreatePost(t, store, blog)
		test.CreatePost(t, store, blog)
		test.CreatePost(t, store, blog)
		test.CreatePost(t, store, blog)
		test.CreatePost(t, store, blog)

		tests := []struct {
			size int
			want int
		}{
			{1, 1},
			{3, 3},
			{5, 5},
		}

		for _, tt := range tests {
			url := fmt.Sprintf("/blogs/%s/posts?size=%d", blog.ID(), tt.size)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)
			handler.ServeHTTP(w, r)

			rr := w.Result()
			respBody, err := io.ReadAll(rr.Body)
			test.AssertNilError(t, err)

			test.AssertEqual(t, rr.StatusCode, 200)

			var resp map[string][]jsonPost
			err = json.Unmarshal(respBody, &resp)
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

func TestHandlePostDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		syncService := test.NewSyncService(t, store, nil, nil)

		app := api.NewApplication(store, syncService)
		handler := app.Handler()

		blog := test.CreateBlog(t, store)
		post := test.CreatePost(t, store, blog)

		url := fmt.Sprintf("/blogs/%s/posts/%s", blog.ID(), post.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", url, nil)
		handler.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp map[string]jsonBlog
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got, ok := resp["post"]
		if !ok {
			t.Fatalf("response missing key: %v", "post")
		}

		test.AssertEqual(t, got.ID, post.ID())

		_, err = store.Post().Read(got.ID)
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
