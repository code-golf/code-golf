package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/medals/{hole}/{lang}
func rankingsMedalsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Hole, PrevHole, NextHole *config.Hole
		HoleID, LangID, Scoring  string
		Pager                    *pager.Pager
		Rows                     []struct {
			Country                                *config.Country
			Unicorn, Diamond, Gold, Silver, Bronze int
			Name                                   string
			Rank, Total                            int
		}
	}{
		HoleID:  param(r, "hole"),
		LangID:  param(r, "lang"),
		Pager:   pager.New(r),
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "all" && data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.Hole = config.HoleByID[data.HoleID]; data.Hole != nil {
		data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)
	}

	if err := session.Database(r).Select(
		&data.Rows,
		`WITH counts AS (
		    SELECT user_id,
		           COUNT(*) FILTER (WHERE medal = 'unicorn') unicorn,
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
		         unicorn, diamond, gold, silver, bronze,
		         country_flag country, login name, COUNT(*) OVER() total
		    FROM counts
		    JOIN users ON id = user_id
		ORDER BY rank, login
		   LIMIT $4 OFFSET $5`,
		data.HoleID,
		data.LangID,
		data.Scoring,
		pager.PerPage,
		data.Pager.Offset,
	); err != nil {
		panic(err)
	}

	if len(data.Rows) > 0 {
		data.Pager.Total = data.Rows[0].Total
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

	description += ". Ranked by golds, then diamonds, then silvers, then bronzes."

	render(w, r, "rankings/medals", data, "Rankings: Medals", description)
}
