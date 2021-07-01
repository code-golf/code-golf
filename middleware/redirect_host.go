package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
)

// Redirect www (or any incorrect domain) to apex.
func RedirectHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == config.Host {
			next.ServeHTTP(w, r)
		} else {
			url := "https://" + config.Host + r.RequestURI
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
	})
}
