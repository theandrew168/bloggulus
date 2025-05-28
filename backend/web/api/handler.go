package api

import (
	"io/fs"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/query"
)

func Handler(
	public fs.FS,
	cmd *command.Command,
	qry *query.Query,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /articles", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Articles API Endpoint"))
	}))
	return mux
}
