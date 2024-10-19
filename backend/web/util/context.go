package util

import (
	"context"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/model"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Convert the string "account" to a contextKey type and assign it to the contextKeyAccount
// constant. We'll use this constant as the key for getting and setting account information
// in the request context.
const contextKeyAccount = contextKey("account")

// The SetContextAccount() method returns a new copy of the request with the provided
// model.Account added to the context. Note that we use our accountContextKey constant
// as the key.
func SetContextAccount(r *http.Request, account *model.Account) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyAccount, account)
	return r.WithContext(ctx)
}

// The GetContextAccount() retrieves the model.Account and an "exists" bool from the request context.
func GetContextAccount(r *http.Request) (*model.Account, bool) {
	account, ok := r.Context().Value(contextKeyAccount).(*model.Account)
	return account, ok
}

// The app's config will be added to every request context via middleware.
const contextKeyConfig = contextKey("config")

func SetContextConfig(r *http.Request, conf config.Config) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyConfig, conf)
	return r.WithContext(ctx)
}

func GetContextConfig(r *http.Request) (config.Config, bool) {
	conf, ok := r.Context().Value(contextKeyConfig).(config.Config)
	return conf, ok
}
