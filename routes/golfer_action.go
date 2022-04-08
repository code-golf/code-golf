package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

const followLimit = 10

// POST /golfers/{golfer}/{action}
func golferActionPOST(w http.ResponseWriter, r *http.Request) {
	action := param(r, "action")
	golfer := session.Golfer(r)
	target := session.GolferInfo(r)

	// Must be logged in and not acting on yourself.
	if golfer == nil || golfer.ID == target.ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx, err := session.Database(r).BeginTx(r.Context(), nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// Check we're not going to hit the follow limit.
	// TODO This probably ought to be folded into the SQL somehow.
	if action == "follow" {
		count := 0

		if err := tx.QueryRow(
			"SELECT COUNT(*) FROM follows WHERE follower_id = $1",
			golfer.ID,
		).Scan(&count); err != nil {
			panic(err)
		} else if count >= followLimit {
			w.WriteHeader(http.StatusBadRequest)
			render(w, r, "golfer/follow-limit", nil, target.Name)
			return
		}
	}

	var sql string
	switch action {
	case "follow":
		sql = "INSERT INTO follows VALUES ($1, $2) ON CONFLICT DO NOTHING"
	case "unfollow":
		sql = "DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2"
	}

	if _, err := tx.Exec(sql, golfer.ID, target.ID); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfers/"+target.Name, http.StatusSeeOther)
}
