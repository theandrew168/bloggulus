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

// Always report that the entire buffer was filled.
func (r *infiniteReader) Read(p []byte) (int, error) {
	return len(p), nil
}

func TestLimitRequestBodySize(t *testing.T) {
	t.Parallel()

	// Prepare the mock ResponseWriter and Request.
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", newInfiniteReader())

	// Prepare the stub handler.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to read more than the request body limit.
		buf := make([]byte, middleware.MaxRequestBodySize+1)
		_, err := r.Body.Read(buf)
		test.AssertErrorAs(t, err, new(*http.MaxBytesError))
	})

	// Wrap the stub handler in the middleware we want to test.
	limitRequestBodySize := middleware.LimitRequestBodySize()
	h := limitRequestBodySize(next)

	// Serve the HTTP request.
	h.ServeHTTP(w, r)

	// Verify that our stub handler was executed.
	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}
