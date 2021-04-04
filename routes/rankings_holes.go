package routes

import (
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
		Country, Login               string
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

	var distinct, table string

	if data.HoleID == "all" {
		distinct = "DISTINCT ON (hole, user_id)"
		table = "summed_leaderboard"
	} else {
		table = "scored_leaderboard"
	}

	rows, err := session.Database(r).Query(
		`WITH leaderboard AS (
		  SELECT `+distinct+`
		         hole,
		         submitted,
		         `+data.Scoring+` strokes,
		         user_id,
		         lang
		    FROM solutions
		    JOIN code ON code_id = id
		   WHERE NOT failing
		     AND $1 IN ('all', hole::text)
		     AND $2 IN ('all', lang::text)
		     AND scoring = $3
		ORDER BY hole, user_id, `+data.Scoring+`, submitted
		), scored_leaderboard AS (
		  SELECT l.hole,
		         1 holes,
		         ROUND(
		             (COUNT(*) OVER (PARTITION BY l.hole) -
		                RANK() OVER (PARTITION BY l.hole ORDER BY strokes) + 1)
		             * (1000.0 / COUNT(*) OVER (PARTITION BY l.hole))
		         ) points,
		         strokes,
		         submitted,
		         l.user_id
		    FROM leaderboard l
		), summed_leaderboard AS (
		  SELECT user_id,
		         COUNT(*)       holes,
		         SUM(points)    points,
		         SUM(strokes)   strokes,
		         MAX(submitted) submitted
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT COALESCE(CASE WHEN show_country THEN country END, ''),
		         holes,
		         login,
		         points,
		         RANK() OVER (ORDER BY points DESC, strokes),
		         strokes,
		         submitted,
		         COUNT(*) OVER()
		    FROM `+table+`
		    JOIN users on user_id = id
		ORDER BY points DESC, strokes, submitted
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
			&r.Country,
			&r.Holes,
			&r.Login,
			&r.Points,
			&r.Rank,
			&r.Strokes,
			&r.Submitted,
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

	description += data.Scoring

	render(w, r, "rankings/holes", data, "Rankings: Holes", description)
}
