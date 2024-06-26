package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

// Let's Go Further - Chapter 15.3
func Authenticate(store *storage.Storage) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add the "Vary: Authorization" header to the response. This indicates to any
			// caches that the response may vary based on the value of the Authorization
			// header in the request.
			w.Header().Add("Vary", "Authorization")

			// Retrieve the value of the Authorization header from the request. This will
			// return the empty string "" if there is no such header found.
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Otherwise, we expect the value of the Authorization header to be in the format
			// "Bearer <token>". We try to split this into its constituent parts, and if the
			// header isn't in the expected format we return a 401 Unauthorized response
			// using the invalidAuthenticationTokenResponse() helper (which we will create
			// in a moment).
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				util.UnauthorizedResponse(w, r)
				return
			}

			// Extract the actual authentication token from the header parts.
			token := headerParts[1]

			// Retrieve the details of the account associated with the token.
			account, err := store.Account().ReadByToken(token)
			if err != nil {
				switch {
				case errors.Is(err, postgres.ErrNotFound):
					util.UnauthorizedResponse(w, r)
				default:
					util.ServerErrorResponse(w, r, err)
				}

				return
			}

			// Call the ContextSetAccount() helper to add the account information to the request context.
			r = util.ContextSetAccount(r, account)

			// Call the next handler in the chain.
			next.ServeHTTP(w, r)
		})
	}
}

func AccountRequired() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := util.ContextGetAccount(r)
			if !ok {
				util.UnauthorizedResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AdminRequired() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			account, ok := util.ContextGetAccount(r)
			if !ok {
				util.UnauthorizedResponse(w, r)
				return
			}

			if !account.IsAdmin() {
				util.UnauthorizedResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
