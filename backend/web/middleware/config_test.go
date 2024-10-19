package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func TestAddConfig(t *testing.T) {
	t.Parallel()

	conf := config.Config{
		Port: "12345",
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := util.GetContextConfig(r)
		test.AssertEqual(t, ok, true)
		test.AssertEqual(t, got.Port, conf.Port)
	})

	addConfig := middleware.AddConfig(conf)
	h := addConfig(next)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}
