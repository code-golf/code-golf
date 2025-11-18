package routes

import (
	"html/template"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/cheevos/{cheevo}
func rankingsCheevosGET(w http.ResponseWriter, r *http.Request) {
	cheevoID := param(r, "cheevo")

	data := struct {
		Cheevo *config.Cheevo
		Pager  *pager.Pager
		Rows   []struct {
			Country            *config.Country
			Earned             time.Time
			Name               string
			Count, Rank, Total int
		}
		Total int
	}{
		Cheevo: config.CheevoByID[cheevoID],
		Pager:  pager.New(r),
		Total:  len(config.CheevoList),
	}

	if cheevoID != "all" && data.Cheevo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := session.Database(r).Select(
		&data.Rows,
		`WITH count AS (
		    SELECT user_id, COUNT(*), MAX(earned) earned
		      FROM cheevos
		     WHERE cheevo = $1 OR $1 IS NULL
		  GROUP BY user_id
		) SELECT count, country_flag country, earned, login name,
		         CASE WHEN $1 IS NULL
		            THEN RANK() OVER(ORDER BY count DESC)
		            ELSE RANK() OVER(ORDER BY earned)
		         END,
		         COUNT(*) OVER() total
		    FROM count JOIN users ON id = user_id
		ORDER BY rank, earned, login
		   LIMIT $2 OFFSET $3`,
		data.Cheevo,
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

	description := template.HTML("All achievements")
	if cheevo := data.Cheevo; cheevo != nil {
		description = template.HTML(cheevo.Emoji+" "+cheevo.Name+" - ") + cheevo.Description
	}

	render(w, r, "rankings/cheevos", data, "Rankings: Achievements", description)
}
