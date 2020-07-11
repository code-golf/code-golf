package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
	"github.com/gofrs/uuid"
)

var hmacKey []byte

func init() {
	var err error
	if hmacKey, err = base64.RawURLEncoding.DecodeString(os.Getenv("HMAC_KEY")); err != nil {
		panic(err)
	}
}

// GolferHandler adds the golfer to the context if logged in.
func GolferHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("__Host-user"); err == nil {
			if i := strings.LastIndexByte(cookie.Value, ':'); i != -1 {
				mac := hmac.New(sha256.New, hmacKey)
				mac.Write([]byte(cookie.Value[:i]))

				if subtle.ConstantTimeCompare(
					[]byte(cookie.Value[i+1:]),
					[]byte(base64.RawURLEncoding.EncodeToString(mac.Sum(nil))),
				) == 1 {
					var golfer golfer.Golfer

					j := strings.IndexByte(cookie.Value, ':')
					golfer.ID, _ = strconv.Atoi(cookie.Value[:j])
					golfer.Name = cookie.Value[j+1 : i]

					r = session.Set(r, "golfer", &golfer)

					// Port them to the new cookie.
					cookie := http.Cookie{
						HttpOnly: true,
						Name:     "__Host-session",
						Path:     "/",
						SameSite: http.SameSiteLaxMode,
						Secure:   true,
					}

					if err := session.Database(r).QueryRow(
						`WITH golfer AS (
						    INSERT INTO users (id, login) VALUES ($1, $2)
						    ON CONFLICT (id) DO UPDATE SET login = excluded.login
						      RETURNING id
						) INSERT INTO sessions (user_id) SELECT * FROM golfer RETURNING id`,
						golfer.ID, golfer.Name,
					).Scan(&cookie.Value); err != nil {
						panic(err)
					}

					http.SetCookie(w, &cookie)

					http.SetCookie(w, &http.Cookie{
						HttpOnly: true,
						MaxAge:   -1,
						Name:     "__Host-user",
						Path:     "/",
						SameSite: http.SameSiteLaxMode,
						Secure:   true,
					})
				}
			}
		}

		if cookie, _ := r.Cookie("__Host-session"); cookie != nil {
			var golfer golfer.Golfer

			if err := session.Database(r).QueryRow(
				`WITH golfer AS (
				    UPDATE sessions SET last_used = DEFAULT WHERE id = $1
				 RETURNING user_id
				) SELECT admin, id, login FROM users JOIN golfer ON id = user_id`,
				uuid.FromStringOrNil(cookie.Value),
			).Scan(&golfer.Admin, &golfer.ID, &golfer.Name); err == nil {
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
