package middleware

import "net/http"

// Limit the size of the request body to 4KB.
const MaxRequestBodySize = 4 * 1024

// Places an upper limit on the size of every request body.
func LimitRequestBodySize() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

			next.ServeHTTP(w, r)
		})
	}
}
