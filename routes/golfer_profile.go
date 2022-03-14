package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}
func golferGET(w http.ResponseWriter, r *http.Request) {
	const limit = 100

	type follow struct {
		Bytes, Chars, Rank int
		Country, Name      string
	}

	type row struct {
		Cheevo *config.Cheevo
		Date   time.Time
		Golfer string
		Hole   *config.Hole
		Lang   *config.Lang
	}

	data := struct {
		Follows  []follow
		Langs    []*config.Lang
		Trophies map[string]map[string]int
		Wall     []row
	}{
		Langs:    config.LangList,
		Trophies: map[string]map[string]int{},
		Wall:     make([]row, 0, limit),
	}

	db := session.Database(r)
	golfer := session.GolferInfo(r).Golfer

	// TODO Support friends/follow.
	rows, err := db.Query(
		`WITH data AS (
		 -- Cheevos
		    SELECT earned       date,
		           trophy::text cheevo,
		           ''           hole,
		           ''           lang,
		           user_id
		      FROM trophies
		     WHERE user_id = ANY(following($1))
		 UNION ALL
		 -- Holes
		    SELECT MAX(submitted) date,
		           ''             cheevo,
		           hole::text     hole,
		           lang::text     lang,
		           user_id
		      FROM solutions
		     WHERE user_id = ANY(following($1))
		  GROUP BY user_id, hole, lang
		) SELECT cheevo, date, login, hole, lang
		    FROM data JOIN users ON id = user_id
		ORDER BY date DESC LIMIT $2`,
		golfer.ID,
		limit,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cheevo, golfer, hole, lang string
		var date time.Time

		if err := rows.Scan(&cheevo, &date, &golfer, &hole, &lang); err != nil {
			panic(err)
		}

		// TODO Parse date into viewers location.
		data.Wall = append(data.Wall, row{
			Cheevo: config.CheevoByID[cheevo],
			Date:   date,
			Golfer: golfer,
			Hole:   config.HoleByID[hole],
			Lang:   config.LangByID[lang],
		})
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = db.Query(
		`WITH summed AS (
		    SELECT user_id, scoring, lang, SUM(points_for_lang)
		      FROM rankings
		  GROUP BY user_id, scoring, lang
		), ranks AS (
		    SELECT user_id, scoring, lang,
		           RANK() OVER (PARTITION BY scoring, lang ORDER BY sum DESC)
		      FROM summed
		) SELECT lang, scoring, rank
		    FROM ranks
		   WHERE rank < 4 AND user_id = $1`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var lang, scoring string
		var rank int

		if err := rows.Scan(&lang, &scoring, &rank); err != nil {
			panic(err)
		}

		if _, ok := data.Trophies[lang]; !ok {
			data.Trophies[lang] = map[string]int{}
		}

		data.Trophies[lang][scoring] = rank
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = db.Query(
		`WITH follows AS (
		    SELECT country_flag, login,
		           COALESCE((SELECT points FROM points
		              WHERE scoring = 'bytes' AND user_id = id), 0) bytes,
		           COALESCE((SELECT points FROM points
		              WHERE scoring = 'chars' AND user_id = id), 0) chars
		      FROM users
		     WHERE id = ANY(following($1))
		) SELECT *, RANK() OVER (ORDER BY bytes DESC, chars DESC)
		    FROM follows
		ORDER BY rank, login`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var f follow

		if err := rows.Scan(
			&f.Country,
			&f.Name,
			&f.Bytes,
			&f.Chars,
			&f.Rank,
		); err != nil {
			panic(err)
		}

		data.Follows = append(data.Follows, f)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/profile", data, golfer.Name)
}
