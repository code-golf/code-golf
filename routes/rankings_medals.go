package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/medals/{hole}/{lang}
func rankingsMedalsGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country, Login                      string
		Rank, Diamond, Gold, Silver, Bronze int
	}

	data := struct {
		HoleID, LangID, Scoring string
		Holes                   []*config.Hole
		Langs                   []*config.Lang
		Pager                   *pager.Pager
		Rows                    []row
	}{
		HoleID:  param(r, "hole"),
		Holes:   config.HoleList,
		LangID:  param(r, "lang"),
		Langs:   config.LangList,
		Pager:   pager.New(r),
		Rows:    make([]row, 0, pager.PerPage),
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "all" && data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rows, err := session.Database(r).Query(
		`WITH counts AS (
		    SELECT user_id,
		           COUNT(*) FILTER (WHERE medal = 'diamond') diamond,
		           COUNT(*) FILTER (WHERE medal = 'gold'   ) gold,
		           COUNT(*) FILTER (WHERE medal = 'silver' ) silver,
		           COUNT(*) FILTER (WHERE medal = 'bronze' ) bronze
		      FROM medals
		     WHERE $1 IN ('all', hole::text)
		       AND $2 IN ('all', lang::text)
		       AND $3 IN ('all', scoring::text)
		  GROUP BY user_id
		) SELECT RANK() OVER(
		             ORDER BY gold DESC, diamond DESC, silver DESC, bronze DESC
		         ),
		         country_flag,
		         login,
		         diamond,
		         gold,
		         silver,
		         bronze,
		         COUNT(*) OVER()
		    FROM counts
		    JOIN users ON id = user_id
		ORDER BY rank, login
		   LIMIT $4 OFFSET $5`,
		data.HoleID,
		data.LangID,
		data.Scoring,
		pager.PerPage,
		data.Pager.Offset,
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
			&r.Diamond,
			&r.Gold,
			&r.Silver,
			&r.Bronze,
			&data.Pager.Total,
		); err != nil {
			panic(err)
		}

		data.Rows = append(data.Rows, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if data.Pager.Calculate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	description := ""
	if hole, ok := config.HoleByID[data.HoleID]; ok {
		description += hole.Name + " in "
	} else {
		description += "All holes in "
	}

	if lang, ok := config.LangByID[data.LangID]; ok {
		description += lang.Name
	} else {
		description += "all languages"
	}

	if data.Scoring != "all" {
		description += " in " + data.Scoring
	}

	description += "."

	render(w, r, "rankings/medals", data, "Rankings: Medals", description)
}
