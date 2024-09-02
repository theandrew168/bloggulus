package page

import (
	_ "embed"
	"net/http"
	"text/template"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/finder"
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

		// blogs, err := find.ListBlogsForAccount(nil)
		// if err != nil {
		// 	http.Error(w, err.Error(), 500)
		// 	return
		// }

		blogs := []finder.BlogForAccount{
			{ID: uuid.New(), Title: "Nice blog", IsFollowing: true},
			{ID: uuid.New(), Title: "Other blog", IsFollowing: true},
			{ID: uuid.New(), Title: "Bad blog", IsFollowing: false},
		}

		data := BlogsPageData{
			Blogs: blogs,
		}
		tmpl.Execute(w, data)
	})
}
