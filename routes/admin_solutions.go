package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
)

type solution struct {
	failing  bool          `json:"-"`
	Pass     bool          `json:"pass"`
	GolferID int           `json:"golfer_id"`
	code     string        `json:"-"`
	Golfer   string        `json:"golfer"`
	HoleID   string        `json:"hole"`
	LangID   string        `json:"lang"`
	Stderr   string        `json:"stderr"`
	Took     time.Duration `json:"took"`
	Total    int           `json:"total"`
}

func testCode(code, holeID string, r *http.Request) bool {
	var passes, fails int

	// Best of three runs.
	for j := 0; passes < 2 && fails < 2; j++ {
		score := play(r.Context(), holeID, "c-sharp", code)
		if score.Pass {
			passes++
		} else {
			fails++
		}
	}

	return passes > fails
}

// Test each deletion and keep it if it works.
func deleteMatches(code, holeID string, r *http.Request, matches [][]int) string {
	// Start with the last match, so that indices remain valid.
	// Note that matches are non-overlapping.
	for i := len(matches) - 1; i >= 0; i-- {
		start := matches[i][0]
		end := matches[i][1]
		newCode := code[:start] + code[end:]
		if testCode(newCode, holeID, r) {
			code = newCode
		}
	}

	return code
}

// If the code can be shortened automatically and the new code passes the tests, returns the new code.
// Otherwise, returns the original code.
func updateCode(originalCode, holeID string, r *http.Request) string {
	resultCode := originalCode
	var re *regexp.Regexp

	// Try removing all of the instances of "System." at once.
	// For quine, removing them one at a time would not work.
	re = regexp.MustCompile(`System\.`)
	newCode := re.ReplaceAllString(resultCode, "")
	if testCode(newCode, holeID, r) {
		resultCode = newCode
	}

	// // Look for using statements and try removing them.
	re = regexp.MustCompile(`using.*?;`)
	resultCode = deleteMatches(resultCode, holeID, r, re.FindAllStringIndex(resultCode, -1))

	// Look for fully qualified names and try removing the namespace.
	re = regexp.MustCompile(`(\w+\.)+`)
	matches := re.FindAllStringIndex(resultCode, -1)
	for i := len(matches) - 1; i >= 0; i-- {
		start := matches[i][0]
		end := matches[i][1]
		for {
			newCode := resultCode[:start] + resultCode[end:]
			if testCode(newCode, holeID, r) {
				resultCode = newCode
				break
			}

			offset := strings.LastIndex(resultCode[start:end-1], ".")
			if offset >= 0 {
				// Try again, removing everything but the last component.
				end = start + offset + 1
				continue
			}

			break
		}
	}

	return resultCode
}

// AdminSolutions serves GET /admin/solutions
func AdminSolutions(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Holes []*config.Hole
		Langs []*config.Lang
	}{config.HoleList, config.LangList}

	render(w, r, "admin/solutions", data, "Admin Solutions")
}

// Wrap hole.Play so we can recover from panics.
func play(ctx context.Context, holeID, langID, code string) (score hole.Scorecard) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	return hole.Play(ctx, holeID, langID, code)
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

			i := 0
			for s := range solutions {
				if s.LangID != "c-sharp" {
					continue
				}

				newCode := updateCode(s.code, s.HoleID, r)
				if s.code != newCode {
					// Indicate that it was updated.
					s.Pass = true

					// update the database.
					if _, err := db.Exec(
						`SELECT save_solution(bytes := octet_length($1), chars := char_length($1), code := $1, hole := $2, lang := $3, user_id := $4)`,
						newCode,
						s.HoleID,
						s.LangID,
						s.GolferID,
					); err != nil {
						panic(err)
					}
				} else {
					// Indicate it was not updated.
					s.Pass = false
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

		rows, err := session.Database(r).QueryContext(
			r.Context(),
			`WITH distinct_solutions AS (
			  SELECT DISTINCT code, failing, login, user_id, hole, lang
			    FROM solutions
			    JOIN users   ON id = user_id
			   WHERE failing IN (true, $1)
			     AND (login = $2 OR $2 = '')
			     AND (hole  = $3 OR $3 IS NULL)
			     AND (lang  = $4 OR $4 IS NULL)
			ORDER BY hole, lang, login
			) SELECT *, COUNT(*) OVER () FROM distinct_solutions`,
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
				&s.failing,
				&s.Golfer,
				&s.GolferID,
				&s.HoleID,
				&s.LangID,
				&s.Total,
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
