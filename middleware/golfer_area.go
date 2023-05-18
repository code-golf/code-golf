package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
)

// GolferArea enforces that a golfer is logged in.
func GolferArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s := session.Get(r); s.Golfer != nil {
			s.GolferInfo = golfer.GetInfo(s.Database, s.Golfer.Name)
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	})
}
