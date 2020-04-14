package middleware

import "github.com/go-chi/chi/middleware"

// Re-export the upstream middleware we use via this namespace.
var (
	Compress        = middleware.Compress
	Logger          = middleware.Logger
	Recoverer       = middleware.Recoverer
	RedirectSlashes = middleware.RedirectSlashes
	WithValue       = middleware.WithValue
)
