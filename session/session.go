package session

import (
	"context"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/jmoiron/sqlx"
)

type key struct{}

type Session struct {
	Database   *sqlx.DB
	Golfer     *golfer.Golfer
	GolferInfo *golfer.GolferInfo
}

func Create(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key{}, &Session{}))
}

func Get(r *http.Request) *Session {
	return r.Context().Value(key{}).(*Session)
}

// Database gets the database handle from the request context.
func Database(r *http.Request) *sqlx.DB {
	return Get(r).Database
}

// Golfer gets the golfer from the request context.
func Golfer(r *http.Request) *golfer.Golfer {
	return Get(r).Golfer
}

// GolferInfo gets the GolferInfo object from the request context.
func GolferInfo(r *http.Request) *golfer.GolferInfo {
	return Get(r).GolferInfo
}
