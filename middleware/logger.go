package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/code-golf/code-golf/session"
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
	// Disable the access logs under e2e, they're just too noisy.
	if _, e2e := os.LookupEnv("E2E"); e2e {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			golfer := "-"
			if g := session.Golfer(r); g != nil {
				golfer = g.Name
			}

			log.Printf(
				"%s%d "+magenta+"%4s"+reset+" %s "+magenta+"%v"+white+" %s %s"+reset,
				statusColours[ww.Status()/100],
				ww.Status(),
				r.Method,
				r.URL.Path,
				time.Since(start).Round(time.Millisecond),
				golfer,
				r.UserAgent(),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
