package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// RankingsMedals serves GET /rankings/medals
func RankingsMedals(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country, Login             string
		Rank, Gold, Silver, Bronze int
	}

	data := struct {
		Holes []hole.Hole
		Langs []lang.Lang
		Rows  []row
	}{
		Holes: hole.List,
		Langs: lang.List,
	}

	rows, err := session.Database(r).Query(
		`SELECT RANK() OVER(ORDER BY gold DESC, silver DESC, bronze DESC),
		        COALESCE(CASE WHEN show_country THEN country END, ''),
		        login,
		        gold,
		        silver,
		        bronze
		   FROM medals
		   JOIN users ON id = user_id
		  LIMIT 25`,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var r row

		if err := rows.Scan(
			&r.Rank,
			&r.Country,
			&r.Login,
			&r.Gold,
			&r.Silver,
			&r.Bronze,
		); err != nil {
			panic(err)
		}

		data.Rows = append(data.Rows, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "rankings/medals", "Rankings: Medals", data)
}
