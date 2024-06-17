package middleware

import (
	"net/http"
)

func RecoverPanic() Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")

					code := http.StatusInternalServerError
					text := http.StatusText(code)
					http.Error(w, text, code)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
