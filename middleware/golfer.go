package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
	"github.com/gofrs/uuid"
)

// GolferHandler adds the golfer to the context if logged in.
func GolferHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, _ := r.Cookie("__Host-session"); cookie != nil {
			var golfer golfer.Golfer

			if err := session.Database(r).QueryRow(
				`WITH golfer AS (
				    UPDATE sessions SET last_used = DEFAULT WHERE id = $1
				 RETURNING user_id
				) SELECT admin, delete, id, keymap, login, time_zone
				    FROM users
				    JOIN golfer ON id = user_id`,
				uuid.FromStringOrNil(cookie.Value),
			).Scan(
				&golfer.Admin,
				&golfer.Delete,
				&golfer.ID,
				&golfer.Keymap,
				&golfer.Name,
				&golfer.TimeZone,
			); err == nil {
				r = session.Set(r, "golfer", &golfer)

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
			} else if !errors.Is(err, sql.ErrNoRows) {
				panic(err)
			}
		}

		next.ServeHTTP(w, r)
	})
}
