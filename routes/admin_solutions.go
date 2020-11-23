package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

type solution struct {
	code     string
	codeID   int
	failing  bool
	Golfer   string        `json:"golfer"`
	GolferID int           `json:"golfer_id"`
	HoleID   string        `json:"hole"`
	LangID   string        `json:"lang"`
	Pass     bool          `json:"pass"`
	Stderr   string        `json:"stderr"`
	Took     time.Duration `json:"took"`
}

// AdminSolutions serves GET /admin/solutions
func AdminSolutions(w http.ResponseWriter, r *http.Request) {
	render(w, r, "admin/solutions", "Admin Solutions", struct {
		Holes []hole.Hole
		Langs []lang.Lang
	}{hole.List, lang.List})
}

// AdminSolutionsRun serves GET /admin/solutions/run
func AdminSolutionsRun(w http.ResponseWriter, r *http.Request) {
	db := session.Database(r)
	solutions := getSolutions(r)

	var mux sync.Mutex
	var wg sync.WaitGroup

	w.Header().Set("Content-Type", "application/x-ndjson")

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for s := range solutions {
				var passes, fails int

				// Best of three runs.
				for j := 0; passes < 2 && fails < 2; j++ {
					score := hole.Play(r.Context(), s.HoleID, s.LangID, s.code)

					s.Took = score.Took

					if score.Pass {
						passes++
					} else {
						fails++
						s.Stderr = string(score.Stderr)
					}
				}

				s.Pass = passes > fails

				// If the last run differs from the DB, update the database.
				//
				// NOTE It's a little confusing that present is called pass
				//      but past is called failing, so == is a mismatch.
				if s.Pass == s.failing {
					if _, err := db.Exec(
						`UPDATE solutions
						    SET failing = $1
						  WHERE code_id = $2
						    AND hole    = $3
						    AND lang    = $4
						    AND user_id = $5`,
						!s.Pass,
						s.codeID,
						s.HoleID,
						s.LangID,
						s.GolferID,
					); err != nil {
						panic(err)
					}
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

func getSolutions(r *http.Request) chan solution {
	solutions := make(chan solution)

	go func() {
		defer close(solutions)

		holeID := r.FormValue("hole")
		langID := r.FormValue("lang")

		rows, err := session.Database(r).QueryContext(
			r.Context(),
			` SELECT DISTINCT code, code_id, failing, login, u.id, hole, lang
			    FROM solutions
			    JOIN code  c ON c.id = code_id
			    JOIN users u ON u.id = user_id
			   WHERE failing IN (true, $1)
			     AND (login = $2 OR $2 = '')
			     AND (hole  = $3 OR $3 IS NULL)
			     AND (lang  = $4 OR $4 IS NULL)
			ORDER BY hole, lang, login`,
			r.FormValue("failing") == "on",
			r.FormValue("golfer"),
			sql.NullString{String: holeID, Valid: holeID != ""},
			sql.NullString{String: langID, Valid: langID != ""},
		)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		for rows.Next() {
			var s solution

			if err := rows.Scan(
				&s.code,
				&s.codeID,
				&s.failing,
				&s.Golfer,
				&s.GolferID,
				&s.HoleID,
				&s.LangID,
			); err != nil {
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
