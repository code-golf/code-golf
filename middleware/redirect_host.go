package middleware

import (
	"net/http"
	"syscall"
)

var host = "code.golf"

func init() {
	if _, dev := syscall.Getenv("DEV"); dev {
		host = "localhost"
	}
}

// Redirect www (or any incorrect domain) to apex.
func RedirectHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == host {
			next.ServeHTTP(w, r)
		} else {
			url := "https://" + host + r.RequestURI
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
	})
}
