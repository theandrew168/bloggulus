package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func filterSpans(r *http.Request) bool {
	// Filter out public file requests (static assets, etc).
	publicFilesPatterns := []string{"GET /favicon.ico", "GET /robots.txt", "GET /css/", "GET /img/", "GET /js/"}
	if slices.Contains(publicFilesPatterns, r.Pattern) {
		return false
	}

	// The "/" catch-all pattern accepts any HTTP method + route. Anything
	// in this bucket is a 404 Not Found and can be ignored.
	if r.Pattern == "/" {
		return false
	}

	return true
}

// Normalize the HTTP method to a valid, uppercase method. If the method
// is not valid, return "HTTP" to guard against cardinality explosions.
func normalizeMethod(method string) string {
	validMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	method = strings.ToUpper(method)
	if slices.Contains(validMethods, method) {
		return method
	}

	return "HTTP"
}

// Normalize the HTTP request pattern to a more general form for tracing. Static
// files are grouped under "public-files", the special exact-match index pattern "/{$}"
// is normalized to "/", and the catch-all path "/" is labeled as "not-found".
func normalizeSpanName(r *http.Request) string {
	staticPatterns := []string{"GET /favicon.ico", "GET /robots.txt", "GET /css/", "GET /img/", "GET /js/"}
	if slices.Contains(staticPatterns, r.Pattern) {
		return "GET public-files"
	}

	if r.Pattern == "GET /{$}" {
		return "GET /"
	}

	// The "/" catch-all pattern accepts any HTTP method + route.
	if r.Pattern == "/" {
		return fmt.Sprintf("%s not-found", normalizeMethod(r.Method))
	}

	return r.Pattern
}

func Trace() Middleware {
	return otelhttp.NewMiddleware("", otelhttp.WithFilter(filterSpans), otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
		return normalizeSpanName(r)
	}))
}
