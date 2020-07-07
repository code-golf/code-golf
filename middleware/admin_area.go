package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/routes"
	"github.com/code-golf/code-golf/session"
)

// AdminArea enforces that an admin golfer is logged in.
func AdminArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if golfer := session.Golfer(r); golfer != nil && golfer.Admin {
			next.ServeHTTP(w, r)
		} else {
			routes.Forbidden(w, r)
		}
	})
}
