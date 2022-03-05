package session

import (
	"context"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/jackc/pgx/v4/pgxpool"
)

type key string

func Set(r *http.Request, k string, v interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key(k), v))
}

// Database gets the database handle from the request context.
func Database(r *http.Request) *pgxpool.Pool {
	return r.Context().Value(key("database")).(*pgxpool.Pool)
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
