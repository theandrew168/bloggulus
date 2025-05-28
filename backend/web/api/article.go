package api

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/web/api/jsonutil"
)

func HandleArticleList(qry *query.Query) http.Handler {
	type response struct {
		Articles []query.Article `json:"articles"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articles, err := qry.ListRecentArticles(2, 0)
		if err != nil {
			http.Error(w, "Error fetching articles", http.StatusInternalServerError)
			return
		}

		jsonutil.Write(w, http.StatusOK, response{
			Articles: articles,
		})
	})
}
