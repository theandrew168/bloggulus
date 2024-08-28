package page

import (
	_ "embed"
	"net/http"
	"strconv"
	"text/template"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
	"golang.org/x/sync/errgroup"
)

//go:embed index.html
var html string

type IndexData struct {
	Search       string
	Articles     []*model.Article
	HasMorePages bool
	NextPage     int
}

func HandleIndex(store *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse(html)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// check search param
		q := r.URL.Query().Get("q")

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
		var articles []*model.Article

		var g errgroup.Group
		if q != "" {
			g.Go(func() error {
				var err error
				count, err = store.Article().CountSearch(q)
				return err
			})
			g.Go(func() error {
				var err error
				articles, err = store.Article().ListSearch(q, limit, offset)
				return err
			})
		} else {
			g.Go(func() error {
				var err error
				count, err = store.Article().Count()
				return err
			})
			g.Go(func() error {
				var err error
				articles, err = store.Article().List(limit, offset)
				return err
			})
		}

		err = g.Wait()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := IndexData{
			Search:       q,
			Articles:     articles,
			HasMorePages: page*size < count,
			NextPage:     page + 1,
		}
		tmpl.Execute(w, data)
	})
}
