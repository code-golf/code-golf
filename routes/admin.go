package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/code-golf/code-golf/hole"
)

type solution struct {
	code     string
	Golfer   string           `json:"golfer"`
	GolferID int              `json:"golfer_id"`
	HoleID   string           `json:"hole"`
	LangID   string           `json:"lang"`
	Failing  bool             `json:"failing"`
	Scores   []hole.Scorecard `json:"scores"`
}

// Admin serves GET /admin
func Admin(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "admin", "Admin", struct {
		Holes []Hole
		Langs []Lang
	}{holes, langs})
}

// AdminSolutions serves GET /admin/solutions
func AdminSolutions(w http.ResponseWriter, r *http.Request) {
	solutions := getSolutions(
		r.Context(),
		r.FormValue("golfer"),
		r.FormValue("hole"),
		r.FormValue("lang"),
	)

	var mux sync.Mutex
	var wg sync.WaitGroup

	w.Header().Set("Content-Type", "application/x-ndjson")

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for s := range solutions {
				// Run the solution upto three times, bail early if it passes.
				for j := 0; j < 3; j++ {
					s.Scores = append(
						s.Scores,
						hole.Play(r.Context(), s.HoleID, s.LangID, s.code),
					)

					if s.Scores[j].Pass {
						break
					}
				}

				// If the last run differs from the DB, update the database.
				//
				// NOTE It's a little confusing that present is called pass
				//      but past is called failing, so == is a mismatch.
				if pass := s.Scores[len(s.Scores)-1].Pass; pass == s.Failing {
					s.Failing = !pass

					if _, err := r.Context().Value("db").(*sql.DB).Exec(
						`UPDATE solutions
						    SET failing = $1
						  WHERE code    = $2
						    AND hole    = $3
						    AND lang    = $4
						    AND user_id = $5`,
						s.Failing,
						s.code,
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

func getSolutions(ctx context.Context, golfer, holeID, langID string) chan solution {
	solutions := make(chan solution)

	go func() {
		defer close(solutions)

		rows, err := ctx.Value("db").(*sql.DB).QueryContext(
			ctx,
			` SELECT code, failing, login, id, hole, lang
				FROM solutions
				JOIN users ON id = user_id
			   WHERE (login = $1 OR $1 = '')
				 AND (hole  = $2 OR $2 IS NULL)
				 AND (lang  = $3 OR $3 IS NULL)
			ORDER BY hole, lang, login`,
			golfer,
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
				&s.Failing,
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
