package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

// Golfer adds the golfer to the context if logged in.
func Golfer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, _ := r.Cookie("__Host-session"); cookie != nil {
			var golfer golfer.Golfer
			var timeZone sql.NullString

			if err := session.Database(r).QueryRow(
				`WITH golfer AS (
				    UPDATE sessions SET last_used = DEFAULT WHERE id = $1
				 RETURNING user_id
				), failing AS (
				    SELECT hole, lang
				      FROM solutions
				      JOIN golfer USING(user_id)
				     WHERE failing
				  GROUP BY hole, lang
				  ORDER BY hole, lang
				)  SELECT u.admin,
				          COALESCE(bytes.points, 0),
				          COALESCE(chars.points, 0),
				          COALESCE(u.country, ''),
				          u.delete,
				          (SELECT COALESCE(json_agg(failing), '[]') FROM failing),
				          u.id,
				          u.layout,
				          u.keymap,
				          u.login,
				          COALESCE(r.login, ''),
				          u.show_country,
				          u.sponsor,
				          u.theme,
				          u.time_zone,
				          ARRAY(
				              SELECT trophy
				                FROM trophies
				               WHERE user_id = u.id
				            ORDER BY trophy
				          ),
				          ARRAY(
				              SELECT followee_id
				                FROM follows
				               WHERE follower_id = u.id
				            ORDER BY followee_id
				          ),
				          ARRAY(
				              SELECT DISTINCT hole
				                FROM solutions
				               WHERE user_id = u.id
				            ORDER BY hole
				          )
				     FROM users  u
				     JOIN golfer g ON u.id = g.user_id
				LEFT JOIN users  r ON r.id = u.referrer_id
				LEFT JOIN points bytes ON u.id = bytes.user_id AND bytes.scoring = 'bytes'
				LEFT JOIN points chars ON u.id = chars.user_id AND chars.scoring = 'chars'`,
				uuid.FromStringOrNil(cookie.Value),
			).Scan(
				&golfer.Admin,
				&golfer.BytesPoints,
				&golfer.CharsPoints,
				&golfer.Country,
				&golfer.Delete,
				&golfer.FailingSolutions,
				&golfer.ID,
				&golfer.Layout,
				&golfer.Keymap,
				&golfer.Name,
				&golfer.Referrer,
				&golfer.ShowCountry,
				&golfer.Sponsor,
				&golfer.Theme,
				&timeZone,
				pq.Array(&golfer.Cheevos),
				pq.Array(&golfer.Following),
				pq.Array(&golfer.Holes),
			); err == nil {
				golfer.TimeZone, _ = time.LoadLocation(timeZone.String)

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
