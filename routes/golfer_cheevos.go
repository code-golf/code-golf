package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /golfers/{golfer}/cheevos
func golferCheevosGET(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer

	type Progress struct {
		Count, Percent, Progress int
		Earned                   *time.Time
	}

	data := map[string]Progress{}

	db := session.Database(r)
	rows, err := db.Query(
		`WITH count AS (
		    SELECT trophy, COUNT(*) FROM trophies GROUP BY trophy
		), earned AS (
		    SELECT trophy, earned FROM trophies WHERE user_id = $1
		) SELECT *, count * 100 / (SELECT COUNT(*) FROM users)
		    FROM count LEFT JOIN earned USING(trophy)`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cheevoID string
		var progress Progress

		if err := rows.Scan(
			&cheevoID,
			&progress.Count,
			&progress.Earned,
			&progress.Percent,
		); err != nil {
			panic(err)
		}

		data[cheevoID] = progress
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	// Caclulate progress
	// TODO Bake it into the cheevos table rather than calculating on the fly.
	cheevoProgress := func(sql string, cheevoIDs []string, args ...any) {
		count := -1

		for _, cheevoID := range cheevoIDs {
			progress := data[cheevoID]

			// No need to calculate progress if they've already earnt it.
			if progress.Earned != nil {
				continue
			}

			if count == -1 {
				if err := db.Get(&count, sql, append(args, golfer.ID)...); err != nil {
					panic(err)
				}
			}

			progress.Progress = count
			data[cheevoID] = progress
		}
	}

	cheevoIDs := func(cheevos []*config.Cheevo) []string {
		ids := make([]string, len(cheevos))
		for i, cheevo := range cheevos {
			ids[i] = cheevo.ID
		}
		return ids
	}

	for _, cheevo := range []struct {
		cheevo       string
		holes, langs []any
	}{
		{
			cheevo: "archivist",
			holes:  []any{"isbn"},
			langs:  []any{"basic", "cobol", "fortran", "lisp"},
		},
		{
			cheevo: "bird-is-the-word",
			holes:  []any{"levenshtein-distance"},
			langs:  []any{"awk", "prolog", "sql", "swift", "tcl", "wren"},
		},
		{
			cheevo: "happy-go-lucky",
			holes:  []any{"happy-numbers", "lucky-numbers"},
			langs:  []any{"go"},
		},
		{
			cheevo: "jeweler",
			holes:  []any{"diamonds"},
			langs:  []any{"crystal", "ruby"},
		},
		{
			cheevo: "mary-had-a-little-lambda",
			holes:  []any{"λ"},
			langs:  []any{"clojure", "coconut", "haskell", "lisp"},
		},
		{
			cheevo: "sounds-quite-nice",
			holes:  []any{"musical-chords"},
			langs:  []any{"c", "c-sharp", "d", "f-sharp"},
		},
	} {
		cheevoProgress(
			`SELECT COUNT(DISTINCT(hole, lang))
			   FROM solutions
			  WHERE NOT failing
			    AND hole    = ANY($1)
			    AND lang    = ANY($2)
			    AND user_id = $3`,
			[]string{cheevo.cheevo},
			pq.Array(cheevo.holes),
			pq.Array(cheevo.langs),
		)
	}

	cheevoProgress(
		`SELECT pangramglot(array_agg(DISTINCT lang))
		   FROM solutions
		  WHERE NOT failing AND hole = 'pangram-grep' AND user_id = $1`,
		[]string{"pangramglot"},
	)

	cheevoProgress(
		`SELECT COALESCE(EXTRACT(days FROM TIMEZONE('UTC', NOW()) - MIN(submitted)), 0)
		   FROM solutions
		  WHERE NOT FAILING AND user_id = $1`,
		[]string{"aged-like-fine-wine"},
	)

	cheevoProgress(
		`SELECT COUNT(DISTINCT hole)
		   FROM solutions
		  WHERE NOT failing AND user_id = $1`,
		cheevoIDs(config.CheevoTree["Total Holes"]),
	)

	cheevoProgress(
		`WITH langs AS (
		    SELECT COUNT(DISTINCT lang)
		      FROM solutions
		     WHERE NOT failing AND user_id = $1
		  GROUP BY hole
		) SELECT COALESCE(MAX(count), 0) FROM langs`,
		[]string{"polyglot", "polyglutton", "omniglot"},
	)

	cheevoProgress(
		"SELECT COALESCE(MAX(points), 0) FROM points WHERE user_id = $1",
		cheevoIDs(config.CheevoTree["Total Points"]),
	)

	render(w, r, "golfer/cheevos", data, golfer.Name)
}
