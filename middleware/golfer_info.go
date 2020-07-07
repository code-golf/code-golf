package middleware

import (
	"context"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/go-chi/chi"
)

var golferInfoKey = key("golferInfo")

// GolferInfoHandler adds the GolferInfo object handle to the context.
func GolferInfoHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		info := golfer.GetInfo(Database(r), name)

		switch {
		case info == nil:
			// TODO Serve custom 404 without import cycle.
			w.WriteHeader(http.StatusNotFound)
		case info.Name != name:
			// TODO Handle /holes suffix.
			http.Redirect(w, r, "/golfers/"+info.Name, http.StatusPermanentRedirect)
		default:
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), golferInfoKey, info)))
		}
	})
}

// GolferInfo gets the GolferInfo object from the request context.
func GolferInfo(r *http.Request) *golfer.GolferInfo {
	info, _ := r.Context().Value(golferInfoKey).(*golfer.GolferInfo)
	return info
}
