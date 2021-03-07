package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	red     = "\033[31;1m"
	green   = "\033[32;1m"
	yellow  = "\033[33;1m"
	blue    = "\033[34;1m"
	magenta = "\033[35;1m"
	cyan    = "\033[36;1m"
	white   = "\033[37;1m"
	reset   = "\033[0m"
)

var statusColours = [...]string{
	"",
	blue,   // 1xx Informational response
	green,  // 2xx Success
	cyan,   // 3xx Redirection
	yellow, // 4xx Client errors
	red,    // 5xx Server errors
}

// Logger logs each request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			log.Printf(
				"%s%d "+magenta+"%4s"+reset+" %s "+magenta+"%v"+white+" %s"+reset,
				statusColours[ww.Status()/100],
				ww.Status(),
				r.Method,
				r.URL.Path,
				time.Since(start).Round(time.Millisecond),
				r.UserAgent(),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
