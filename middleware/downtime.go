package middleware

import "net/http"

const html = `<!doctype html>

<title>Code Golf</title>

<h1>503 Service Unavailable</h1>

<p>Code Golf is down for maintenance, please try again in a few minutes.`

// Downtime serves a 503 Service Unavailable.
func Downtime(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(html))
	})
}
