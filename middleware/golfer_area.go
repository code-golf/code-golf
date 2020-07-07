package middleware

import (
	"context"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
)

// GolferArea enforces that a golfer is logged in.
func GolferArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g := Golfer(r); g != nil {
			info := golfer.GetInfo(Database(r), g.Name)
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), golferInfoKey, info)))
		} else {
			// TODO Serve custom 403 without import cycle.
			w.WriteHeader(http.StatusForbidden)
		}
	})
}
