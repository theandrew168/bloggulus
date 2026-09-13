package middleware

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Trace() Middleware {
	return otelhttp.NewMiddleware("")
}
