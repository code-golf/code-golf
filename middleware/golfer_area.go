package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
)

// GolferArea enforces that a golfer is logged in.
func GolferArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g := session.Golfer(r); g != nil {
			info := golfer.GetInfo(session.Database(r), g.Name)
			next.ServeHTTP(w, session.Set(r, "golfer-info", info))
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	})
}
