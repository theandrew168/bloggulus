package web_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web"
)

func TestHandleIndex(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	app := web.NewApplication(logger, storage)

	post := test.CreateMockPost(t, storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

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

	page := string(body)
	if !strings.Contains(strings.ToLower(page), strings.ToLower(post.Title)) {
		t.Fatalf("expected recent post title on page")
	}
}

func TestHandleIndexSearch(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, storage)

	// generate some searchable post data
	post := domain.NewPost(
		test.RandomURL(32),
		"python rust",
		test.RandomTime(),
		"content about python and rust",
		blog,
	)

	// create a searchable post
	err := storage.Post.Create(&post)
	if err != nil {
		t.Fatal(err)
	}

	app := web.NewApplication(logger, storage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?q=python+rust", nil)

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

	page := string(body)
	if !strings.Contains(strings.ToLower(page), strings.ToLower(post.Title)) {
		t.Fatalf("expected searched post title on page")
	}
}
