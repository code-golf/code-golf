package middleware

import "github.com/go-chi/chi/v5/middleware"

// Re-export the upstream middleware we use via this namespace.
var (
	Compress        = middleware.Compress
	Recoverer       = middleware.Recoverer
	RedirectSlashes = middleware.RedirectSlashes
	WithValue       = middleware.WithValue
)
