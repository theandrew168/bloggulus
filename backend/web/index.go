package web

import (
	"io"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/ui"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

// TODO: Rename p / s to page / size.
func HandleIndexPage(find *finder.Finder) http.Handler {
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
		var articles []finder.Article

		// Two levels of decision making here:
		// 1. Is the user logged in?
		// 2. Is the user searching?
		var g errgroup.Group
		if isLoggedIn {
			if search != "" {
				g.Go(func() error {
					var err error
					count, err = find.CountSearchArticlesByAccount(account, search)
					return err
				})
				g.Go(func() error {
					var err error
					articles, err = find.SearchArticlesByAccount(account, search, limit, offset)
					return err
				})
			} else {
				g.Go(func() error {
					var err error
					count, err = find.CountArticlesByAccount(account)
					return err
				})
				g.Go(func() error {
					var err error
					articles, err = find.ListArticlesByAccount(account, limit, offset)
					return err
				})
			}
		} else {
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
		}

		err = g.Wait()
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		page := ui.IndexPage(ui.IndexPageData{
			PageLayoutData: util.GetPageLayoutData(r, w),

			Search:       search,
			Articles:     articles,
			HasMorePages: p*s < count,
			NextPage:     p + 1,
		})

		util.Render(w, r, 200, func(w io.Writer) error {
			return page.Render(w)
		})
	})
}
