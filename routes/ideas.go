package routes

import (
	"net/http"
	"sort"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /ideas
func ideasGET(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		Hole                     *config.Hole
		ID, ThumbsDown, ThumbsUp int
		Title                    string
	}

	data := struct {
		Holes, Ideas, Langs []idea
	}{Holes: make([]idea, len(config.ExpHoleList))}

	for i, hole := range config.ExpHoleList {
		data.Holes[i] = idea{Hole: hole, ID: hole.Experiment, Title: hole.Name}
	}

	rows, err := session.Database(r).Query(
		"SELECT * FROM ideas ORDER BY thumbs_up - thumbs_down DESC, title")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

rows:
	for rows.Next() {
		var i idea
		if err := rows.Scan(&i.ID, &i.ThumbsDown, &i.ThumbsUp, &i.Title); err != nil {
			panic(err)
		}

		for j, hole := range data.Holes {
			if hole.ID == i.ID {
				data.Holes[j].ThumbsDown = i.ThumbsDown
				data.Holes[j].ThumbsUp = i.ThumbsUp
				continue rows
			}
		}

		if strings.Contains(strings.ToLower(i.Title), "lang") {
			data.Langs = append(data.Langs, i)
		} else {
			data.Ideas = append(data.Ideas, i)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	sort.SliceStable(data.Holes, func(i, j int) bool {
		return data.Holes[i].ThumbsUp-data.Holes[i].ThumbsDown >
			data.Holes[j].ThumbsUp-data.Holes[j].ThumbsDown
	})

	render(w, r, "ideas", data, "Ideas")
}
