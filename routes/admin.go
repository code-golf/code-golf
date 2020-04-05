package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/code-golf/code-golf/hole"
)

type solution struct {
	code    string
	Golfer  string        `json:"golfer"`
	HoleID  string        `json:"hole"`
	LangID  string        `json:"lang"`
	Failing bool          `json:"failing"`
	Pass    bool          `json:"pass"`
	Took    time.Duration `json:"took"`
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
				score := hole.Play(r.Context(), s.HoleID, s.LangID, s.code)

				s.Pass = score.Pass
				s.Took = score.Took

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
			` SELECT code, failing, login, hole, lang
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
				&s.HoleID,
				&s.LangID,
			); err != nil {
				panic(err)
			}

			solutions <- s
		}

		if err := rows.Err(); err != nil && err != context.Canceled {
			panic(err)
		}
	}()

	return solutions
}
