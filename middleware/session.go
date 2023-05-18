package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// Session creates the session and stores it in the request context.
func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, session.Create(r))
	})
}
