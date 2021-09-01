package routes

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func cookie(r *http.Request, name string) (value string) {
	if c, _ := r.Cookie(name); c != nil {
		value = c.Value
	}
	return
}

func param(r *http.Request, key string) string {
	value, _ := url.QueryUnescape(chi.URLParam(r, key))
	return value
}
