package util

import "net/http"

func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") != ""
}
