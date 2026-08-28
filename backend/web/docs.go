package web

import (
	"io"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandlePrivacy() http.Handler {
	tmpl := page.NewPrivacy()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := page.PrivacyData{
			BaseData: util.GetTemplateBaseData(r, w),
		}
		util.Render(w, r, http.StatusOK, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}
