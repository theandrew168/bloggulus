package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestHandleReadBlog(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := api.NewApplication(storage, logger)

	blog := test.CreateMockBlog(storage, t)

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

	var env map[string]core.Blog
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
}

func TestHandleReadBlogNotFound(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := api.NewApplication(storage, logger)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/blog/999999999", nil)

	router := app.Router()
	router.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != 404 {
		t.Fatalf("want %v, got %v", 404 , resp.StatusCode)
	}
}

func TestHandleReadBlogs(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := api.NewApplication(storage, logger)

	test.CreateMockBlog(storage, t)

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

	var env map[string][]core.Blog
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
}

func TestHandleReadBlogsPagination(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()
	app := api.NewApplication(storage, logger)

	// create 5 blogs to test with
	test.CreateMockBlog(storage, t)
	test.CreateMockBlog(storage, t)
	test.CreateMockBlog(storage, t)
	test.CreateMockBlog(storage, t)
	test.CreateMockBlog(storage, t)

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

		var env map[string][]core.Blog
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
}
