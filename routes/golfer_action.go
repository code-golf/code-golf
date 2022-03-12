package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// POST /golfers/{golfer}/{action}
func golferActionPOST(w http.ResponseWriter, r *http.Request) {
	golfer := session.Golfer(r)
	target := session.GolferInfo(r)

	// Must be logged in and not acting on yourself.
	if golfer == nil || golfer.ID == target.ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var sql string
	switch param(r, "action") {
	case "follow":
		sql = "INSERT INTO follows VALUES ($1, $2) ON CONFLICT DO NOTHING"
	case "unfollow":
		sql = "DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2"
	}

	if _, err := session.Database(r).Exec(sql, golfer.ID, target.ID); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/"+target.Name, http.StatusSeeOther)
}
