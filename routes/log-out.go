package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// POST /log-out
func logOutPOST(w http.ResponseWriter, r *http.Request) {
	if cookie, _ := r.Cookie("__Host-session"); cookie != nil {
		session.Database(r).MustExec(
			"DELETE FROM sessions WHERE id = uuid_or_null($1)", cookie.Value)
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
