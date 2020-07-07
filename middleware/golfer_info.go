package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/routes"
	"github.com/code-golf/code-golf/session"
	"github.com/go-chi/chi"
)

// GolferInfoHandler adds the GolferInfo object handle to the context.
func GolferInfoHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		info := golfer.GetInfo(session.Database(r), name)

		switch {
		case info == nil:
			routes.NotFound(w, r)
		case info.Name != name:
			// TODO Handle /holes suffix.
			http.Redirect(w, r, "/golfers/"+info.Name, http.StatusPermanentRedirect)
		default:
			next.ServeHTTP(w, session.Set(r, "golfer-info", info))
		}
	})
}
