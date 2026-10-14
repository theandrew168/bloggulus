package middleware

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// shouldTrace returns false for low-value, high-volume endpoints.
func shouldTrace(r *http.Request) bool {
	switch {
	case r.URL.Path == "/ping":
		return false
	case r.URL.Path == "/metrics":
		return false
	case r.URL.Path == "/favicon.ico":
		return false
	case r.URL.Path == "/robots.txt":
		return false
	case strings.HasPrefix(r.URL.Path, "/css/"):
		return false
	case strings.HasPrefix(r.URL.Path, "/img/"):
		return false
	case strings.HasPrefix(r.URL.Path, "/js/"):
		return false
	default:
		return true
	}
}

// Normalize the HTTP request pattern to a more general form for tracing. The
// special exact-match index pattern "/{$}" is normalized to "/", and the catch-all
// path "/" is labeled as "not-found".
func formatSpanName(operation string, r *http.Request) string {
	if r.Pattern == "GET /{$}" {
		return "GET /"
	}

	// The "/" catch-all pattern accepts any HTTP method + route.
	if r.Pattern == "/" {
		return "Not Found"
	}

	return r.Pattern
}

func Trace() Middleware {
	return otelhttp.NewMiddleware("", otelhttp.WithFilter(shouldTrace), otelhttp.WithSpanNameFormatter(formatSpanName))
}
