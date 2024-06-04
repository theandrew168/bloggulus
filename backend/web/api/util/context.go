package util

import (
	"context"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Convert the string "account" to a contextKey type and assign it to the accountContextKey
// constant. We'll use this constant as the key for getting and setting account information
// in the request context.
const accountContextKey = contextKey("account")

// The ContextSetAccount() method returns a new copy of the request with the provided
// admin.Account added to the context. Note that we use our accountContextKey constant
// as the key.
func ContextSetAccount(r *http.Request, account *admin.Account) *http.Request {
	ctx := context.WithValue(r.Context(), accountContextKey, account)
	return r.WithContext(ctx)
}

// The ContextGetAccount() retrieves the admin.Account from the request context. The only
// time that we'll use this helper is when we logically expect there to be admin.Account
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
// As we discussed earlier in the book, it's OK to panic in those circumstances.
func ContextGetAccount(r *http.Request) *admin.Account {
	account, ok := r.Context().Value(accountContextKey).(*admin.Account)
	if !ok {
		panic("missing account value in request context")
	}
	return account
}
