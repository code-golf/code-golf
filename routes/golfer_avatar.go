package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{name}/avatar
func golferAvatarGET(w http.ResponseWriter, r *http.Request) {
	// TODO This will need to be more complex when multi-OAuth is supported
	//      and GitHub can't be presumed.
	var id string
	if err := session.Database(r).Get(
		&id,
		`SELECT c.id
		   FROM users       u
		   JOIN connections c ON c.user_id = u.id AND c.connection = 'github'
		  WHERE login = $1`,
		param(r, "name"),
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	w.Header().Set("Cache-Control", "max-age=300")

	url := "//avatars.githubusercontent.com/u/" + id
	if size := param(r, "size"); size != "" {
		url += "?s=" + size
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
