package routes

import (
	"net/http"
	"strings"
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
               scoring,
               bytes,
               chars,
               case when scoring = 'chars' then chars else bytes end strokes,
               submitted
          FROM solutions
          JOIN code  ON code_id = code.id
          JOIN users ON user_id = users.id
         WHERE NOT failing
           AND $1 IN ('all-langs', lang::text)
     )  SELECT t1.hole,
               t1.lang,
               login,
               t1.scoring,
               bytes,
               chars,
               t1.strokes,
               rank,
               COUNT(*) - 1 tie_count,
               t1.submitted
          FROM solution_lengths AS t1
    INNER JOIN (
        SELECT RANK() OVER (PARTITION BY hole, lang, scoring ORDER BY strokes) rank,
               hole,
               lang,
               scoring,
               strokes,
               submitted
          FROM solution_lengths
    ) AS t2
            ON t1.hole = t2.hole
           AND t1.lang = t2.lang
           AND t1.scoring = t2.scoring
           AND t1.strokes = t2.strokes
           AND t2.submitted <= t1.submitted
      GROUP BY t1.hole,
               t1.lang,
               login,
               t1.scoring,
               bytes,
               chars,
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
		Hole                                  hole.Hole
		Lang                                  lang.Lang
		Login, Scoring                        string
		Bytes, Chars, Strokes, Rank, TieCount int
		Submitted                             time.Time
	}

	var recents []recent
	var previous recent

	for rows.Next() {
		var holeID, langID string
		var r recent

		if err := rows.Scan(
			&holeID,
			&langID,
			&r.Login,
			&r.Scoring,
			&r.Bytes,
			&r.Chars,
			&r.Strokes,
			&r.Rank,
			&r.TieCount,
			&r.Submitted,
		); err != nil {
			panic(err)
		}

		r.Hole = hole.ByID[holeID]
		r.Lang = lang.ByID[langID]
		r.Scoring = strings.Title(r.Scoring)

		// If all of the information in two rows matches, other than the scoring, collapse them into one.
		if previous.Login == r.Login &&
			previous.Hole.ID == r.Hole.ID &&
			previous.Lang.ID == r.Lang.ID &&
			previous.Strokes == r.Strokes &&
			(previous.Rank == r.Rank || previous.Rank > 3 && r.Rank > 3) &&
			(r.Rank > 3 || (previous.TieCount > 0) == (r.TieCount > 0)) {
			recents[len(recents)-1].Scoring = "Both"
		} else {
			recents = append(recents, r)
			previous = r
		}
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
