package api

import (
	"net/http"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
	"github.com/theandrew168/bloggulus/backend/web/validator"
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
		v := validator.New()
		qs := r.URL.Query()

		// check search param
		q := qs.Get("q")

		// check pagination params
		page := util.ReadInt(qs, "page", 1, v)
		v.Check(page >= 1, "page", "must be greater than or equal to 1")

		size := util.ReadInt(qs, "size", 20, v)
		v.Check(size >= 1, "size", "must be greater than or equal to 1")
		v.Check(size <= 50, "size", "must be less than or equal to 50")

		if !v.Valid() {
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		limit, offset := util.PageSizeToLimitOffset(page, size)

		var count int
		var articles []*model.Article
		var err error

		if q != "" {
			count, err = store.Article().CountSearch(q)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}

			articles, err = store.Article().ListSearch(q, limit, offset)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}
		} else {
			count, err = store.Article().Count()
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}

			articles, err = store.Article().List(limit, offset)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}
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
