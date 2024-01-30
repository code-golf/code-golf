package golfer

import (
	"database/sql"
	"errors"

	"github.com/code-golf/code-golf/config"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

// Get a Golfer given a session ID, updates the session's last used time.
func Get(db *sqlx.DB, sessionID uuid.UUID) *Golfer {
	var golfer Golfer

	if err := db.Get(
		&golfer,
		`WITH golfer AS (
		    UPDATE sessions SET last_used = DEFAULT WHERE id = $1
		    RETURNING user_id
		), failing AS (
		    SELECT hole, lang
		      FROM solutions
		      JOIN golfer USING(user_id)
		     WHERE failing
		  GROUP BY hole, lang
		  ORDER BY hole, lang
		)  SELECT u.admin                                   admin,
		          COALESCE(bytes.points, 0)                 bytes_points,
		          COALESCE(chars.points, 0)                 chars_points,
		          u.country                                 country,
		          u.delete                                  delete,
		          (SELECT COALESCE(json_agg(failing), '[]')
		             FROM failing)                          failing_solutions,
		          u.id                                      id,
		          u.layout                                  layout,
		          u.keymap                                  keymap,
		          u.login                                   name,
		          u.pronouns                                pronouns,
		          COALESCE(r.login, '')                     referrer,
		          u.show_country                            show_country,
		          u.sponsor                                 sponsor,
		          u.theme                                   theme,
		          u.time_zone                               time_zone,
		          ARRAY(SELECT trophy
		                  FROM trophies
		                 WHERE user_id = u.id
		              ORDER BY trophy)                      cheevos,
		          ARRAY(SELECT followee_id
		                  FROM follows
		                 WHERE follower_id = u.id
		              ORDER BY followee_id)                 following,
		          ARRAY(SELECT DISTINCT hole
		                  FROM solutions
		                 WHERE user_id = u.id
		              ORDER BY hole)                        holes
		     FROM users  u
		     JOIN golfer g     ON u.id = g.user_id
		LEFT JOIN users  r     ON r.id = u.referrer_id
		LEFT JOIN points bytes ON u.id = bytes.user_id AND bytes.scoring = 'bytes'
		LEFT JOIN points chars ON u.id = chars.user_id AND chars.scoring = 'chars'`,
		sessionID,
	); errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		panic(err)
	}

	return &golfer
}

func GetInfo(db *sqlx.DB, name string) *GolferInfo {
	info := GolferInfo{
		CheevosTotal: len(config.CheevoList),
		HolesTotal:   len(config.HoleList),
		LangsTotal:   len(config.LangList),
	}

	if err := db.Get(
		&info,
		`WITH medals AS (
		   SELECT user_id,
		          COUNT(*) FILTER (WHERE medal = 'diamond') diamond,
		          COUNT(*) FILTER (WHERE medal = 'gold'   ) gold,
		          COUNT(*) FILTER (WHERE medal = 'silver' ) silver,
		          COUNT(*) FILTER (WHERE medal = 'bronze' ) bronze
		     FROM medals
		 GROUP BY user_id
		)  SELECT admin,
		          COALESCE(bronze, 0)                   bronze,
		          ARRAY(
		            SELECT trophy
		              FROM trophies
		             WHERE user_id = users.id
		          ORDER BY trophy
		          )                                     cheevos,
		          country_flag                          country,
		          COALESCE(diamond, 0)                  diamond,
		          COALESCE(gold, 0)                     gold,
		          (SELECT COUNT(DISTINCT hole)
		             FROM solutions
		            WHERE user_id = id AND NOT FAILING) holes,
		          id,
		          (SELECT COUNT(DISTINCT lang)
		             FROM solutions
		            WHERE user_id = id AND NOT FAILING) langs,
		          login                                 name,
		          COALESCE(bytes.points, 0)             bytes_points,
		          COALESCE(chars.points, 0)             chars_points,
		          pronouns                              pronouns,
		          ARRAY(
		            SELECT login
		              FROM users u
		             WHERE referrer_id = users.id
		          ORDER BY login
		          )                                     referrals,
		          COALESCE(silver, 0)                   silver,
		          sponsor,
		          started
		     FROM users
		LEFT JOIN medals       ON id = medals.user_id
		LEFT JOIN points bytes ON id = bytes.user_id AND bytes.scoring = 'bytes'
		LEFT JOIN points chars ON id = chars.user_id AND chars.scoring = 'chars'
		    WHERE login = $1`,
		name,
	); errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		panic(err)
	}

	return &info
}
