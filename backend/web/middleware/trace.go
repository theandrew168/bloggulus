package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type spanNameFormatter func(operation string, r *http.Request) string

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

func newSpanNameFormatter(prefix string) spanNameFormatter {
	// Normalize the HTTP request pattern to a more general form for tracing. The
	// special exact-match index pattern "/{$}" is normalized to "/", and the catch-all
	// path "/" is labeled as "not-found".
	return func(operation string, r *http.Request) string {
		if r.Pattern == "GET /{$}" {
			return fmt.Sprintf("GET %s/", prefix)
		}

		// The "/" catch-all pattern accepts any HTTP method + route.
		if r.Pattern == "/" {
			return "Not Found"
		}

		// Split "GET /foo" into ["GET", "/foo"] and prepend the prefix to the path.
		parts := strings.Split(r.Pattern, " ")
		if len(parts) == 2 {
			return fmt.Sprintf("%s %s%s", parts[0], prefix, parts[1])
		}

		return r.Pattern
	}
}

func Trace(prefix string) Middleware {
	formatter := newSpanNameFormatter(prefix)
	return otelhttp.NewMiddleware("", otelhttp.WithFilter(shouldTrace), otelhttp.WithSpanNameFormatter(formatter))
}
