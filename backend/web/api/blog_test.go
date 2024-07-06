package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/feed/mock"
	"github.com/theandrew168/bloggulus/backend/postgres"
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

func TestHandleBlogCreate(t *testing.T) {
	t.Parallel()

	blog := feed.Blog{
		FeedURL: "https://example.com/index.xml",
		SiteURL: "https://example.com",
		Title:   "Example Blog",
	}
	feed, err := mock.GenerateAtomFeed(blog)
	test.AssertNilError(t, err)

	feeds := map[string]string{
		blog.FeedURL: feed,
	}

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		syncService := test.NewSyncService(t, store, feeds, nil)

		h := api.HandleBlogCreate(store, syncService)

		req := struct {
			FeedURL string `json:"feedURL"`
		}{
			FeedURL: blog.FeedURL,
		}

		reqBody, err := json.Marshal(req)
		test.AssertNilError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/blogs", bytes.NewReader(reqBody))
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Blog jsonBlog `json:"blog"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Blog
		test.AssertEqual(t, got.FeedURL, blog.FeedURL)
		test.AssertEqual(t, got.SiteURL, blog.SiteURL)
		test.AssertEqual(t, got.Title, blog.Title)

		// Ensure the blog got created in the database.
		_, err = store.Blog().Read(got.ID)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestHandleBlogRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)

		h := api.HandleBlogRead(store)

		url := fmt.Sprintf("/blogs/%s", blog.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)
		r.SetPathValue("blogID", blog.ID().String())
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Blog jsonBlog `json:"blog"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Blog
		test.AssertEqual(t, got.ID, blog.ID())

		return postgres.ErrRollback
	})
}

func TestHandleBlogReadNotFound(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	h := api.HandleBlogRead(store)

	blogID := uuid.New()
	url := fmt.Sprintf("/blogs/%s", blogID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)
	r.SetPathValue("blogID", blogID.String())
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, 404)
}

func TestHandleBlogList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		test.CreateBlog(t, store)

		h := api.HandleBlogList(store)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/blogs", nil)
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Blogs []jsonBlog `json:"blogs"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Blogs
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

		h := api.HandleBlogList(store)

		for _, tt := range tests {
			url := fmt.Sprintf("/blogs?size=%d", tt.size)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)
			h.ServeHTTP(w, r)

			rr := w.Result()
			respBody, err := io.ReadAll(rr.Body)
			test.AssertNilError(t, err)

			test.AssertEqual(t, rr.StatusCode, 200)

			var resp struct {
				Blogs []jsonBlog `json:"blogs"`
			}
			err = json.Unmarshal(respBody, &resp)
			test.AssertNilError(t, err)

			got := resp.Blogs
			test.AssertEqual(t, len(got), tt.want)
		}
		return postgres.ErrRollback
	})
}

func TestHandleBlogDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)

		h := api.HandleBlogDelete(store)

		url := fmt.Sprintf("/blogs/%s", blog.ID())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", url, nil)
		r.SetPathValue("blogID", blog.ID().String())
		h.ServeHTTP(w, r)

		rr := w.Result()
		respBody, err := io.ReadAll(rr.Body)
		test.AssertNilError(t, err)

		test.AssertEqual(t, rr.StatusCode, 200)

		var resp struct {
			Blog jsonBlog `json:"blog"`
		}
		err = json.Unmarshal(respBody, &resp)
		test.AssertNilError(t, err)

		got := resp.Blog
		test.AssertEqual(t, got.ID, blog.ID())

		_, err = store.Blog().Read(got.ID)
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
