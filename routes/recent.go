package routes

import (
	"net/http"
	"time"
)

// Recent serves GET /recent/{lang}
func Recent(w http.ResponseWriter, r *http.Request) {
	langID := param(r, "lang")

	if langID == "" {
		langID = "all-langs"
	}

	if _, ok := langByID[langID]; langID != "all-langs" && !ok {
		render(w, r, http.StatusNotFound, "404", "", nil)
		return
	}

	rows, err := db(r).Query(
		`WITH solution_lengths AS (
        SELECT hole,
               lang,
               login,
               LENGTH(code) strokes,
               submitted
          FROM solutions
          JOIN users on user_id = id
         WHERE NOT failing
           AND $1 IN ('all-langs', lang::text)
     )  SELECT t1.hole,
               t1.lang,
               login,
               t1.strokes,
               rank,
               COUNT(*) - 1 tie_count,
               t1.submitted
          FROM solution_lengths AS t1
    INNER JOIN (
        SELECT RANK() OVER (PARTITION BY hole, lang ORDER BY strokes) rank,
               hole,
               lang,
               strokes,
               submitted
          FROM solution_lengths
    ) AS t2
            ON t1.hole = t2.hole
           AND t1.lang = t2.lang
           AND t1.strokes = t2.strokes
           AND t2.submitted <= t1.submitted
      GROUP BY t1.hole,
               t1.lang,
               login,
               t1.strokes,
               t1.submitted,
               rank
      ORDER BY t1.submitted DESC LIMIT 100`,
		langID,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	type recent struct {
		Hole      Hole
		Lang      Lang
		Login     string
		Strokes   int
		Rank      int
		TieCount  int
		Submitted time.Time
	}

	var recents []recent

	for rows.Next() {
		var holeID, langID string
		var r recent

		if err := rows.Scan(
			&holeID,
			&langID,
			&r.Login,
			&r.Strokes,
			&r.Rank,
			&r.TieCount,
			&r.Submitted,
		); err != nil {
			panic(err)
		}

		r.Hole = holeByID[holeID]
		r.Lang = langByID[langID]

		recents = append(recents, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	var title = "Recent Solutions in "

	if langID == "all-langs" {
		title += "All Langs"
	} else {
		title += langByID[langID].Name
	}

	data := struct {
		LangID  string
		Langs   []Lang
		Recents []recent
	}{
		LangID:  langID,
		Langs:   langs,
		Recents: recents,
	}

	render(w, r, http.StatusOK, "recent", title, data)
}
