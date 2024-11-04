package routes

import (
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
	data := struct {
		Hole, PrevHole, NextHole              *config.Hole
		HoleID, LangID, OtherScoring, Scoring string
		Holes                                 []*config.Hole
		Langs                                 []*config.Lang
		Pager                                 *pager.Pager
		Recent                                bool
		Rows                                  []struct {
			Country                             config.NullCountry
			Holes, Rank, Points, Strokes, Total int
			Lang                                *config.Lang
			Name                                string
			OtherStrokes                        *int
			Submitted                           time.Time
		}
	}{
		HoleID:  param(r, "hole"),
		Holes:   config.HoleList,
		LangID:  param(r, "lang"),
		Langs:   config.LangList,
		Pager:   pager.New(r),
		Recent:  strings.HasPrefix(r.URL.Path, "/rankings/recent-holes"),
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

	var sql string
	var bind []any

	// TODO Try and merge these SQL queries?
	if data.HoleID == "all" && data.LangID == "all" {
		sql = `WITH foo AS (
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
			) SELECT country_flag                                 country,
			         holes                                        holes,
			         login                                        name,
			         points                                       points,
			         RANK() OVER (ORDER BY points DESC, strokes)  rank,
			         strokes                                      strokes,
			         submitted                                    submitted,
			         COUNT(*) OVER()                              total
			    FROM summed
			    JOIN users ON user_id = id
			ORDER BY rank, submitted
			   LIMIT $2 OFFSET $3`

		bind = []any{data.Scoring, pager.PerPage, data.Pager.Offset, holeWhere}
	} else if data.HoleID == "all" {
		sql = `WITH summed AS (
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
			) SELECT country_flag                                 country,
			         holes                                        holes,
			         login                                        name,
			         points                                       points,
			         RANK() OVER (ORDER BY points DESC, strokes)  rank,
			         strokes                                      strokes,
			         submitted                                    submitted,
			         COUNT(*) OVER()                              total
			    FROM summed
			    JOIN users ON user_id = id
			ORDER BY rank, submitted
			   LIMIT $3 OFFSET $4`

		bind = []any{data.LangID, data.Scoring, pager.PerPage, data.Pager.Offset, holeWhere}
	} else if data.LangID == "all" {
		sql = `SELECT country_flag     country,
			          lang             lang,
			          login            name,
			          other_strokes    other_strokes,
			          points           points,
			          rank_overall     rank,
			          strokes          strokes,
			          submitted        submitted,
			          COUNT(*) OVER()  total
			     FROM rankings
			     JOIN users ON user_id = id
			    WHERE hole = $1 AND scoring = $2
			 ORDER BY rank_overall, submitted
			    LIMIT $3 OFFSET $4`

		bind = []any{data.HoleID, data.Scoring, pager.PerPage, data.Pager.Offset}
	} else {
		sql = `SELECT country_flag     country,
			          lang             lang,
			          login            name,
			          other_strokes    other_strokes,
			          points_for_lang  points,
			          rank             rank,
			          strokes          strokes,
			          submitted        submitted,
			          COUNT(*) OVER()  total
			     FROM rankings
			     JOIN users ON user_id = id
			    WHERE hole = $1 AND lang = $2 AND scoring = $3
			 ORDER BY rank, submitted
			    LIMIT $4 OFFSET $5`

		bind = []any{data.HoleID, data.LangID, data.Scoring, pager.PerPage, data.Pager.Offset}
	}

	if err := session.Database(r).Select(&data.Rows, sql, bind...); err != nil {
		panic(err)
	}

	if len(data.Rows) > 0 {
		data.Pager.Total = data.Rows[0].Total
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
