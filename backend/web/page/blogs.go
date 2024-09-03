package page

import (
	_ "embed"
	"net/http"
	"text/template"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

//go:embed blogs.html
var blogsHTML string

type BlogsPageData struct {
	Search string
	Blogs  []finder.BlogForAccount
}

func HandleBlogsPage(find *finder.Finder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("blogs").Parse(blogsHTML)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		account, ok := util.ContextGetAccount(r)
		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		blogs, err := find.ListBlogsForAccount(account)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := BlogsPageData{
			Blogs: blogs,
		}
		tmpl.Execute(w, data)
	})
}
