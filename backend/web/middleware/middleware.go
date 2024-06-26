package middleware

import (
	"net/http"
)

// Based on:
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
// https://www.youtube.com/watch?v=tIm8UkSf6RA

// Represents a piece of HTTP middleware.
type Middleware func(http.Handler) http.Handler

// Apply a sequence of middleware to a handler (in the provided order).
func Use(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Chain multiple middleware together (in the provided order).
func Chain(mws ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			h = mws[i](h)
		}
		return h
	}
}
