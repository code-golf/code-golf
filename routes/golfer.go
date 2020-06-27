package routes

import (
	"database/sql"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
)

// Golfer serves GET /golfers/{golfer}
func Golfer(w http.ResponseWriter, r *http.Request) {
	type EarnedTrophy struct {
		Count, Percent int
		Earned         sql.NullTime
		Trophy         Trophy
	}

	data := struct {
		Max      int
		Trophies []EarnedTrophy
	}{
		Trophies: make([]EarnedTrophy, 0, len(trophies)),
	}

	tx, err := db(r).BeginTx(
		r.Context(),
		&sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true},
	)
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	if err := tx.QueryRow(
		"SELECT COUNT(DISTINCT user_id) FROM trophies",
	).Scan(&data.Max); err != nil {
		panic(err)
	}

	rows, err := tx.Query(
		`WITH count AS (
		    SELECT trophy, COUNT(*) FROM trophies GROUP BY trophy
		), earned AS (
		    SELECT trophy, earned FROM trophies WHERE user_id = $1
		)  SELECT *
		     FROM count
		LEFT JOIN earned USING(trophy)
		 ORDER BY count DESC, trophy`,
		r.Context().Value("golferInfo").(*golfer.GolferInfo).ID,
	)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var trophy EarnedTrophy

		if err := rows.Scan(
			&trophy.Trophy.ID,
			&trophy.Count,
			&trophy.Earned,
		); err != nil {
			panic(err)
		}

		trophy.Percent = trophy.Count * 100 / data.Max
		trophy.Trophy = trophiesByID[trophy.Trophy.ID]

		data.Trophies = append(data.Trophies, trophy)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	render(w, r, http.StatusOK, "golfer", "JRaspass", data)
}

// GolferHoles serves GET /golfers/{golfer}/holes
func GolferHoles(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "golfer-holes", "JRaspass", langs)
}
