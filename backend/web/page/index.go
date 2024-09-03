package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed index.html
var IndexHTML string

type IndexData struct {
	layout.BaseData

	Search       string
	Articles     []finder.Article
	HasMorePages bool
	NextPage     int
}

func Index() (*template.Template, error) {
	// Create the template and list necessary sources.
	tmpl := template.New("page")
	sources := []string{
		layout.BaseHTML,
		IndexHTML,
	}

	// Parse each source required to render this page.
	for _, source := range sources {
		_, err := tmpl.Parse(source)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

// TODO: Cache the parsed template? Lazy parse it?
// I almost think each page should be an object with methods?
// That way you can namespace the top-level and sub-level renders.
// Things like: Render(), RenderBlogs(), or RenderBlogRow().
func RenderIndex(wr io.Writer, data IndexData) error {
	tmpl, err := Index()
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, data)
}
