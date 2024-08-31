package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/middleware"
)

// Simulates an io.Reader of infinite size.
type infiniteReader struct{}

// Make a new infinite io.Reader.
func newInfiniteReader() *infiniteReader {
	r := infiniteReader{}
	return &r
}

// Always fill the requested buffer with zeros.
func (r *infiniteReader) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		p[i] = 0
	}
	return len(p), nil
}

func TestLimitRequestBodySize(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", newInfiniteReader())

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to read 5KB from the 4KB limit.
		buf := make([]byte, 5*1024)
		_, err := r.Body.Read(buf)
		test.AssertErrorContains(t, err, "request body too large")
	})

	limitRequestBodySize := middleware.LimitRequestBodySize()
	h := limitRequestBodySize(next)
	h.ServeHTTP(w, r)

	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}
