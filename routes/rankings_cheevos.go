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

	type row struct {
		Country, Login string
		Earned         time.Time
		Rank, Count    int
	}

	data := struct {
		Cheevo *config.Cheevo
		Pager  *pager.Pager
		Rows   []row
		Total  int
	}{
		Cheevo: config.CheevoByID[cheevoID],
		Pager:  pager.New(r),
		Rows:   make([]row, 0, pager.PerPage),
		Total:  len(config.CheevoList),
	}

	if cheevoID != "" && data.Cheevo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rows, err := session.Database(r).Query(
		`WITH count AS (
		    SELECT user_id, COUNT(*), MAX(earned) earned
		      FROM trophies
		     WHERE $1 IN ('', trophy::text)
		  GROUP BY user_id
		) SELECT count, country_flag, earned, login,
		         CASE WHEN $1 = ''
		            THEN RANK() OVER(ORDER BY count DESC)
		            ELSE RANK() OVER(ORDER BY earned)
		         END,
		         COUNT(*) OVER()
		    FROM count JOIN users ON id = user_id
		ORDER BY rank, earned, login
		   LIMIT $2 OFFSET $3`,
		cheevoID,
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
			&r.Count,
			&r.Country,
			&r.Earned,
			&r.Login,
			&r.Rank,
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

	description := interface{}("All achievements")
	if cheevo := data.Cheevo; cheevo != nil {
		description = template.HTML(cheevo.Emoji+" "+cheevo.Name+" - ") + cheevo.Description
	}

	render(w, r, "rankings/cheevos", data, "Rankings: Achievements", description)
}
