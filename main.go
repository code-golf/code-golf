package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/code-golf/code-golf/discord"
	"github.com/code-golf/code-golf/github"
	"github.com/code-golf/code-golf/routes"
	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.Ltime)

	rand.Seed(time.Now().UnixNano())

	db, err := sql.Open("postgres", "")
	if err != nil {
		panic(err)
	}

	// Every 30 seconds.
	go func() {
		for range time.Tick(10 * time.Second) {
			for _, view := range []string{"medals", "rankings", "points"} {
				if _, err := db.Exec(
					"REFRESH MATERIALIZED VIEW CONCURRENTLY " + view,
				); err != nil {
					log.Println(err)
				}
			}

			// Once points is refreshed, award points based cheevos.
			if _, err := db.Exec(
				`INSERT INTO trophies(user_id, trophy)
				      SELECT user_id, 'big-brother'::cheevo
				        FROM points
				       WHERE points >= 1984
				   UNION ALL
				      SELECT user_id, 'its-over-9000'
				        FROM points
				       WHERE points > 9000
				   UNION ALL
				      SELECT user_id, 'twenty-kiloleagues'
				        FROM points
				       WHERE points >= 20000
				   UNION ALL
				      SELECT user_id, 'marathon-runner'
				        FROM points
				       WHERE points >= 42195
				   UNION ALL
				      SELECT user_id, '0xdead'
				        FROM points
				       WHERE points >= 57005
				 ON CONFLICT DO NOTHING`,
			); err != nil {
				log.Println(err)
			}
		}
	}()

	// Every 5 minutes.
	go func() {
		for range time.Tick(5 * time.Minute) {
			// Various GitHub API requests.
			github.Run(db, false)

			if err := github.Wiki(db); err != nil {
				log.Println(err)
			}

			if err := discord.AwardRoles(db); err != nil {
				log.Println(err)
			}
		}
	}()

	// Every hour.
	go func() {
		for range time.Tick(time.Hour) {
			// Update GitHub usernames.
			github.Run(db, true)

			for _, job := range [...]struct{ name, sql string }{
				{
					"expired sessions",
					`DELETE FROM sessions
					  WHERE last_used < TIMEZONE('UTC', NOW()) - INTERVAL '30 days'`,
				},
				{
					"superfluous users",
					`DELETE FROM users u
					  WHERE NOT EXISTS (SELECT FROM sessions WHERE user_id = u.id)
						AND NOT EXISTS (SELECT FROM trophies WHERE user_id = u.id)`,
				},
				{
					"users scheduled for deletion",
					"DELETE FROM users WHERE delete < TIMEZONE('UTC', NOW())",
				},
			} {
				if res, err := db.Exec(job.sql); err != nil {
					log.Println(err)
				} else if rows, _ := res.RowsAffected(); rows != 0 {
					log.Printf("Deleted %d %s\n", rows, job.name)
				}
			}
		}
	}()

	log.Println("Listening…")

	// Dev.
	if _, dev := os.LookupEnv("DEV"); dev {
		// Redirect HTTP to HTTPS.
		go func() {
			panic(http.ListenAndServe(":80",
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, "https://localhost"+r.RequestURI,
						http.StatusMovedPermanently)
				})))
		}()

		// Serve HTTPS.
		panic(http.ListenAndServeTLS(
			":443", "localhost.pem", "localhost-key.pem", routes.Router(db)))
	}

	// Live only listens on HTTP, TLS is handled by Caddy.
	panic(http.ListenAndServe(":80", routes.Router(db)))
}
