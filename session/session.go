package session

import (
	"context"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/jmoiron/sqlx"
)

type key string

func Set(r *http.Request, k string, v any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key(k), v))
}

// Database gets the database handle from the request context.
func Database(r *http.Request) *sqlx.DB {
	return r.Context().Value(key("database")).(*sqlx.DB)
}

// Golfer gets the golfer from the request context.
func Golfer(r *http.Request) *golfer.Golfer {
	golfer, _ := r.Context().Value(key("golfer")).(*golfer.Golfer)
	return golfer
}

// GolferInfo gets the GolferInfo object from the request context.
func GolferInfo(r *http.Request) *golfer.GolferInfo {
	info, _ := r.Context().Value(key("golfer-info")).(*golfer.GolferInfo)
	return info
}
