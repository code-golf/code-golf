package main

import (
	"crypto/tls"
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
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	log.SetFlags(log.Ltime)

	rand.Seed(time.Now().UnixNano())

	db, err := sql.Open("postgres", "")
	if err != nil {
		panic(err)
	}

	certManager := autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"code.golf", "www.code.golf",
			"code-golf.io", "www.code-golf.io", // Legacy domain.
		),
	}

	server := http.Server{
		Addr:    ":1443",
		Handler: routes.Router(db),
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // HTTP/2-required.
			},
			CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
			MinVersion:       tls.VersionTLS12,
		},
	}

	var crt, key string
	if _, dev := os.LookupEnv("DEV"); dev {
		crt = "localhost.pem"
		key = "localhost-key.pem"
	} else {
		server.TLSConfig.GetCertificate = certManager.GetCertificate
	}

	// Every 30 seconds.
	go func() {
		for range time.Tick(30 * time.Second) {
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
				      SELECT user_id, 'its-over-9000'::cheevo
				        FROM points
				       WHERE points > 9000
				   UNION ALL
				      SELECT user_id, 'twenty-kiloleagues'::cheevo
				        FROM points
				       WHERE points >= 20000
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

	// Redirect HTTP to HTTPS and handle ACME challenges.
	go func() {
		panic(http.ListenAndServe(":1080", certManager.HTTPHandler(nil)))
	}()

	// Serve HTTPS.
	panic(server.ListenAndServeTLS(crt, key))
}
