package web

import (
	"io"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

// TODO: Rename p / s to page / size.
func HandleIndexPage(find *finder.Finder) http.Handler {
	tmpl := page.NewIndex()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		data := page.IndexData{
			BaseData: util.TemplateBaseData(r, w),

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
