package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// RankingsHoles serves GET /rankings/holes/{hole}/{lang}/{scoring}
func RankingsHoles(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country, Lang, Login         string
		Holes, Rank, Points, Strokes int
		Submitted                    time.Time
	}

	data := struct {
		HoleID, LangID, Scoring string
		Holes                   []hole.Hole
		Langs                   []lang.Lang
		Pager                   *pager.Pager
		Rows                    []row
	}{
		HoleID:  param(r, "hole"),
		Holes:   hole.List,
		LangID:  param(r, "lang"),
		Langs:   lang.List,
		Pager:   pager.New(r),
		Rows:    make([]row, 0, pager.PerPage),
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && hole.ByID[data.HoleID].ID == "" ||
		data.LangID != "all" && lang.ByID[data.LangID].ID == "" ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		NotFound(w, r)
		return
	}

	var rows *sql.Rows
	var err error

	// TODO Try and merge these SQL queries?
	if data.HoleID == "all" && data.LangID == "all" {
		rows, err = session.Database(r).Query(
			`WITH foo AS (
			    SELECT user_id, hole, lang, points, strokes, submitted,
			           ROW_NUMBER() OVER (
			               PARTITION BY user_id, hole ORDER BY points DESC
			           )
			      FROM rankings
			     WHERE scoring = $1
			), summed AS (
			    SELECT user_id,
			           COUNT(*)       holes,
			           SUM(points)    points,
			           SUM(strokes)   strokes,
			           MAX(submitted) submitted
			      FROM foo
			     WHERE row_number = 1
			  GROUP BY user_id
			) SELECT COALESCE(CASE WHEN show_country THEN country END, ''),
			         holes,
			         $1,
			         login,
			         points,
			         RANK() OVER (ORDER BY points DESC, strokes),
			         strokes,
			         submitted,
			         COUNT(*) OVER()
			    FROM summed
			    JOIN users ON user_id = id
			ORDER BY rank, submitted
			   LIMIT $2 OFFSET $3`,
			data.Scoring,
			pager.PerPage,
			data.Pager.Offset,
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
			     WHERE lang    = $1
			       AND scoring = $2
			  GROUP BY user_id
			) SELECT COALESCE(CASE WHEN show_country THEN country END, ''),
			         holes,
			         $1,
			         login,
			         points,
			         RANK() OVER (ORDER BY points DESC, strokes),
			         strokes,
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
		)
	} else if data.LangID == "all" {
		rows, err = session.Database(r).Query(
			` SELECT COALESCE(CASE WHEN show_country THEN country END, ''),
			         1,
			         lang,
			         login,
			         points,
			         RANK() OVER (ORDER BY points DESC, strokes),
			         strokes,
			         submitted,
			         COUNT(*) OVER()
			    FROM rankings
			    JOIN users ON user_id = id
			   WHERE hole = $1 AND scoring = $2
			ORDER BY rank, submitted
			   LIMIT $3 OFFSET $4`,
			data.HoleID,
			data.Scoring,
			pager.PerPage,
			data.Pager.Offset,
		)
	} else {
		rows, err = session.Database(r).Query(
			` SELECT COALESCE(CASE WHEN show_country THEN country END, ''),
			         1,
			         lang,
			         login,
			         points_for_lang,
			         RANK() OVER (ORDER BY points_for_lang DESC, strokes),
			         strokes,
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

		if err := rows.Scan(
			&r.Country,
			&r.Holes,
			&r.Lang,
			&r.Login,
			&r.Points,
			&r.Rank,
			&r.Strokes,
			&r.Submitted,
			&data.Pager.Total,
		); err != nil {
			panic(err)
		}

		r.Lang = lang.ByID[r.Lang].Name

		data.Rows = append(data.Rows, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if data.Pager.Calculate() {
		NotFound(w, r)
		return
	}

	description := ""
	if hole := hole.ByID[data.HoleID]; hole.ID != "" {
		description += hole.Name + " in "
	} else {
		description += "All holes in "
	}

	if lang := lang.ByID[data.LangID]; lang.ID != "" {
		description += lang.Name + " in "
	} else {
		description += "all languages in "
	}

	description += data.Scoring + "."

	render(w, r, "rankings/holes", data, "Rankings: Holes", description)
}
