package middleware

import "net/http"

const (
	MaxRequestBodySize = 4096
)

func LimitRequestBodySize() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

			next.ServeHTTP(w, r)
		})
	}
}
