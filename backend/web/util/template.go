package util

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

func TemplateBaseData(r *http.Request, w http.ResponseWriter) layout.BaseData {
	data := layout.BaseData{}

	toastCookie, err := r.Cookie(ToastCookieName)
	if err == nil {
		data.Toast = toastCookie.Value

		cookie := NewExpiredCookie(ToastCookieName)
		http.SetCookie(w, &cookie)
	}

	return data
}
