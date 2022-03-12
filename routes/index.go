package routes

import (
	"net/http"
	"sort"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /
func indexGET(w http.ResponseWriter, r *http.Request) {
	for _, param := range []string{"scoring", "sort"} {
		if value := r.FormValue(param); value != "" {
			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				Name:     "__Host-" + param,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
				Value:    value,
			})

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	type Card struct {
		Hole   *config.Hole
		Lang   *config.Lang
		Points int
	}

	data := struct {
		Cards         []Card
		Scoring, Sort string
		Sorts         []string
	}{
		Cards:   make([]Card, 0, len(config.HoleList)),
		Scoring: "bytes",
		Sorts:   []string{"alphabetical", "category", "points"},
	}

	if golfer := session.Golfer(r); golfer == nil {
		for _, hole := range config.HoleList {
			data.Cards = append(data.Cards, Card{Hole: hole})
		}
	} else {
		if cookie(r, "__Host-scoring") == "chars" {
			data.Scoring = "chars"
		}

		rows, err := session.Database(r).Query(
			`WITH points AS (
			    SELECT DISTINCT ON (hole) hole, lang, points
			      FROM rankings
			     WHERE scoring = $1 AND user_id = $2
			  ORDER BY hole, points DESC, lang
			)  SELECT hole, COALESCE(lang::text, ''), COALESCE(points, 0)
			     FROM unnest(enum_range(NULL::hole)) hole
			LEFT JOIN points USING(hole)`,
			data.Scoring,
			golfer.ID,
		)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var card Card
			var holeID, langID string

			if err := rows.Scan(&holeID, &langID, &card.Points); err != nil {
				panic(err)
			}

			if hole, ok := config.HoleByID[holeID]; ok {
				card.Hole = hole
				card.Lang = config.LangByID[langID]

				data.Cards = append(data.Cards, card)
			}
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}

		switch data.Sort = cookie(r, "__Host-sort"); data.Sort {
		case "alphabetical-desc":
			sort.Slice(data.Cards, func(i, j int) bool {
				return strings.ToLower(data.Cards[i].Hole.Name) >
					strings.ToLower(data.Cards[j].Hole.Name)
			})
		case "category-asc":
			sort.Slice(data.Cards, func(i, j int) bool {
				if data.Cards[i].Hole.Category == data.Cards[j].Hole.Category {
					return strings.ToLower(data.Cards[i].Hole.Name) <
						strings.ToLower(data.Cards[j].Hole.Name)
				}
				return data.Cards[i].Hole.Category < data.Cards[j].Hole.Category
			})
		case "category-desc":
			sort.Slice(data.Cards, func(i, j int) bool {
				if data.Cards[i].Hole.Category == data.Cards[j].Hole.Category {
					return strings.ToLower(data.Cards[i].Hole.Name) >
						strings.ToLower(data.Cards[j].Hole.Name)
				}
				return data.Cards[i].Hole.Category > data.Cards[j].Hole.Category
			})
		case "points-asc":
			sort.Slice(data.Cards, func(i, j int) bool {
				if data.Cards[i].Points == data.Cards[j].Points {
					return strings.ToLower(data.Cards[i].Hole.Name) <
						strings.ToLower(data.Cards[j].Hole.Name)
				}
				return data.Cards[i].Points < data.Cards[j].Points
			})
		case "points-desc":
			sort.Slice(data.Cards, func(i, j int) bool {
				if data.Cards[i].Points == data.Cards[j].Points {
					return strings.ToLower(data.Cards[i].Hole.Name) >
						strings.ToLower(data.Cards[j].Hole.Name)
				}
				return data.Cards[i].Points > data.Cards[j].Points
			})
		default:
			data.Sort = "alphabetical-asc"

			sort.Slice(data.Cards, func(i, j int) bool {
				return strings.ToLower(data.Cards[i].Hole.Name) <
					strings.ToLower(data.Cards[j].Hole.Name)
			})
		}
	}

	w.Header().Set(
		"Strict-Transport-Security",
		"max-age=31536000;includeSubDomains;preload",
	)
	render(w, r, "index", data)
}
