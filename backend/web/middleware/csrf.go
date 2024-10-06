package middleware

import (
	"net/http"

	"github.com/justinas/nosurf"

	"github.com/theandrew168/bloggulus/backend/web/util"
)

func PreventCSRF() Middleware {
	return func(next http.Handler) http.Handler {
		baseCookie := util.NewBaseCookie()
		baseCookie.Name = "bloggulus_csrf_token"

		csrfHandler := nosurf.New(next)
		csrfHandler.SetBaseCookie(baseCookie)
		return csrfHandler
	}
}
