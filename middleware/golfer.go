package middleware

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
	"github.com/gofrs/uuid"
)

// Golfer adds the golfer to the context if logged in.
func Golfer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, _ := r.Cookie("__Host-session"); cookie != nil {
			sessionID := uuid.FromStringOrNil(cookie.Value)
			if golfer := golfer.Get(session.Database(r), sessionID); golfer != nil {
				session.Get(r).Golfer = golfer

				// Refresh the cookie.
				http.SetCookie(w, &http.Cookie{
					HttpOnly: true,
					MaxAge:   int(30 * 24 * time.Hour / time.Second),
					Name:     "__Host-session",
					Path:     "/",
					SameSite: http.SameSiteLaxMode,
					Secure:   true,
					Value:    cookie.Value,
				})
			}
		}

		next.ServeHTTP(w, r)
	})
}
