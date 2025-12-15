package routes

import (
	"net/http"
	"net/url"
	"regexp"
)

func cookie(r *http.Request, name string) (value string) {
	if c, _ := r.Cookie(name); c != nil {
		value = c.Value
	}
	return
}

func param(r *http.Request, key string) string {
	value, _ := url.QueryUnescape(r.PathValue(key))
	return value
}

func redir(templateURL string) http.HandlerFunc {
	re := regexp.MustCompile("{[^{}]*}")
	return func(w http.ResponseWriter, r *http.Request) {
		url := re.ReplaceAllStringFunc(templateURL, func(s string) string {
			// slice to remove the surrounding {}
			return param(r, s[1:len(s)-1])
		})

		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}

		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
