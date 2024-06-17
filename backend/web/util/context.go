package util

import (
	"context"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/model/admin"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Convert the string "account" to a contextKey type and assign it to the contextKeyAccount
// constant. We'll use this constant as the key for getting and setting account information
// in the request context.
const contextKeyAccount = contextKey("account")

// The ContextSetAccount() method returns a new copy of the request with the provided
// admin.Account added to the context. Note that we use our accountContextKey constant
// as the key.
func ContextSetAccount(r *http.Request, account *admin.Account) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyAccount, account)
	return r.WithContext(ctx)
}

// The ContextGetAccount() retrieves the admin.Account and an "exists" bool from the request context.
func ContextGetAccount(r *http.Request) (*admin.Account, bool) {
	account, ok := r.Context().Value(contextKeyAccount).(*admin.Account)
	return account, ok
}
