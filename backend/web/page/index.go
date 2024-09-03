package page

import (
	_ "embed"
	"net/http"
	"strconv"
	"text/template"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/util"
	"golang.org/x/sync/errgroup"
)

//go:embed index.html
var indexHTML string

type IndexPageData struct {
	Search       string
	Articles     []finder.Article
	HasMorePages bool
	NextPage     int
}

func HandleIndexPage(find *finder.Finder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("page").Parse(indexHTML)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// check search param
		search := r.URL.Query().Get("q")

		// check page param
		page, err := strconv.Atoi(r.URL.Query().Get("p"))
		if err != nil {
			page = 1
		}

		if page < 1 {
			page = 1
		}

		size := 20
		limit, offset := util.PageSizeToLimitOffset(page, size)

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

		data := IndexPageData{
			Search:       search,
			Articles:     articles,
			HasMorePages: page*size < count,
			NextPage:     page + 1,
		}
		tmpl.Execute(w, data)
	})
}
