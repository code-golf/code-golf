package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

type solution struct {
	failing  bool          `json:"-"`
	Pass     bool          `json:"pass"`
	codeID   int           `json:"-"`
	GolferID int           `json:"golfer_id"`
	code     string        `json:"-"`
	Golfer   string        `json:"golfer"`
	HoleID   string        `json:"hole"`
	LangID   string        `json:"lang"`
	Stderr   string        `json:"stderr"`
	Took     time.Duration `json:"took"`
}

func testCode(code, holeID string, r *http.Request) bool {
	var passes, fails int

	// Best of three runs.
	for j := 0; passes < 2 && fails < 2; j++ {
		score := hole.Play(r.Context(), holeID, "c-sharp", code)
		if score.Pass {
			passes++
		} else {
			fails++
		}
	}

	return passes > fails
}

func substitueArgs(originalCode, code, args, holeID string, r *http.Request) string {
	re := regexp.MustCompile(`\b` + args + `\b`)
	resultCode := originalCode

	// Try first of two strategies.
	code1 := re.ReplaceAllString(code, "args")
	if testCode(code1, holeID, r) && len(code1) < len(resultCode) {
		resultCode = code1
	}

	// Second strategy: declare new variable and assign args to it.
	// The variable would have to be used in at least four places for this to help.
	// Ex: var a=args;aaaa
	// Vs: argsargsargsargs
	code2 := "var " + args + "=args;" + code
	if testCode(code2, holeID, r) && len(code2) < len(resultCode) {
		resultCode = code2
	}

	return resultCode
}

// If the code can be shortened automatically and the new code passes the tests, returns the new code.
// Otherwise, returns the original code.
func updateCode(originalCode, holeID string, recursionLevel int, r *http.Request) string {
	// Find the Class
	regex := regexp.MustCompile(`(public\s+|internal\s+)?class\s+\w+\s*({{?)`)
	match := regex.FindStringIndex(originalCode)

	if match == nil {
		// Can't find class
		return originalCode
	}

	submatch := regex.FindStringSubmatch(originalCode)
	// If this is in a Quine, the { } may be doubled to escape them.
	endingThing := strings.Repeat("}", len(submatch[2]))

	start := match[0]
	end := match[1]
	code1 := originalCode[:start] + originalCode[end:]
	endIndex := strings.LastIndex(code1, endingThing)
	if endIndex == -1 {
		// Can't find end of class
		return originalCode
	}

	code1 = code1[:endIndex] + code1[endIndex+len(endingThing):]

	var code2, argsName string

	// Maybe Main uses an expression body.
	regex = regexp.MustCompile(`(public\s+|private\s+)?static\s+(public\s+|private\s+)?void\s+Main\s*\(\s*(string\[\]\s*(?P<args>\w+))?\s*\)\s*=>\s*`)
	match = regex.FindStringIndex(code1)
	if match != nil {
		start = match[0]
		end = match[1]
		code2 = code1[:start] + code1[end:]
		argsName = regex.FindStringSubmatch(code1)[4]
	} else {
		// Find Main with opening square backets
		regex = regexp.MustCompile(`(public\s+|private\s+)?static\s+(public\s+|private\s+)?void\s+Main\s*\(\s*(string\[\]\s*(?P<args>\w+))?\s*\)\s*({{?)\s*`)
		match = regex.FindStringIndex(code1)

		if match == nil {
			// Can't find Main
			return originalCode
		}

		submatch = regex.FindStringSubmatch(code1)
		argsName = submatch[4]

		// If this is in a Quine, the { } may be doubled to escape them.
		endingThing = strings.Repeat("}", len(submatch[5]))

		start = match[0]
		end = match[1]
		code2 = code1[:start] + code1[end:]
		endIndex = strings.LastIndex(code2, endingThing)
		if endIndex == -1 {
			// Can't find end of Main
			return originalCode
		}

		code2 = code2[:endIndex] + code2[endIndex+len(endingThing):]
	}

	resultCode := originalCode

	if len(argsName) > 0 {
		resultCode = substitueArgs(originalCode, code2, argsName, holeID, r)
	} else if testCode(code2, holeID, r) {
		resultCode = code2
	} else if holeID == "quine" && recursionLevel == 0 {
		// If we are in a quine, may need to repeat once to fix it.
		code3 := updateCode(code2, holeID, recursionLevel+1, r)
		if code3 != code2 {
			resultCode = code3
		}
	}

	return resultCode
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
				if s.LangID != "c-sharp" {
					continue
				}

				newCode := updateCode(s.code, s.HoleID, 0, r)
				if s.code != newCode {
					// Indicate that it was updated.
					s.Pass = true

					// update the database.
					if _, err := db.Exec(
						`SELECT save_solution(code := $1, hole := $2, lang := $3, user_id := $4)`,
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
