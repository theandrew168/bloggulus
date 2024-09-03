package web

import (
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleIndexPage(find *finder.Finder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := page.Index()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// check search param
		search := r.URL.Query().Get("q")

		// check p param
		p, err := strconv.Atoi(r.URL.Query().Get("p"))
		if err != nil {
			p = 1
		}

		if p < 1 {
			p = 1
		}

		s := 20
		limit, offset := util.PageSizeToLimitOffset(p, s)

		var count int
		var articles []finder.Article

		var g errgroup.Group
		if search != "" {
			g.Go(func() error {
				var err error
				count, err = find.CountSearchArticles(search)
				return err
			})
			g.Go(func() error {
				var err error
				articles, err = find.SearchArticles(search, limit, offset)
				return err
			})
		} else {
			g.Go(func() error {
				var err error
				count, err = find.CountArticles()
				return err
			})
			g.Go(func() error {
				var err error
				articles, err = find.ListArticles(limit, offset)
				return err
			})
		}

		err = g.Wait()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := page.IndexData{
			Search:       search,
			Articles:     articles,
			HasMorePages: p*s < count,
			NextPage:     p + 1,
		}
		tmpl.Execute(w, data)
	})
}
