package util

import (
	"net/http"

	"github.com/justinas/nosurf"

	"github.com/theandrew168/bloggulus/backend/web/layout"
)

func GetTemplateBaseData(r *http.Request, w http.ResponseWriter) layout.BaseData {
	data := layout.BaseData{}

	account, isLoggedIn := GetContextAccount(r)
	if isLoggedIn {
		data.Account = account
	}

	csrfToken := nosurf.Token(r)
	if csrfToken != "" {
		data.CSRFToken = csrfToken
	}

	toastCookie, err := r.Cookie(ToastCookieName)
	if err == nil {
		data.Toast = toastCookie.Value

		cookie := NewExpiredCookie(ToastCookieName)
		http.SetCookie(w, &cookie)
	}

	conf, ok := GetContextConfig(r)
	if ok && conf.PlausibleDataDomain != "" {
		data.PlausibleDataDomain = conf.PlausibleDataDomain
	}
	if ok && conf.GoatCounterCode != "" {
		data.GoatCounterCode = conf.GoatCounterCode
	}

	return data
}
