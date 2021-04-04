package routes

import (
	"net/http"
	"sort"
	"strings"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
)

// Index serves GET /
func Index(w http.ResponseWriter, r *http.Request) {
	if scoring := r.URL.Query().Get("scoring"); scoring != "" {
		http.SetCookie(w, &http.Cookie{
			HttpOnly: true,
			Name:     "__Host-scoring",
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Secure:   true,
			Value:    scoring,
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

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

	data := struct {
		Cards    []Card
		Scoring  string
		Scorings []string
		Sort     string
		Sorts    []string
	}{
		Scoring:  "bytes",
		Scorings: []string{"bytes", "chars"},
		Sorts:    []string{"alphabetical", "golfers", "rank"},
	}

	if sort, err := r.Cookie("__Host-sort"); err == nil {
		data.Sort = sort.Value
	}

	if scoring, err := r.Cookie("__Host-scoring"); err == nil {
		for _, value := range data.Scorings {
			if scoring.Value == value {
				data.Scoring = value
				break
			}
		}
	}

	var userID int
	if golfer := session.Golfer(r); golfer != nil {
		userID = golfer.ID
	}

	db := session.Database(r)

	rows, err := db.Query(
		`WITH ranks AS (
		    SELECT hole,
		           RANK() OVER (PARTITION BY hole ORDER BY `+data.Scoring+`),
		           user_id
		      FROM solutions
		      JOIN code ON code_id = id
		     WHERE scoring = $2
		       AND NOT failing
		) SELECT COUNT(*),
		         (SELECT COALESCE(MIN(rank), 0) FROM ranks WHERE hole = r.hole AND user_id = $1),
		         hole
		    FROM ranks r
		GROUP BY hole`,
		userID,
		data.Scoring,
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
			return strings.ToLower(data.Cards[i].Hole.Name) >
				strings.ToLower(data.Cards[j].Hole.Name)
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
			return strings.ToLower(data.Cards[i].Hole.Name) <
				strings.ToLower(data.Cards[j].Hole.Name)
		})
	}

	w.Header().Set(
		"Strict-Transport-Security",
		"max-age=31536000;includeSubDomains;preload",
	)
	render(w, r, "index", data)
}
