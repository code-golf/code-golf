package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	Tested   time.Time     `json:"tested"`
	Took     time.Duration `json:"took"`
	Total    int           `json:"total"`
}

// GET /admin/solutions
func adminSolutionsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Holes                []*config.Hole
		Langs                []*config.Lang
		TestedFrom, TestedTo time.Time
	}{Holes: config.AllHoleList, Langs: config.AllLangList}

	if err := session.Database(r).Get(
		&data,
		`SELECT MIN(TIMEZONE($1, TIMEZONE('UTC', tested))) tested_from,
		        MAX(TIMEZONE($1, TIMEZONE('UTC', tested))) tested_to
		   FROM solutions`,
		session.Golfer(r).TimeZone,
	); err != nil {
		panic(err)
	}

	render(w, r, "admin/solutions", data, "Admin Solutions")
}

// GET /admin/solutions/run
func adminSolutionsRunGET(w http.ResponseWriter, r *http.Request) {
	db := session.Database(r)
	solutions := getSolutions(r)

	var mux sync.Mutex
	var wg sync.WaitGroup

	w.Header().Set("Content-Type", "application/x-ndjson")

	noNewFailures := r.FormValue("no-new-failures") == "on"
	workers, _ := strconv.Atoi(r.FormValue("workers"))
	for range workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for s := range solutions {
				// Run each solution up to three times.
				for range 3 {
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

					s.Stderr = run.Stderr
					s.Took = run.Time

					if run.Pass {
						s.Pass = true
						break
					}
				}

				// If we passed, or we're okay saving failures, or we used to
				// fail then save to at least update the tested time.
				if s.Pass || !noNewFailures || s.Failing {
					db.MustExec(
						`UPDATE solutions
						    SET failing = $1, tested = DEFAULT
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
				w.(http.Flusher).Flush()
				mux.Unlock()
			}
		}()
	}

	wg.Wait()
}

// GET /admin/solutions/{hole}/{lang}/{golferID}
func adminSolutionGET(w http.ResponseWriter, r *http.Request) {
	var data []string

	if err := session.Database(r).Select(
		&data,
		"SELECT DISTINCT code FROM solutions WHERE hole = $1 AND lang = $2 AND user_id = $3",
		param(r, "hole"),
		param(r, "lang"),
		param(r, "golferID"),
	); err != nil {
		panic(err)
	}

	render(w, r, "admin/solution", data, "Admin Solution")
}

func getSolutions(r *http.Request) chan solution {
	solutions := make(chan solution)

	go func() {
		defer close(solutions)

		rows, err := session.Database(r).QueryxContext(
			r.Context(),
			`WITH distinct_solutions AS (
			  SELECT DISTINCT code, failing, login golfer, user_id golfer_id,
			                  hole hole_id, lang lang_id, tested
			    FROM solutions
			    JOIN users   ON id = user_id
			   WHERE failing IN (true, $1)
			     AND (login = $2 OR $2 = '')
			     AND (hole  = $3 OR $3 IS NULL)
			     AND (lang  = $4 OR $4 IS NULL)
			     AND DATE(TIMEZONE($5, TIMEZONE('UTC', tested))) BETWEEN $6 AND $7
			ORDER BY tested
			) SELECT *, COUNT(*) OVER () total FROM distinct_solutions`,
			r.FormValue("failing") == "on",
			r.FormValue("golfer"),
			null.NullIfZero(r.FormValue("hole")),
			null.NullIfZero(r.FormValue("lang")),
			session.Golfer(r).TimeZone,
			r.FormValue("tested-from"),
			r.FormValue("tested-to"),
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
