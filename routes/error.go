package routes

import (
	"fmt"
	"net/http"

	"github.com/code-golf/code-golf/middleware"
)

// errorMiddleware writes HTML/JSON bodies for 4xx/5xx status codes.
// Can't be in middleware as it uses routes.render and that would be a cycle.
func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		// Write an error body for 4xx & 5xx if we have yet to write a body.
		if code := ww.Status(); code >= 400 && ww.BytesWritten() == 0 {
			if ww.Header().Get("Content-Type") == "application/json" {
				ww.Write([]byte("null\n"))
			} else {
				render(w, r, "error",
					fmt.Sprintf("%d %s", code, http.StatusText(code)))
			}
		}
	})
}
