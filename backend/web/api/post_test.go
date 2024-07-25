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
		h := api.HandlePostRead(store)

		blog := test.CreateBlog(t, store)
		post := test.CreatePost(t, store, blog)

		url := fmt.Sprintf("/blogs/%s/posts/%s", blog.ID(), post.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)
		r.SetPathValue("blogID", blog.ID().String())
		r.SetPathValue("postID", post.ID().String())
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Post jsonPost `json:"post"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Post
		test.AssertEqual(t, got.ID, post.ID())

		return postgres.ErrRollback
	})
}

func TestHandlePostReadNotFound(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)

		h := api.HandlePostRead(store)

		postID := uuid.New()
		url := fmt.Sprintf("/blogs/%s/posts/%s", blog.ID(), postID)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)
		r.SetPathValue("blogID", blog.ID().String())
		r.SetPathValue("postID", postID.String())
		h.ServeHTTP(w, r)

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
		blog := test.CreateBlog(t, store)
		test.CreatePost(t, store, blog)

		h := api.HandlePostList(store)

		url := fmt.Sprintf("/blogs/%s/posts", blog.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)
		r.SetPathValue("blogID", blog.ID().String())
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Count int        `json:"count"`
			Posts []jsonPost `json:"posts"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		if resp.Count < 1 || len(resp.Posts) < 1 {
			t.Fatalf("expected at least one post")
		}

		return postgres.ErrRollback
	})
}

func TestHandlePostListPagination(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
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

		h := api.HandlePostList(store)

		for _, tt := range tests {
			url := fmt.Sprintf("/blogs/%s/posts?size=%d", blog.ID(), tt.size)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)
			r.SetPathValue("blogID", blog.ID().String())
			h.ServeHTTP(w, r)

			rr := w.Result()
			respBody, err := io.ReadAll(rr.Body)
			test.AssertNilError(t, err)

			test.AssertEqual(t, rr.StatusCode, 200)

			var resp struct {
				Posts []jsonPost `json:"posts"`
			}
			err = json.Unmarshal(respBody, &resp)
			test.AssertNilError(t, err)

			got := resp.Posts
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
		blog := test.CreateBlog(t, store)
		post := test.CreatePost(t, store, blog)

		h := api.HandlePostDelete(store)

		url := fmt.Sprintf("/blogs/%s/posts/%s", blog.ID(), post.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", url, nil)
		r.SetPathValue("blogID", blog.ID().String())
		r.SetPathValue("postID", post.ID().String())
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Post jsonPost `json:"post"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Post
		test.AssertEqual(t, got.ID, post.ID())

		_, err = store.Post().Read(got.ID)
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
