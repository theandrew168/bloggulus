package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/api"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func TestHandleBlogFollow(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	h := api.HandleBlogFollow(store)

	account, _ := test.CreateAccount(t, store)
	blog := test.CreateBlog(t, store)

	url := fmt.Sprintf("/blogs/%s/follow", blog.ID())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", url, nil)
	r.SetPathValue("blogID", blog.ID().String())
	h.ServeHTTP(w, util.ContextSetAccount(r, account))

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusNoContent)
}

func TestHandleBlogUnfollow(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	h := api.HandleBlogUnfollow(store)

	account, _ := test.CreateAccount(t, store)
	blog := test.CreateBlog(t, store)

	err := store.AccountBlog().Create(account, blog)
	test.AssertNilError(t, err)

	url := fmt.Sprintf("/blogs/%s/unfollow", blog.ID())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", url, nil)
	r.SetPathValue("blogID", blog.ID().String())
	h.ServeHTTP(w, util.ContextSetAccount(r, account))

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusNoContent)
}

func TestHandleBlogFollowing(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	h := api.HandleBlogFollowing(store)

	account, _ := test.CreateAccount(t, store)
	blog := test.CreateBlog(t, store)

	url := fmt.Sprintf("/blogs/%s/following", blog.ID())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)
	r.SetPathValue("blogID", blog.ID().String())
	h.ServeHTTP(w, util.ContextSetAccount(r, account))

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusNotFound)

	test.CreateAccountBlog(t, store, account, blog)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", url, nil)
	r.SetPathValue("blogID", blog.ID().String())
	h.ServeHTTP(w, util.ContextSetAccount(r, account))

	rr = w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusNoContent)
}
