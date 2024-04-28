package routes

import (
	"cmp"
	"database/sql"
	"net/http"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /
func homeGET(w http.ResponseWriter, r *http.Request) {
	type Card struct {
		Hole   *config.Hole
		Lang   *config.Lang
		Points int
	}

	data := struct {
		Cards     []Card
		LangsUsed map[string]bool
		Settings  []*config.Setting
	}{
		Cards:     make([]Card, 0, len(config.HoleList)),
		LangsUsed: map[string]bool{},
		Settings:  config.Settings["home"],
	}

	if golfer := session.Golfer(r); golfer == nil {
		for _, hole := range config.HoleList {
			data.Cards = append(data.Cards, Card{Hole: hole})
		}
	} else {
		var rows *sql.Rows
		var err error

		if lang := golfer.Settings["home"]["points-for"]; lang == "all" {
			rows, err = session.Database(r).Query(
				`WITH points AS (
				    SELECT DISTINCT ON (hole) hole, lang, points
				      FROM rankings
				     WHERE scoring = $1 AND user_id = $2
				  ORDER BY hole, points DESC, lang
				)  SELECT hole, lang, COALESCE(points, 0)
				     FROM unnest(enum_range(NULL::hole)) hole
				LEFT JOIN points USING(hole)`,
				golfer.Settings["home"]["scoring"],
				golfer.ID,
			)
		} else {
			rows, err = session.Database(r).Query(
				`WITH points AS (
				    SELECT hole, lang, points_for_lang
				      FROM rankings
				     WHERE scoring = $1 AND user_id = $2 AND lang = $3
				)  SELECT hole, lang, COALESCE(points_for_lang, 0)
				     FROM unnest(enum_range(NULL::hole)) hole
				LEFT JOIN points USING(hole)`,
				golfer.Settings["home"]["scoring"],
				golfer.ID,
				lang,
			)
		}

		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var card Card
			var hole config.NullHole // NULL when DB updated before config.
			var lang config.NullLang // NULL when not solved.

			if err := rows.Scan(&hole, &lang, &card.Points); err != nil {
				panic(err)
			}

			if lang.Valid {
				data.LangsUsed[lang.Lang.ID] = true
			}

			if hole.Valid {
				card.Hole = hole.Hole
				card.Lang = lang.Lang
				data.Cards = append(data.Cards, card)
			}
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}

		cmpHoleNameLowercase := func(a, b Card) int {
			return cmp.Compare(strings.ToLower(a.Hole.Name),
				strings.ToLower(b.Hole.Name))
		}

		switch golfer.Settings["home"]["order-by"] {
		case "alphabetical-asc": // name desc.
			slices.SortFunc(data.Cards, cmpHoleNameLowercase)
		case "alphabetical-desc": // name desc.
			slices.SortFunc(data.Cards, func(a, b Card) int {
				return cmpHoleNameLowercase(b, a)
			})
		case "category-asc": // category asc, name asc.
			slices.SortFunc(data.Cards, func(a, b Card) int {
				if c := cmp.Compare(a.Hole.Category, b.Hole.Category); c != 0 {
					return c
				}
				return cmpHoleNameLowercase(a, b)
			})
		case "category-desc": // category desc, name asc.
			slices.SortFunc(data.Cards, func(a, b Card) int {
				if c := cmp.Compare(b.Hole.Category, a.Hole.Category); c != 0 {
					return c
				}
				return cmpHoleNameLowercase(a, b)
			})
		case "points-asc": // points asc, name asc.
			slices.SortFunc(data.Cards, func(a, b Card) int {
				if c := cmp.Compare(a.Points, b.Points); c != 0 {
					return c
				}
				return cmpHoleNameLowercase(a, b)
			})
		case "points-desc": // points desc, name asc.
			slices.SortFunc(data.Cards, func(a, b Card) int {
				if c := cmp.Compare(b.Points, a.Points); c != 0 {
					return c
				}
				return cmpHoleNameLowercase(a, b)
			})
		case "released-asc": // released asc, name asc.
			slices.SortFunc(data.Cards, func(a, b Card) int {
				if c := cmp.Compare(
					a.Hole.Released.String(), b.Hole.Released.String(),
				); c != 0 {
					return c
				}
				return cmpHoleNameLowercase(a, b)
			})
		case "released-desc": // released desc, name asc.
			slices.SortFunc(data.Cards, func(a, b Card) int {
				if c := cmp.Compare(
					b.Hole.Released.String(), a.Hole.Released.String(),
				); c != 0 {
					return c
				}
				return cmpHoleNameLowercase(a, b)
			})
		}
	}

	w.Header().Set(
		"Strict-Transport-Security",
		"max-age=31536000;includeSubDomains;preload",
	)
	render(w, r, "home", data)
}
