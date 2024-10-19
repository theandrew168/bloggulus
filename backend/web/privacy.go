package web

import (
	"io"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandlePrivacyPolicy() http.Handler {
	tmpl := page.NewPrivacyPolicy()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := page.PrivacyPolicyData{
			BaseData: util.TemplateBaseData(r, w),
		}
		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}
