package routes

import (
	"net/http"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// Ideas serves GET /ideas
func Ideas(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		ID, ThumbsDown, ThumbsUp int
		Title                    string
	}

	data := struct {
		Holes            []*config.Hole
		Ideas, Languages []idea
	}{Holes: config.ExpHoleList}

	rows, err := session.Database(r).Query(
		r.Context(),
		"SELECT * FROM ideas ORDER BY thumbs_up - thumbs_down DESC, title",
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var i idea

		if err := rows.Scan(&i.ID, &i.ThumbsDown, &i.ThumbsUp, &i.Title); err != nil {
			panic(err)
		}
		if strings.Contains(strings.ToLower(i.Title), "lang") {
			data.Languages = append(data.Languages, i)
		} else {
			data.Ideas = append(data.Ideas, i)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "ideas", data, "Ideas")
}
