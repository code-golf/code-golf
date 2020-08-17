package routes

import (
	"net/http"
	"sort"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// Index serves GET /
func Index(w http.ResponseWriter, r *http.Request) {
	if sort := r.URL.Query().Get("sort"); sort != "" {
		http.SetCookie(w, &http.Cookie{
			HttpOnly: true,
			Name:     "__Host-sort",
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Secure:   true,
			Value:    sort,
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	type Card struct {
		Golfers, Rank int
		Hole          hole.Hole
	}

	type Fail struct {
		Hole hole.Hole
		Lang lang.Lang
	}

	data := struct {
		Cards []Card
		Fails []Fail
		Sort  string
		Sorts []string
	}{
		Sorts: []string{"alphabetical", "golfers", "rank"},
	}

	if sort, err := r.Cookie("__Host-sort"); err == nil {
		data.Sort = sort.Value
	}

	var userID int
	if golfer := session.Golfer(r); golfer != nil {
		userID = golfer.ID
	}

	db := session.Database(r)

	if userID != 0 {
		rows, err := db.Query(
			` SELECT hole, lang
			    FROM solutions
			   WHERE failing AND user_id = $1
			ORDER BY hole`,
			userID,
		)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		for rows.Next() {
			var holeID, langID string

			if err := rows.Scan(&holeID, &langID); err != nil {
				panic(err)
			}

			data.Fails = append(data.Fails, Fail{hole.ByID[holeID], lang.ByID[langID]})
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	rows, err := db.Query(
		`WITH ranks AS (
		    SELECT hole, RANK() OVER (PARTITION BY hole ORDER BY chars), user_id
		      FROM solutions
		      JOIN code ON code_id = id
		     WHERE NOT failing
		) SELECT COUNT(*),
		         (SELECT COALESCE(MIN(rank), 0) FROM ranks WHERE hole = r.hole AND user_id = $1),
		         hole
		    FROM ranks r
		GROUP BY hole`,
		userID,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var card Card
		var holeID string

		if err := rows.Scan(&card.Golfers, &card.Rank, &holeID); err != nil {
			panic(err)
		}

		card.Hole = hole.ByID[holeID]

		data.Cards = append(data.Cards, card)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	switch data.Sort {
	case "alphabetical-desc":
		sort.Slice(data.Cards, func(i, j int) bool {
			return data.Cards[i].Hole.Name > data.Cards[j].Hole.Name
		})
	case "golfers-asc":
		sort.Slice(data.Cards, func(i, j int) bool {
			if data.Cards[i].Golfers == data.Cards[j].Golfers {
				return data.Cards[i].Hole.Name < data.Cards[j].Hole.Name
			}
			return data.Cards[i].Golfers < data.Cards[j].Golfers
		})
	case "golfers-desc":
		sort.Slice(data.Cards, func(i, j int) bool {
			if data.Cards[i].Golfers == data.Cards[j].Golfers {
				return data.Cards[i].Hole.Name < data.Cards[j].Hole.Name
			}
			return data.Cards[i].Golfers > data.Cards[j].Golfers
		})
	case "rank-asc":
		sort.Slice(data.Cards, func(i, j int) bool {
			if data.Cards[i].Rank == data.Cards[j].Rank {
				return data.Cards[i].Hole.Name < data.Cards[j].Hole.Name
			}
			return data.Cards[i].Rank < data.Cards[j].Rank
		})
	case "rank-desc":
		sort.Slice(data.Cards, func(i, j int) bool {
			if data.Cards[i].Rank == data.Cards[j].Rank {
				return data.Cards[i].Hole.Name < data.Cards[j].Hole.Name
			}
			return data.Cards[i].Rank > data.Cards[j].Rank
		})
	default:
		data.Sort = "alphabetical-asc"

		sort.Slice(data.Cards, func(i, j int) bool {
			return data.Cards[i].Hole.Name < data.Cards[j].Hole.Name
		})
	}

	w.Header().Set(
		"Strict-Transport-Security",
		"max-age=31536000;includeSubDomains;preload",
	)
	render(w, r, "index", "", data)
}
