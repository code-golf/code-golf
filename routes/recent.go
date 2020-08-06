package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// Recent serves GET /recent/{lang}
func Recent(w http.ResponseWriter, r *http.Request) {
	langID := param(r, "lang")

	if langID == "" {
		http.Redirect(w, r, "/recent/all-langs", http.StatusPermanentRedirect)
		return
	}

	if _, ok := lang.ByID[langID]; langID != "all-langs" && !ok {
		NotFound(w, r)
		return
	}

	rows, err := session.Database(r).Query(
		`WITH solution_lengths AS (
        SELECT hole,
               lang,
               login,
               bytes,
               chars,
               submitted
          FROM solutions
          JOIN users on user_id = id
         WHERE NOT failing
           AND $1 IN ('all-langs', lang::text)
     )  SELECT t1.hole,
               t1.lang,
               login,
               bytes,
               t1.chars,
               rank,
               COUNT(*) - 1 tie_count,
               t1.submitted
          FROM solution_lengths AS t1
    INNER JOIN (
        SELECT RANK() OVER (PARTITION BY hole, lang ORDER BY chars) rank,
               hole,
               lang,
               chars,
               submitted
          FROM solution_lengths
    ) AS t2
            ON t1.hole = t2.hole
           AND t1.lang = t2.lang
           AND t1.chars = t2.chars
           AND t2.submitted <= t1.submitted
      GROUP BY t1.hole,
               t1.lang,
               login,
               bytes,
               t1.chars,
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
		Hole                         hole.Hole
		Lang                         lang.Lang
		Login                        string
		Bytes, Chars, Rank, TieCount int
		Submitted                    time.Time
	}

	var recents []recent

	for rows.Next() {
		var holeID, langID string
		var r recent

		if err := rows.Scan(
			&holeID,
			&langID,
			&r.Login,
			&r.Bytes,
			&r.Chars,
			&r.Rank,
			&r.TieCount,
			&r.Submitted,
		); err != nil {
			panic(err)
		}

		r.Hole = hole.ByID[holeID]
		r.Lang = lang.ByID[langID]

		recents = append(recents, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	title := "Recent Solutions in "

	if langID == "all-langs" {
		title += "All Langs"
	} else {
		title += lang.ByID[langID].Name
	}

	data := struct {
		LangID  string
		Langs   []lang.Lang
		Recents []recent
	}{
		LangID:  langID,
		Langs:   lang.List,
		Recents: recents,
	}

	render(w, r, "recent", title, data)
}
