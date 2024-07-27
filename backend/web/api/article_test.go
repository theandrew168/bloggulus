package api_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
)

type jsonArticle struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	BlogTitle   string    `json:"blogTitle"`
	BlogURL     string    `json:"blogURL"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
}

func TestHandleArticleList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	test.CreatePost(t, store, blog)

	h := api.HandleArticleList(store)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/articles", nil)
	h.ServeHTTP(w, r)

	rr := w.Result()
	respBody, err := io.ReadAll(rr.Body)
	test.AssertNilError(t, err)

	test.AssertEqual(t, rr.StatusCode, 200)

	var resp struct {
		Count    int           `json:"count"`
		Articles []jsonArticle `json:"articles"`
	}
	err = json.Unmarshal(respBody, &resp)
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, resp.Count, 1)
	test.AssertAtLeast(t, len(resp.Articles), 1)
}
