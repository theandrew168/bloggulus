package web

import (
	"io"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

// TODO: Rename p / s to page / size.
func HandleIndexPage(qry *query.Query) http.Handler {
	tmpl := page.NewIndex()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, isLoggedIn := util.GetContextAccount(r)

		// check search param
		search := r.URL.Query().Get("q")

		// check page param
		p, err := strconv.Atoi(r.URL.Query().Get("p"))
		if err != nil {
			p = 1
		}

		if p < 1 {
			p = 1
		}

		// assume size is always 20 (for now...)
		s := 20
		limit, offset := util.PageSizeToLimitOffset(p, s)

		var count int
		var articles []query.Article

		// Two levels of decision making here:
		// 1. Is the user logged in?
		// 2. Is the user searching?
		var g errgroup.Group
		if isLoggedIn {
			if search != "" {
				g.Go(func() error {
					var err error
					count, err = qry.CountRelevantArticlesByAccount(account, search)
					return err
				})
				g.Go(func() error {
					var err error
					articles, err = qry.ListRelevantArticlesByAccount(account, search, limit, offset)
					return err
				})
			} else {
				g.Go(func() error {
					var err error
					count, err = qry.CountRecentArticlesByAccount(account)
					return err
				})
				g.Go(func() error {
					var err error
					articles, err = qry.ListRecentArticlesByAccount(account, limit, offset)
					return err
				})
			}
		} else {
			if search != "" {
				g.Go(func() error {
					var err error
					count, err = qry.CountRelevantArticles(search)
					return err
				})
				g.Go(func() error {
					var err error
					articles, err = qry.ListRelevantArticles(search, limit, offset)
					return err
				})
			} else {
				g.Go(func() error {
					var err error
					count, err = qry.CountRecentArticles()
					return err
				})
				g.Go(func() error {
					var err error
					articles, err = qry.ListRecentArticles(limit, offset)
					return err
				})
			}
		}

		err = g.Wait()
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		data := page.IndexData{
			BaseData: util.GetTemplateBaseData(r, w),

			Search:       search,
			Articles:     articles,
			HasMorePages: p*s < count,
			NextPage:     p + 1,
		}
		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}
