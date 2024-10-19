package middleware

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func AddConfig(conf config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = util.SetContextConfig(r, conf)

			next.ServeHTTP(w, r)
		})
	}
}
