package middleware

import (
	"net/http"
)

// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
// https://www.youtube.com/watch?v=tIm8UkSf6RA

// Adapter type represents a piece of HTTP middleware.
type Adapter func(http.Handler) http.Handler

// Apply a sequence of middleware to a handler (in the given order).
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}
