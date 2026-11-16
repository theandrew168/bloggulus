package api

import (
	"net/http"

	webquery "github.com/theandrew168/bloggulus/backend/query/web"
)

func Handler(
	qry *webquery.Query,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /{$}", HandleIndex())
	mux.Handle("GET /openapi.yaml", HandleOpenAPISpec())
	mux.Handle("GET /articles", HandleArticles(qry))
	return mux
}
