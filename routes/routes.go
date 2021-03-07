package routes

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func param(r *http.Request, key string) string {
	value, _ := url.QueryUnescape(chi.URLParam(r, key))
	return value
}
