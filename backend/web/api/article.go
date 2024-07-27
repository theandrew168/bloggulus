package api

import (
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

type jsonArticle struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	BlogTitle   string    `json:"blogTitle"`
	BlogURL     string    `json:"blogURL"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
}

func marshalArticle(article *model.Article) jsonArticle {
	a := jsonArticle{
		Title:       article.Title(),
		URL:         article.URL(),
		BlogTitle:   article.BlogTitle(),
		BlogURL:     article.BlogURL(),
		PublishedAt: article.PublishedAt(),
		Tags:        article.Tags(),
	}
	return a
}

func HandleArticleList(store *storage.Storage) http.Handler {
	type response struct {
		Count    int           `json:"count"`
		Articles []jsonArticle `json:"articles"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := util.NewErrors()
		qs := r.URL.Query()

		// check search param
		q := qs.Get("q")

		// check pagination params
		page := util.ReadInt(qs, "page", 1, e)
		e.CheckField(page >= 1, "Page must be greater than or equal to 1", "page")

		size := util.ReadInt(qs, "size", 20, e)
		e.CheckField(size >= 1, "Size must be greater than or equal to 1", "size")
		e.CheckField(size <= 50, "Size must be less than or equal to 50", "size")

		if !e.Valid() {
			util.FailedValidationResponse(w, r, e)
			return
		}

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

		err := g.Wait()
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Count: count,
			// use make here to encode JSON as "[]" instead of "null" if empty
			Articles: make([]jsonArticle, 0),
		}

		for _, article := range articles {
			resp.Articles = append(resp.Articles, marshalArticle(article))
		}

		err = util.WriteJSON(w, 200, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	})
}
