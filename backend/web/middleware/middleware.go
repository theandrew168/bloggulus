package middleware

import (
	"net/http"
	"slices"
)

// Based on:
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
// https://www.youtube.com/watch?v=tIm8UkSf6RA

// Represents a piece of HTTP middleware.
type Middleware func(http.Handler) http.Handler

// Apply a sequence of middleware to a handler (in the provided order).
func Use(h http.Handler, mws ...Middleware) http.Handler {
	// Due to how these functions wrap the handler, we apply them
	// in reverse order so that the first one supplied is the first
	// one that runs.
	for _, mw := range slices.Backward(mws) {
		h = mw(h)
	}
	return h
}

// Chain multiple middleware together for delayed application to a handler.
func Chain(mws ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		return Use(h, mws...)
	}
}
