package api

import (
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/query"
)

func Handler(
	conf config.Config,
	cmd *command.Command,
	qry *query.Query,
) http.Handler {
	githubConf := oauth2.Config{
		Endpoint:     github.Endpoint,
		ClientID:     conf.GithubClientID,
		ClientSecret: conf.GithubClientSecret,
		RedirectURL:  conf.GithubRedirectURI,
		Scopes:       []string{},
	}
	googleConf := oauth2.Config{
		Endpoint:     google.Endpoint,
		ClientID:     conf.GoogleClientID,
		ClientSecret: conf.GoogleClientSecret,
		RedirectURL:  conf.GoogleRedirectURI,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
	}

	mux := http.NewServeMux()
	mux.Handle("GET /articles", HandleArticleList(qry))
	mux.Handle("GET /signin/github", HandleOAuthSignin(&githubConf))
	mux.Handle("GET /signin/google", HandleOAuthSignin(&googleConf))
	return mux
}
