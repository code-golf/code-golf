package routes

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/holes/{hole}/{lang}/{scoring}
func rankingsHolesGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Distribution                          []struct{ Frequency, Strokes int }
		Strokes                               int
		Hole, PrevHole, NextHole              *config.Hole
		HoleID, LangID, OtherScoring, Scoring string
		Pager                                 *pager.Pager
		Rows                                  []struct {
			AvatarURL, Name                          string
			Country                                  *config.Country
			Holes, Points, Rank, Row, Strokes, Total int
			Lang                                     *config.Lang
			OtherStrokes                             *int
			Submitted                                time.Time
			Time                                     time.Duration
		}
	}{
		HoleID:  param(r, "hole"),
		LangID:  param(r, "lang"),
		Pager:   pager.New(r),
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && config.AllHoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.AllLangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.Hole = config.AllHoleByID[data.HoleID]; data.Hole != nil {
		data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)
	}

	if data.Scoring == "bytes" {
		data.OtherScoring = "chars"
	} else {
		data.OtherScoring = "bytes"
	}

	var query string
	var bind []any

	if data.HoleID == "all" && data.LangID == "all" {
		query = `
			  SELECT avatar_url                                  avatar_url,
			         country_flag                                country,
			         holes                                       holes,
			         name                                        name,
			         points                                      points,
			         RANK() OVER (ORDER BY points DESC, strokes) rank,
			         strokes                                     strokes,
			         submitted                                   submitted,
			         COUNT(*) OVER()                             total
			    FROM points
			    JOIN golfers_with_avatars ON user_id = id
			   WHERE scoring = $1
			ORDER BY rank, submitted
			   LIMIT $2 OFFSET $3`

		bind = []any{data.Scoring, pager.PerPage, data.Pager.Offset}
	} else if data.HoleID == "all" {
		query = `WITH summed AS (
			    SELECT user_id,
			           COUNT(*)             holes,
			           SUM(points_for_lang) points,
			           SUM(strokes)         strokes,
			           MAX(submitted)       submitted
			      FROM rankings
			     WHERE lang    = $1
			       AND scoring = $2
			       AND NOT experimental_hole
			  GROUP BY user_id
			) SELECT avatar_url                                   avatar_url,
			         country_flag                                 country,
			         holes                                        holes,
			         name                                         name,
			         points                                       points,
			         RANK() OVER (ORDER BY points DESC, strokes)  rank,
			         strokes                                      strokes,
			         submitted                                    submitted,
			         COUNT(*) OVER()                              total
			    FROM summed
			    JOIN golfers_with_avatars ON user_id = id
			ORDER BY rank, submitted
			   LIMIT $3 OFFSET $4`

		bind = []any{data.LangID, data.Scoring, pager.PerPage, data.Pager.Offset}
	} else if data.LangID == "all" {
		query = `SELECT avatar_url       avatar_url,
			          country_flag     country,
			          lang             lang,
			          name             name,
			          other_strokes    other_strokes,
			          points           points,
			          rank_overall     rank,
			          row_overall      row,
			          strokes          strokes,
			          submitted        submitted,
			          time_ms * 1e6    time,
			          COUNT(*) OVER()  total
			     FROM rankings
			     JOIN golfers_with_avatars ON user_id = id
			    WHERE hole = $1 AND scoring = $2
			 ORDER BY rank_overall, submitted
			    LIMIT $3 OFFSET $4`

		bind = []any{data.HoleID, data.Scoring, pager.PerPage, data.Pager.Offset}
	} else {
		query = `SELECT avatar_url       avatar_url,
			          country_flag     country,
			          lang             lang,
			          name             name,
			          other_strokes    other_strokes,
			          points_for_lang  points,
			          rank             rank,
			          row              row,
			          strokes          strokes,
			          submitted        submitted,
			          time_ms * 1e6    time,
			          COUNT(*) OVER()  total
			     FROM rankings
			     JOIN golfers_with_avatars ON user_id = id
			    WHERE hole = $1 AND lang = $2 AND scoring = $3
			 ORDER BY rank, submitted
			    LIMIT $4 OFFSET $5`

		bind = []any{data.HoleID, data.LangID, data.Scoring, pager.PerPage, data.Pager.Offset}
	}

	if err := session.Database(r).Select(&data.Rows, query, bind...); err != nil {
		panic(err)
	}

	if len(data.Rows) > 0 {
		data.Pager.Total = data.Rows[0].Total
	}

	if data.LangID != "all" && data.HoleID != "all" {
		query = `SELECT strokes  strokes,
					  COUNT(*) frequency
				 FROM rankings
				WHERE hole = $1 AND lang = $2 AND scoring = $3
			 GROUP BY strokes
			 ORDER BY strokes`
		if err := session.Database(r).Select(&data.Distribution, query, data.HoleID, data.LangID, data.Scoring); err != nil {
			panic(err)
		}
		golfer := session.Golfer(r)
		if golfer != nil {
			if err := session.Database(r).Get(
				&data,
				`SELECT strokes  strokes
				FROM rankings
				WHERE hole = $1 AND lang = $2 AND scoring = $3 AND user_id = $4`,
				data.HoleID,
				data.LangID,
				data.Scoring,
				golfer.ID,
			); err != nil && !errors.Is(err, sql.ErrNoRows) {
				panic(err)
			}
		}
	}

	if data.Pager.Calculate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var desc strings.Builder
	if data.Hole != nil {
		desc.WriteString(data.Hole.Name)
		desc.WriteString(" in ")
	} else {
		desc.WriteString("All holes in ")
	}

	if lang, ok := config.AllLangByID[data.LangID]; ok {
		desc.WriteString(lang.Name)
		desc.WriteString(" in ")
	} else {
		desc.WriteString("all languages in ")
	}

	desc.WriteString(data.Scoring)
	desc.WriteByte('.')

	render(w, r, "rankings/holes", data, "Rankings: Holes", template.HTML(desc.String()))
}
