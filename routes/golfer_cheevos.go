package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
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

	for _, cheevo := range config.CheevoTree["Hole/Lang Specific"] {
		if len(cheevo.Holes) == 0 || len(cheevo.Langs) == 0 {
			continue
		}

		progress := data[cheevo.ID]

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
			pq.Array(cheevo.Holes),
			pq.Array(cheevo.Langs),
			golfer.ID,
		); err != nil {
			panic(err)
		}

		for _, hole := range cheevo.Holes {
			for _, lang := range cheevo.Langs {
				step := Step{Path: "/" + hole.ID + "#" + lang.ID}

				step.Name = hole.Name
				if len(cheevo.Langs) > len(cheevo.Holes) {
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
		if cheevo.ID == "archivist" {
			progress.Steps = slices.DeleteFunc(progress.Steps,
				func(s Step) bool { return !s.Complete })

			progress.Steps = append(progress.Steps, slices.Repeat(
				[]Step{{Name: "???"}},
				cheevo.Target-len(progress.Steps),
			)...)
		}

		progress.Progress = len(completedSteps)
		data[cheevo.ID] = progress
	}

	// Holes of the Week streaks.
	{
		var weeks []time.Time
		if err := db.Select(
			&weeks,
			` SELECT week
			    FROM weekly_solves
			   WHERE user_id = $1 AND week >= this_week() - 7 * 7
			ORDER BY week DESC`,
			golfer.ID,
		); err != nil {
			panic(err)
		}

		var currentStreak int
		for i, week := range weeks {
			// To be current the streak must be up to this week or last.
			if i == 0 {
				thisWeek := config.ThisWeek()

				if !(week.Equal(thisWeek) || week.Equal(thisWeek.AddDate(0, 0, -7))) {
					break
				}

				// To continue the streak must be a week after the previous.
			} else if !week.Equal(weeks[i-1].AddDate(0, 0, -7)) {
				break
			}

			currentStreak++
		}

		for _, cheevo := range config.CheevoTree["Holes of the Week"] {
			progress := data[cheevo.ID]
			progress.Progress = currentStreak
			data[cheevo.ID] = progress
		}
	}

	if progress := data["count-to-ten"]; progress.Earned == nil {
		var completedSteps pq.Int32Array
		var hole *config.Hole
		if err := db.QueryRow(
			` SELECT hole, array_agg(DISTINCT LENGTH(langs.name))
			    FROM stable_passing_solutions
			    JOIN holes ON hole = holes.id
			    JOIN langs ON lang = langs.id
			   WHERE category = 'Sequence'
			     AND LENGTH(langs.name) <= 10
			     AND user_id = $1
			GROUP BY hole
			ORDER BY COUNT(DISTINCT LENGTH(langs.name)) DESC, hole
			   LIMIT 1`,
			golfer.ID,
		).Scan(&hole, &completedSteps); err != nil && !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}

		for i := 1; i <= 10; i++ {
			step := Step{
				Name:     strconv.Itoa(i),
				Complete: slices.Contains(completedSteps, int32(i)),
			}

			if hole != nil {
				step.Path = "/" + hole.ID
			}

			progress.Steps = append(progress.Steps, step)
		}

		progress.Progress = len(completedSteps)
		data["count-to-ten"] = progress
	}

	if progress := data["never-eat-shredded-wheat"]; progress.Earned == nil {
		var completedLangs []*config.Lang
		if err := db.Select(
			&completedLangs,
			`WITH letters AS (
			    SELECT substr(lang::text, 1, 1) lang_letter, lang
			      FROM stable_passing_solutions
			     WHERE hole = 'arrows' AND user_id = $1
			) SELECT DISTINCT ON (lang_letter) lang
			    FROM letters
			   WHERE lang_letter IN ('n', 'e', 's', 'w')
			ORDER BY lang_letter, lang`,
			golfer.ID,
		); err != nil {
			panic(err)
		}

		for _, c := range []string{"N", "E", "S", "W"} {
			step := Step{Name: c + "..."}

			for _, lang := range completedLangs {
				if strings.ToUpper(lang.Name[:1]) == c {
					step.Complete = true
					step.Name = lang.Name
					break
				}
			}

			progress.Steps = append(progress.Steps, step)
		}

		progress.Progress = len(completedLangs)
		data["never-eat-shredded-wheat"] = progress
	}

	if progress := data["pangramglot"]; progress.Earned == nil {
		var completedSteps []string
		if err := db.Select(
			&completedSteps,
			`SELECT unnest(letters(array_agg(DISTINCT lang)))
			   FROM stable_passing_solutions
			  WHERE hole = 'pangram-grep' AND user_id = $1`,
			golfer.ID,
		); err != nil {
			panic(err)
		}

		for c := 'A'; c <= 'Z'; c++ {
			progress.Steps = append(progress.Steps, Step{
				Name:     string(c),
				Complete: slices.Contains(completedSteps, string(c)),
			})
		}

		progress.Progress = len(completedSteps)
		data["pangramglot"] = progress
	}

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

	if progress := data["smörgåsbord"]; progress.Earned == nil {
		if err := db.Select(
			&progress.Steps,
			`  SELECT category name, bool_or(s.hole IS NOT NULL) complete
			     FROM holes
			LEFT JOIN stable_passing_solutions s
			       ON s.hole = id AND s.user_id = $1
			 GROUP BY category
			 ORDER BY category`,
			golfer.ID,
		); err != nil {
			panic(err)
		}

		for _, step := range progress.Steps {
			if step.Complete {
				progress.Progress++
			}
		}

		data["smörgåsbord"] = progress
	}

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
