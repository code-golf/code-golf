package routes

import (
	"database/sql"
	"net/http"
	"net/url"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// RankingsHoles serves GET /rankings/holes/{hole}/{lang}/{scoring}
func RankingsHoles(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country, Login               string
		Holes, Rank, Points, Strokes int
		OtherStrokes                 *int
		Lang                         *config.Lang
		Submitted                    time.Time
	}

	data := struct {
		HoleID, LangID, OtherScoring, Query, Scoring, Sorting string
		HolesSorting, PointsSorting, StrokesSorting           *url.URL
		Holes                                                 []*config.Hole
		Langs                                                 []*config.Lang
		Pager                                                 *pager.Pager
		Rows                                                  []row
	}{
		HoleID:  param(r, "hole"),
		Holes:   config.HoleList,
		LangID:  param(r, "lang"),
		Langs:   config.LangList,
		Pager:   pager.New(r),
		Rows:    make([]row, 0, pager.PerPage),
		Scoring: param(r, "scoring"),
		Sorting: r.FormValue("sorting"),
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		NotFound(w, r)
		return
	}

	if data.Scoring == "bytes" {
		data.OtherScoring = "chars"
	} else {
		data.OtherScoring = "bytes"
	}

	if data.Sorting != "holes" && data.Sorting != "points" && data.Sorting != "strokes" ||
		data.HoleID != "all" && data.Sorting == "holes" ||
		data.HoleID != "all" && data.LangID != "all" && data.Sorting != "points" {
		data.Sorting = "points"
	}

	if data.Sorting != "points" {
		// This is used when navigating between holes.
		data.Query = "?sorting=" + data.Sorting
	}

	if data.HoleID == "all" && data.Sorting != "holes" {
		data.HolesSorting = changeSorting(r.URL, "holes")
	}

	if data.HoleID == "all" || data.LangID == "all" {
		if data.Sorting != "points" {
			data.PointsSorting = changeSorting(r.URL, "points")
		}
		if data.Sorting != "strokes" {
			data.StrokesSorting = changeSorting(r.URL, "strokes")
		}
	}

	var rows *sql.Rows
	var err error
	var order string

	if data.Sorting == "holes" {
		order = "ORDER BY holes DESC, points DESC, strokes"
	} else if data.Sorting == "points" {
		order = "ORDER BY points DESC, strokes"
	} else if data.Sorting == "strokes" && data.HoleID == "all" {
		// Sorting by strokes doesn't seem to useful here. Sort by holes first.
		order = "ORDER BY holes DESC, strokes, points DESC"
	} else {
		order = "ORDER BY strokes, points DESC"
	}

	// TODO Try and merge these SQL queries?
	if data.HoleID == "all" && data.LangID == "all" {
		rows, err = session.Database(r).Query(
			`WITH foo AS (
			    SELECT user_id, hole, lang, points, strokes, submitted,
			           ROW_NUMBER() OVER (
			               PARTITION BY user_id, hole ORDER BY points DESC, strokes
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
			         '',
			         login,
			         points,
			         RANK() OVER (`+order+`),
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
			         RANK() OVER (`+order+`),
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
		)
	} else if data.LangID == "all" {
		rows, err = session.Database(r).Query(
			` SELECT COALESCE(CASE WHEN show_country THEN country END, ''),
			         1,
			         lang,
			         login,
			         points,
			         RANK() OVER (`+order+`),
			         strokes,
			         other_strokes,
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
		var langID string

		if err := rows.Scan(
			&r.Country,
			&r.Holes,
			&langID,
			&r.Login,
			&r.Points,
			&r.Rank,
			&r.Strokes,
			&r.OtherStrokes,
			&r.Submitted,
			&data.Pager.Total,
		); err != nil {
			panic(err)
		}

		r.Lang = config.LangByID[langID]

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
	if hole, ok := config.HoleByID[data.HoleID]; ok {
		description += hole.Name + " in "
	} else {
		description += "All holes in "
	}

	if lang, ok := config.LangByID[data.LangID]; ok {
		description += lang.Name + " in "
	} else {
		description += "all languages in "
	}

	description += data.Scoring + "."

	render(w, r, "rankings/holes", data, "Rankings: Holes", description)
}

func changeSorting(u *url.URL, sorting string) *url.URL {
	q := u.Query()
	q.Del("page")
	if sorting == "points" {
		q.Del("sorting")
	} else {
		q.Set("sorting", sorting)
	}

	new := *u
	new.RawQuery = q.Encode()

	return &new
}
