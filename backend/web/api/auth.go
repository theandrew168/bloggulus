package api

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/theandrew168/bloggulus/backend/random"
	"github.com/theandrew168/bloggulus/backend/web/api/jsonutil"
)

func HandleOAuthSignin(conf *oauth2.Config) http.Handler {
	type response struct {
		State       string `json:"state"`
		AuthCodeURL string `json:"authCodeURL"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := random.BytesBase64(16)
		if err != nil {
			panic(err)
		}

		authCodeURL := conf.AuthCodeURL(state)

		jsonutil.Write(w, http.StatusOK, response{
			State:       state,
			AuthCodeURL: authCodeURL,
		})
	})
}
