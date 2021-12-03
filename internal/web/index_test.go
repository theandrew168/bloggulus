package web_test

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
	"github.com/theandrew168/bloggulus/internal/web"
)

func TestHandleIndex(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()

	post := test.CreateMockPost(storage, t)

	app := web.NewApplication(storage, logger)

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
		t.Errorf("want %v, got %v", 200, resp.StatusCode)
	}

	page := string(body)
	if !strings.Contains(strings.ToLower(page), strings.ToLower(post.Title)) {
		t.Errorf("expected recent post title on page")
	}
}

func TestHandleIndexSearch(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()

	blog := test.CreateMockBlog(storage, t)

	// generate some searchable post data
	post := core.NewPost(
		test.RandomURL(32),
		"python rust",
		test.RandomTime(),
		blog,
	)

	// create a searchable post
	err := storage.CreatePost(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	app := web.NewApplication(storage, logger)

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
		t.Errorf("want %v, got %v", 200, resp.StatusCode)
	}

	page := string(body)
	if !strings.Contains(strings.ToLower(page), strings.ToLower(post.Title)) {
		t.Errorf("expected searched post title on page")
	}
}
