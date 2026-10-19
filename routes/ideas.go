package routes

import (
	"cmp"
	"net/http"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /ideas
func ideasGET(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		Category, Title          string
		Solved                   bool
		Hole                     *config.Hole
		Golfers                  int
		ID, ThumbsDown, ThumbsUp int
		Lang                     *config.Lang
	}

	var data struct{ Holes, Ideas, Langs []idea }

	// FIXME This is only a slice because of semiprime/sphenic.
	holes := map[int][]*config.Hole{}
	holeIDs := make([]string, len(config.ExpHoleList))
	for i, hole := range config.ExpHoleList {
		holes[hole.Experiment] = append(holes[hole.Experiment], hole)
		holeIDs[i] = hole.ID
	}

	langs := map[int]*config.Lang{}
	langIDs := make([]string, len(config.ExpLangList))
	for i, lang := range config.ExpLangList {
		langs[lang.Experiment] = lang
		langIDs[i] = lang.ID
	}

	// No golfer has an ID of zero, so nothing is solved when logged out.
	golferID := 0
	if golfer := session.Golfer(r); golfer != nil {
		golferID = golfer.ID
	}

	holeSolves := ideaSolves(r, "hole", holeIDs, golferID)
	langSolves := ideaSolves(r, "lang", langIDs, golferID)

	var ideas []idea
	if err := session.Database(r).Select(
		&ideas,
		"SELECT * FROM ideas ORDER BY thumbs_up - thumbs_down DESC, title",
	); err != nil {
		panic(err)
	}

	for _, i := range ideas {
		if holes, ok := holes[i.ID]; ok {
			for _, hole := range holes {
				i.Golfers = holeSolves[hole.ID].Golfers
				i.Hole = hole
				i.Solved = holeSolves[hole.ID].Solved
				i.Title = hole.Name
				data.Holes = append(data.Holes, i)
			}
		} else if lang, ok := langs[i.ID]; ok {
			i.Golfers = langSolves[lang.ID].Golfers
			i.Lang = lang
			i.Title = lang.Name
			data.Langs = append(data.Langs, i)
		} else {
			data.Ideas = append(data.Ideas, i)
		}
	}

	// Sort by vote difference, then case-insensitive title.
	for _, ideas := range [][]idea{data.Holes, data.Langs} {
		slices.SortFunc(ideas, func(a, b idea) int {
			return cmp.Or(
				cmp.Compare(b.ThumbsUp-b.ThumbsDown, a.ThumbsUp-a.ThumbsDown),
				cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title)),
			)
		})
	}

	render(w, r, "ideas", data, "Ideas")
}

type ideaSolve struct {
	Golfers int
	ID      string
	Solved  bool
}

// How many golfers have a passing solution for each of the given holes or
// langs, and whether the golfer is one of them. Column is "hole" or "lang".
func ideaSolves(r *http.Request, column string, ids []string, golferID int) map[string]ideaSolve {
	var rows []ideaSolve
	if err := session.Database(r).Select(
		&rows,
		`SELECT COUNT(DISTINCT user_id) golfers, `+column+` id,
		        bool_or(user_id = $2) solved
		   FROM solutions
		  WHERE `+column+` = ANY($1) AND NOT failing
	   GROUP BY `+column,
		pq.Array(ids),
		golferID,
	); err != nil {
		panic(err)
	}

	solves := map[string]ideaSolve{}
	for _, row := range rows {
		solves[row.ID] = row
	}
	return solves
}
