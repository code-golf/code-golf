package routes

import (
	"cmp"
	"net/http"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

type Card struct {
	Hole   *config.Hole
	Lang   *config.Lang
	Points int
}

var cardList, expCardList []Card

func init() {
	for _, hole := range config.HoleList {
		cardList = append(cardList, Card{Hole: hole})
	}
	for _, hole := range config.ExpHoleList {
		expCardList = append(expCardList, Card{Hole: hole})
	}
}

// Get prev/next hole with order based on homepage settings.
func getPrevNextHole(r *http.Request, hole *config.Hole) (prev, next *config.Hole) {
	cards := expCardList
	if hole.Experiment == 0 {
		cards = getHomeCards(r)
	}
	i := slices.IndexFunc(cards, func(c Card) bool { return c.Hole.ID == hole.ID })

	if i == 0 {
		prev = cards[len(cards)-1].Hole
	} else {
		prev = cards[i-1].Hole
	}

	if i == len(cards)-1 {
		next = cards[0].Hole
	} else {
		next = cards[i+1].Hole
	}

	return
}

// Get homepage cards with order based on homepage settings.
func getHomeCards(r *http.Request) (cards []Card) {
	golfer := session.Golfer(r)
	if golfer == nil {
		return cardList
	}

	var query string
	var bind []any

	if lang := golfer.Settings["home"]["points-for"]; lang == "all" {
		query = `WITH points AS (
			   SELECT DISTINCT ON (hole) hole, lang, points
			     FROM rankings
			    WHERE scoring = $1 AND user_id = $2
			 ORDER BY hole, points DESC, lang
			)  SELECT id hole, lang, COALESCE(points, 0) points
			     FROM holes
			LEFT JOIN points ON id = hole WHERE experiment = 0`

		bind = []any{golfer.Settings["home"]["scoring"], golfer.ID}
	} else {
		query = `WITH points AS (
			   SELECT hole, lang, points_for_lang
			     FROM rankings
			    WHERE scoring = $1 AND user_id = $2 AND lang = $3
			)  SELECT id hole, lang, COALESCE(points_for_lang, 0) points
			     FROM holes
			LEFT JOIN points ON id = hole WHERE experiment = 0`

		bind = []any{golfer.Settings["home"]["scoring"], golfer.ID, lang}
	}

	if err := session.Database(r).Select(&cards, query, bind...); err != nil {
		panic(err)
	}

	cmpHoleNameLowercase := func(a, b Card) int {
		return cmp.Compare(strings.ToLower(a.Hole.Name),
			strings.ToLower(b.Hole.Name))
	}

	switch golfer.Settings["home"]["order-by"] {
	case "alphabetical-asc": // name desc.
		slices.SortFunc(cards, cmpHoleNameLowercase)
	case "alphabetical-desc": // name desc.
		slices.SortFunc(cards, func(a, b Card) int {
			return cmpHoleNameLowercase(b, a)
		})
	case "category-asc": // category asc, name asc.
		slices.SortFunc(cards, func(a, b Card) int {
			if c := cmp.Compare(a.Hole.Category, b.Hole.Category); c != 0 {
				return c
			}
			return cmpHoleNameLowercase(a, b)
		})
	case "category-desc": // category desc, name asc.
		slices.SortFunc(cards, func(a, b Card) int {
			if c := cmp.Compare(b.Hole.Category, a.Hole.Category); c != 0 {
				return c
			}
			return cmpHoleNameLowercase(a, b)
		})
	case "points-asc": // points asc, name asc.
		slices.SortFunc(cards, func(a, b Card) int {
			if c := cmp.Compare(a.Points, b.Points); c != 0 {
				return c
			}
			return cmpHoleNameLowercase(a, b)
		})
	case "points-desc": // points desc, name asc.
		slices.SortFunc(cards, func(a, b Card) int {
			if c := cmp.Compare(b.Points, a.Points); c != 0 {
				return c
			}
			return cmpHoleNameLowercase(a, b)
		})
	case "released-asc": // released asc, name asc.
		slices.SortFunc(cards, func(a, b Card) int {
			if c := cmp.Compare(
				a.Hole.Released.String(), b.Hole.Released.String(),
			); c != 0 {
				return c
			}
			return cmpHoleNameLowercase(a, b)
		})
	case "released-desc": // released desc, name asc.
		slices.SortFunc(cards, func(a, b Card) int {
			if c := cmp.Compare(
				b.Hole.Released.String(), a.Hole.Released.String(),
			); c != 0 {
				return c
			}
			return cmpHoleNameLowercase(a, b)
		})
	}

	return
}
