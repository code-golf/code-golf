package routes

import (
	"net/http"
	"sort"

	"github.com/JRaspass/code-golf/cookie"
	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Fail struct {
		Hole Hole
		Lang Lang
	}

	type Row struct {
		LangID, Login string
		Me            bool
		Rank, Strokes int
	}

	data := struct {
		Fails []Fail
		Holes []Hole
		Rows  map[string][]*Row
	}{
		Holes: make([]Hole, len(holes)),
		Rows:  make(map[string][]*Row),
	}

	// Sort by difficulty, then name. FIXME Seems to disagree wrt τ.
	copy(data.Holes, holes)
	sort.Slice(data.Holes, func(i, j int) bool {
		if data.Holes[i].Difficulty == data.Holes[j].Difficulty {
			return data.Holes[i].Name < data.Holes[j].Name
		}

		return data.Holes[i].Difficulty < data.Holes[j].Difficulty
	})

	userID, _ := cookie.Read(r)

	if userID != 0 {
		rows, err := db.Query(
			` SELECT hole, lang
			    FROM solutions
			   WHERE failing AND user_id = $1
			ORDER BY hole`,
			userID,
		)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		for rows.Next() {
			var holeID, langID string

			if err := rows.Scan(&holeID, &langID); err != nil {
				panic(err)
			}

			data.Fails = append(data.Fails, Fail{holeByID[holeID], langByID[langID]})
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         lang,
		         LENGTH(code) strokes,
		         submitted,
		         user_id
		    FROM solutions
		   WHERE NOT failing
		ORDER BY hole, user_id, LENGTH(code), submitted
		), ranked_leaderboard AS (
		  SELECT hole,
		         lang,
		         RANK()       OVER (PARTITION BY hole ORDER BY strokes),
		         ROW_NUMBER() OVER (PARTITION BY hole ORDER BY strokes, submitted),
		         strokes,
		         submitted,
		         user_id
		    FROM leaderboard
		) SELECT hole,
		         lang,
		         login,
		         user_id = $1,
		         rank,
		         row_number,
		         strokes
		    FROM ranked_leaderboard
		    JOIN users on user_id = id
		   WHERE row_number < 6 OR user_id = $1
		ORDER BY row_number`,
		userID,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var holeID string
		var row Row
		var rowNumber int

		if err := rows.Scan(
			&holeID,
			&row.LangID,
			&row.Login,
			&row.Me,
			&row.Rank,
			&rowNumber,
			&row.Strokes,
		); err != nil {
			panic(err)
		}

		if rowNumber < 6 {
			data.Rows[holeID] = append(data.Rows[holeID], &row)
		} else {
			data.Rows[holeID][3] = nil
			data.Rows[holeID][4] = &row
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	w.Header().Set(
		"Strict-Transport-Security",
		"max-age=31536000;includeSubDomains;preload",
	)
	Render(w, r, http.StatusOK, "index", data)
}
