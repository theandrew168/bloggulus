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

	blog := test.CreateMockBlog(storage, t)

	app := api.NewApplication(storage, logger)

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
		t.Errorf("want %v, got %v", 200, resp.StatusCode)
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
		t.Errorf("want %v, got %v", blog.ID, got.ID)
	}
}

func TestHandleReadBlogs(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	logger := test.NewLogger()

	test.CreateMockBlog(storage, t)

	app := api.NewApplication(storage, logger)

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
		t.Errorf("want %v, got %v", 200, resp.StatusCode)
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
		t.Errorf("expected at least one blog")
	}
}

// TODO
func TestHandleReadBlogsPagination(t *testing.T) {

}
