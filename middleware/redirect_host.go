package middleware

import (
	"net/http"
	"syscall"

	"github.com/code-golf/code-golf/session"
)

var (
	host     = "code.golf"
	betaHost = "beta.code.golf"
)

func init() {
	if _, e2e := syscall.Getenv("E2E"); e2e {
		host = "app:1443"
	} else if _, dev := syscall.Getenv("DEV"); dev {
		host = "localhost"
	}
}

// Redirect www (or any incorrect domain) to apex.
func RedirectHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Host {
		case betaHost:
			r = session.Set(r, "beta", true)
			fallthrough
		case host:
			next.ServeHTTP(w, r)
		default:
			url := "https://" + host + r.RequestURI
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
	})
}
