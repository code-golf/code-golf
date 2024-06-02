package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/buildkite/terminal-to-html/v3"
	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/null"
	"github.com/code-golf/code-golf/session"
)

type solution struct {
	Code     string        `json:"-"`
	Failing  bool          `json:"-"`
	Golfer   string        `json:"golfer"`
	GolferID int           `json:"golfer_id"`
	HoleID   string        `json:"hole"`
	LangID   string        `json:"lang"`
	Pass     bool          `json:"pass"`
	Stderr   string        `json:"stderr"`
	Took     time.Duration `json:"took"`
	Total    int           `json:"total"`
}

// GET /admin/solutions
func adminSolutionsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Holes []*config.Hole
		Langs []*config.Lang
	}{config.HoleList, config.LangList}

	render(w, r, "admin/solutions", data, "Admin Solutions")
}

// GET /admin/solutions/run
func adminSolutionsRunGET(w http.ResponseWriter, r *http.Request) {
	db := session.Database(r)
	solutions := getSolutions(r)

	var mux sync.Mutex
	var wg sync.WaitGroup

	w.Header().Set("Content-Type", "application/x-ndjson")

	for range 3 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			i := 0
			for s := range solutions {
				// Run each solution up to three times.
				for j := 0; j < 3; j++ {
					// Get the first failing (or last overall) run.
					var run hole.Run
					for _, r := range hole.Play(
						r.Context(),
						config.AllHoleByID[s.HoleID],
						config.AllLangByID[s.LangID],
						s.Code,
					) {
						run = r
						if !r.Pass {
							break
						}
					}

					s.Stderr = terminal.Render([]byte(run.Stderr))
					s.Took = run.Time

					if run.Pass {
						s.Pass = true
						break
					}
				}

				// If the last run differs from the DB, update the database.
				//
				// NOTE It's a little confusing that present is called pass
				//      but past is called failing, so == is a mismatch.
				if s.Pass == s.Failing {
					db.MustExec(
						`UPDATE solutions
						    SET failing = $1
						  WHERE code    = $2
						    AND hole    = $3
						    AND lang    = $4
						    AND user_id = $5`,
						!s.Pass,
						s.Code,
						s.HoleID,
						s.LangID,
						s.GolferID,
					)
				}

				b, err := json.Marshal(s)
				if err != nil {
					panic(err)
				}

				b = append(b, '\n')

				mux.Lock()
				w.Write(b)
				if i++; i%10 == 0 {
					w.(http.Flusher).Flush()
				}
				mux.Unlock()
			}
		}()
	}

	wg.Wait()
}

func getSolutions(r *http.Request) chan solution {
	solutions := make(chan solution)

	go func() {
		defer close(solutions)

		holeID := r.FormValue("hole")
		langID := r.FormValue("lang")

		rows, err := session.Database(r).QueryxContext(
			r.Context(),
			`WITH distinct_solutions AS (
			  SELECT DISTINCT code, failing, login golfer, user_id golfer_id,
			                  hole hole_id, lang lang_id
			    FROM solutions
			    JOIN users   ON id = user_id
			   WHERE failing IN (true, $1)
			     AND (login = $2 OR $2 = '')
			     AND (hole  = $3 OR $3 IS NULL)
			     AND (lang  = $4 OR $4 IS NULL)
			ORDER BY hole, lang, login
			) SELECT *, COUNT(*) OVER () total FROM distinct_solutions`,
			r.FormValue("failing") == "on",
			r.FormValue("golfer"),
			null.New(holeID, holeID != ""),
			null.New(langID, langID != ""),
		)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var s solution
			if err := rows.StructScan(&s); err != nil {
				panic(err)
			}

			solutions <- s
		}

		if err := rows.Err(); err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	return solutions
}
