package routes

import (
	"net/http"
	"slices"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /golfers/{golfer}/cheevos
func golferCheevosGET(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer

	type Step struct {
		Complete   bool
		Name, Path string
	}

	type Progress struct {
		Count, Percent, Progress int
		Earned                   *time.Time
		Steps                    []Step
	}

	data := map[string]Progress{}

	db := session.Database(r)
	rows, err := db.Query(
		`WITH count AS (
		    SELECT cheevo, COUNT(*) FROM cheevos GROUP BY cheevo
		), earned AS (
		    SELECT cheevo, earned FROM cheevos WHERE user_id = $1
		) SELECT *, count * 100 / (SELECT COUNT(*) FROM users)
		    FROM count LEFT JOIN earned USING(cheevo)`,
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

	// Calculate progress
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
			cheevo: "alphabet-soup",
			holes:  []any{"scrambled-sort"},
			langs:  []any{"c", "d", "j", "k", "r", "v"},
		},
		{
			cheevo: "archivist",
			holes:  []any{"isbn"},
			langs:  []any{"basic", "cobol", "common-lisp", "forth", "fortran"},
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
			langs:  []any{"clojure", "coconut", "common-lisp", "haskell", "scheme"},
		},
		{
			cheevo: "sounds-quite-nice",
			holes:  []any{"musical-chords"},
			langs:  []any{"c", "c-sharp", "d", "f-sharp"},
		},
	} {
		progress := data[cheevo.cheevo]

		// No need to calculate progress if they've already earnt it.
		if progress.Earned != nil {
			continue
		}

		var completedSteps []struct{ Hole, Lang string }
		if err := db.Select(
			&completedSteps,
			`SELECT DISTINCT hole, lang
			   FROM stable_passing_solutions
			  WHERE hole    = ANY($1)
			    AND lang    = ANY($2)
			    AND user_id = $3`,
			pq.Array(cheevo.holes),
			pq.Array(cheevo.langs),
			golfer.ID,
		); err != nil {
			panic(err)
		}

		for _, holeID := range cheevo.holes {
			for _, langID := range cheevo.langs {
				hole := config.HoleByID[holeID.(string)]
				lang := config.LangByID[langID.(string)]

				step := Step{Path: "/" + hole.ID + "#" + lang.ID}

				step.Name = hole.Name
				if len(cheevo.langs) > len(cheevo.holes) {
					step.Name = lang.Name
				}

				for _, c := range completedSteps {
					if c.Hole == hole.ID && c.Lang == lang.ID {
						step.Complete = true
						break
					}
				}

				progress.Steps = append(progress.Steps, step)
			}
		}

		// Replace incomplete steps with "???" for "hidden" step cheevos.
		if cheevo.cheevo == "archivist" {
			progress.Steps = slices.DeleteFunc(progress.Steps,
				func(s Step) bool { return !s.Complete })

			progress.Steps = append(progress.Steps, slices.Repeat(
				[]Step{{Name: "???"}},
				config.CheevoByID[cheevo.cheevo].Target-len(progress.Steps),
			)...)
		}

		progress.Progress = len(completedSteps)
		data[cheevo.cheevo] = progress
	}

	cheevoProgress(
		`SELECT pangramglot(array_agg(DISTINCT lang))
		   FROM stable_passing_solutions
		  WHERE hole = 'pangram-grep' AND user_id = $1`,
		[]string{"pangramglot"},
	)

	cheevoProgress(
		`SELECT COALESCE(EXTRACT(days FROM TIMEZONE('UTC', NOW()) - MIN(submitted)), 0)
		   FROM stable_passing_solutions
		  WHERE user_id = $1`,
		[]string{"aged-like-fine-wine"},
	)

	cheevoProgress(
		`WITH langs AS (
		    SELECT COUNT(DISTINCT lang)
		      FROM stable_passing_solutions
		     WHERE user_id = $1
		  GROUP BY hole
		) SELECT COALESCE(MAX(count), 0) FROM langs`,
		[]string{"polyglot", "polyglutton", "omniglot", "omniglutton"},
	)

	cheevoProgress(
		`WITH distinct_holes AS (
		    SELECT DISTINCT hole
		      FROM stable_passing_solutions
		     WHERE user_id = $2
		) SELECT COUNT(DISTINCT $1::hstore->hole::text) FROM distinct_holes`,
		[]string{"smörgåsbord"},
		config.HoleCategoryHstore,
	)

	cheevoProgress(
		`SELECT COUNT(DISTINCT hole)
		   FROM stable_passing_solutions
		  WHERE user_id = $1`,
		cheevoIDs(config.CheevoTree["Total Holes"]),
	)

	cheevoProgress(
		"SELECT COALESCE(MAX(points), 0) FROM points WHERE user_id = $1",
		cheevoIDs(config.CheevoTree["Total Points"]),
	)

	render(w, r, "golfer/cheevos", data, golfer.Name)
}
