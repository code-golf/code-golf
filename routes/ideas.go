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

	golfersByHole := ideaGolfers(r, "hole", holeIDs)
	golfersByLang := ideaGolfers(r, "lang", langIDs)

	// Experimental holes the golfer has a passing solution for.
	var solvedHoles map[string]bool
	if golfer := session.Golfer(r); golfer != nil {
		solvedHoles = ideaSolvedHoles(r, golfer.ID)
	}

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
				i.Golfers = golfersByHole[hole.ID]
				i.Hole = hole
				i.Solved = solvedHoles[hole.ID]
				i.Title = hole.Name
				data.Holes = append(data.Holes, i)
			}
		} else if lang, ok := langs[i.ID]; ok {
			i.Golfers = golfersByLang[lang.ID]
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

// How many golfers have a passing solution for each of the given holes or
// langs, column is either "hole" or "lang".
func ideaGolfers(r *http.Request, column string, ids []string) map[string]int {
	var rows []struct {
		Golfers int
		ID      string
	}
	if err := session.Database(r).Select(
		&rows,
		`SELECT COUNT(DISTINCT user_id) golfers, `+column+` id
		   FROM solutions
		  WHERE `+column+` = ANY($1) AND NOT failing
	   GROUP BY `+column,
		pq.Array(ids),
	); err != nil {
		panic(err)
	}

	golfers := map[string]int{}
	for _, row := range rows {
		golfers[row.ID] = row.Golfers
	}
	return golfers
}

// The holes the golfer has a passing solution for.
func ideaSolvedHoles(r *http.Request, golferID int) map[string]bool {
	var holes []string
	if err := session.Database(r).Select(
		&holes,
		"SELECT DISTINCT hole FROM solutions WHERE user_id = $1 AND NOT failing",
		golferID,
	); err != nil {
		panic(err)
	}

	solved := map[string]bool{}
	for _, hole := range holes {
		solved[hole] = true
	}
	return solved
}
