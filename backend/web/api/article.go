package api

import (
	"net/http"

	"golang.org/x/sync/errgroup"

	webquery "github.com/theandrew168/bloggulus/backend/query/web"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleArticles(qry *webquery.Query) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ParsePage(r.URL.Query())
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		size, err := ParseSize(r.URL.Query())
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		limit, offset := util.PageSizeToLimitOffset(page, size)

		var total int
		var articles []webquery.Article

		var g errgroup.Group

		search := r.URL.Query().Get("q")
		if search != "" {
			g.Go(func() error {
				var err error
				total, err = qry.Article().CountRelevant(r.Context(), search)
				return err
			})
			g.Go(func() error {
				var err error
				articles, err = qry.Article().ListRelevant(r.Context(), search, limit, offset)
				return err
			})
		} else {
			g.Go(func() error {
				var err error
				total, err = qry.Article().CountRecent(r.Context())
				return err
			})
			g.Go(func() error {
				var err error
				articles, err = qry.Article().ListRecent(r.Context(), limit, offset)
				return err
			})
		}

		err = g.Wait()
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		WriteJSON(w, http.StatusOK, ItemsResponse[webquery.Article]{
			Items: articles,
			Total: total,
			Page:  page,
			Size:  size,
		})
	})
}
