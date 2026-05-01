package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/session"
	"github.com/mileusna/useragent"
)

type ua useragent.UserAgent

func (u *ua) Scan(src any) error {
	*u = ua(useragent.Parse(src.(string)))
	return nil
}

// GET /golfer/settings/sessions
func golferSettingsSessionsGET(w http.ResponseWriter, r *http.Request) {
	var sessions []struct {
		Browser  ua
		ID       string
		LastUsed time.Time
	}

	if err := session.Database(r).Select(
		&sessions,
		` SELECT browser, id, last_used
		    FROM sessions
		   WHERE user_id = $1
		ORDER BY last_used DESC`,
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	render(w, r, "golfer/settings", sessions, "Settings: Sessions")
}

// POST /golfer/settings/sessions
func golferSettingsSessionsPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	session.Database(r).MustExec(
		"DELETE FROM sessions WHERE id = $1 AND user_id = $2",
		r.FormValue("id"),
		session.Golfer(r).ID,
	)

	http.Redirect(w, r, "/golfer/settings/sessions", http.StatusSeeOther)
}
