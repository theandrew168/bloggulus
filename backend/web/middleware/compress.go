package middleware

import (
	"net/http"

	"github.com/klauspost/compress/gzhttp"
)

func CompressFiles() Middleware {
	return func(next http.Handler) http.Handler {
		return gzhttp.GzipHandler(next)
	}
}
