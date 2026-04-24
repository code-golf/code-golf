package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/uuid"
)

// POST /log-out
func logOutPOST(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("__Host-session"); err == nil {
		if sessionID, err := uuid.Parse(cookie.Value); err == nil {
			session.Database(r).MustExec(
				"DELETE FROM sessions WHERE id = $1", sessionID)
		}
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		Name:     "__Host-session",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
