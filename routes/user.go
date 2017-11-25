package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func user(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps[0].Value

	var html []byte

	switch err := db.QueryRow(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id
		    FROM solutions
		ORDER BY hole, user_id, LENGTH(code), submitted
		), scored_leaderboard AS (
		  SELECT hole,
		         ROUND(
		             (
		                 (SELECT COUNT(distinct user_id) FROM solutions WHERE hole = l.hole)
		                 -
		                 RANK() OVER (PARTITION BY hole ORDER BY strokes)
		                 +
		                 1
		             )
		             *
		             (
		                 100.0
		                 /
		                 (SELECT COUNT(distinct user_id) FROM solutions WHERE hole = l.hole)
		             )
		         ) score,
		         strokes,
		         submitted,
		         user_id
		    FROM leaderboard l
		), summed_leaderboard AS (
		  SELECT user_id,
		         SUM(score),
		         COUNT(*),
		         MAX(submitted)
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT CONCAT(
		             login,
		             '?s=100"><h1>',
		             login,
		             '</h1><table><tr><td>',
		             TO_CHAR(sum, 'FM9,999'),
		             '<td>points<tr><td>',
		             count,
		             '<td>holes</table><hr>'
		         )
		    FROM summed_leaderboard
		    JOIN users ON id = user_id
		   WHERE login = $1`,
		user,
	).Scan(&html); err {
	case sql.ErrNoRows:
		print404(w, r)
	case nil:
		printHeader(w, r, 200)
		w.Write([]byte("<link rel=stylesheet href=" + userCssPath + `><main><img src="//avatars.githubusercontent.com/`))
		w.Write(html)
	default:
		panic(err)
	}
}
