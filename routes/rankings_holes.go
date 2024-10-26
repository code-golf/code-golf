package routes

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /rankings/holes/{hole}/{lang}/{scoring}
// GET /rankings/recent-holes/{lang}/{scoring}
func rankingsHolesGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country                      config.NullCountry
		Holes, Rank, Points, Strokes int
		Lang                         *config.Lang
		Name                         string
		OtherStrokes                 *int
		Submitted                    time.Time
	}

	data := struct {
		Hole, PrevHole, NextHole              *config.Hole
		HoleID, LangID, OtherScoring, Scoring string
		Holes                                 []*config.Hole
		Langs                                 []*config.Lang
		Pager                                 *pager.Pager
		Recent                                bool
		Rows                                  []row
	}{
		HoleID:  param(r, "hole"),
		Holes:   config.HoleList,
		LangID:  param(r, "lang"),
		Langs:   config.LangList,
		Pager:   pager.New(r),
		Recent:  strings.HasPrefix(r.URL.Path, "/rankings/recent-holes"),
		Rows:    make([]row, 0, pager.PerPage),
		Scoring: param(r, "scoring"),
	}

	var holeWhere any
	if data.Recent {
		data.HoleID = "all"
		holeWhere = pq.Array(config.RecentHoles)
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.Hole = config.HoleByID[data.HoleID]; data.Hole != nil {
		data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)
	}

	if data.Scoring == "bytes" {
		data.OtherScoring = "chars"
	} else {
		data.OtherScoring = "bytes"
	}

	var rows *sql.Rows
	var err error

	// TODO Try and merge these SQL queries?
	if data.HoleID == "all" && data.LangID == "all" {
		rows, err = session.Database(r).Query(
			`WITH foo AS (
			    SELECT user_id, hole, lang, points, strokes, submitted,
			           ROW_NUMBER() OVER (
			               PARTITION BY user_id, hole ORDER BY points DESC, strokes
			           )
			      FROM rankings
			     WHERE scoring = $1 AND (hole = ANY($4) OR $4 IS NULL)
			), summed AS (
			    SELECT user_id,
			           COUNT(*)       holes,
			           SUM(points)    points,
			           SUM(strokes)   strokes,
			           MAX(submitted) submitted
			      FROM foo
			     WHERE row_number = 1
			  GROUP BY user_id
			) SELECT country_flag,
			         holes,
			         NULL,
			         login,
			         points,
			         RANK() OVER (ORDER BY points DESC, strokes),
			         strokes,
			         NULL other_strokes,
			         submitted,
			         COUNT(*) OVER()
			    FROM summed
			    JOIN users ON user_id = id
			ORDER BY rank, submitted
			   LIMIT $2 OFFSET $3`,
			data.Scoring,
			pager.PerPage,
			data.Pager.Offset,
			holeWhere,
		)
	} else if data.HoleID == "all" {
		rows, err = session.Database(r).Query(
			`WITH summed AS (
			    SELECT user_id,
			           COUNT(*)             holes,
			           SUM(points_for_lang) points,
			           SUM(strokes)         strokes,
			           MAX(submitted)       submitted
			      FROM rankings
			     WHERE (hole = ANY($5) OR $5 IS NULL)
			       AND lang    = $1
			       AND scoring = $2
			  GROUP BY user_id
			) SELECT country_flag,
			         holes,
			         $1,
			         login,
			         points,
			         RANK() OVER (ORDER BY points DESC, strokes),
			         strokes,
			         NULL other_strokes,
			         submitted,
			         COUNT(*) OVER()
			    FROM summed
			    JOIN users ON user_id = id
			ORDER BY rank, submitted
			   LIMIT $3 OFFSET $4`,
			data.LangID,
			data.Scoring,
			pager.PerPage,
			data.Pager.Offset,
			holeWhere,
		)
	} else if data.LangID == "all" {
		rows, err = session.Database(r).Query(
			` SELECT country_flag,
			         1,
			         lang,
			         login,
			         points,
			         rank_overall,
			         strokes,
			         other_strokes,
			         submitted,
			         COUNT(*) OVER()
			    FROM rankings
			    JOIN users ON user_id = id
			   WHERE hole = $1 AND scoring = $2
			ORDER BY rank_overall, submitted
			   LIMIT $3 OFFSET $4`,
			data.HoleID,
			data.Scoring,
			pager.PerPage,
			data.Pager.Offset,
		)
	} else {
		rows, err = session.Database(r).Query(
			` SELECT country_flag,
			         1,
			         lang,
			         login,
			         points_for_lang,
			         rank,
			         strokes,
			         other_strokes,
			         submitted,
			         COUNT(*) OVER()
			    FROM rankings
			    JOIN users ON user_id = id
			   WHERE hole = $1 AND lang = $2 AND scoring = $3
			ORDER BY rank, submitted
			   LIMIT $4 OFFSET $5`,
			data.HoleID,
			data.LangID,
			data.Scoring,
			pager.PerPage,
			data.Pager.Offset,
		)
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r row
		var lang config.NullLang

		if err := rows.Scan(
			&r.Country,
			&r.Holes,
			&lang,
			&r.Name,
			&r.Points,
			&r.Rank,
			&r.Strokes,
			&r.OtherStrokes,
			&r.Submitted,
			&data.Pager.Total,
		); err != nil {
			panic(err)
		}

		r.Lang = lang.Lang
		data.Rows = append(data.Rows, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if data.Pager.Calculate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var desc strings.Builder
	if hole, ok := config.HoleByID[data.HoleID]; ok {
		desc.WriteString(hole.Name)
		desc.WriteString(" in ")
	} else if data.Recent {
		desc.WriteString("Ten most recent holes in ")
	} else {
		desc.WriteString("All holes in ")
	}

	if lang, ok := config.LangByID[data.LangID]; ok {
		desc.WriteString(lang.Name)
		desc.WriteString(" in ")
	} else {
		desc.WriteString("all languages in ")
	}

	desc.WriteString(data.Scoring)

	if data.Recent {
		desc.WriteString(". <p>")

		for i, hole := range config.RecentHoles {
			if i > 0 {
				desc.WriteString(", ")

				if i == len(config.RecentHoles)-1 {
					desc.WriteString("and ")
				}
			}

			desc.WriteString(`<a href="/`)
			desc.WriteString(hole.ID)
			if data.LangID != "all" {
				desc.WriteByte('#')
				desc.WriteString(data.LangID)
			}
			desc.WriteString(`">`)
			desc.WriteString(hole.Name)
			desc.WriteString("</a>")
		}
	}

	desc.WriteByte('.')

	name := "rankings/holes"
	title := "Rankings: Holes"
	if data.Recent {
		name = "rankings/recent-holes"
		title = "Rankings: Recent Holes"
	}

	render(w, r, name, data, title, template.HTML(desc.String()))
}
