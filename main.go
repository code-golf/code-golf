package main

import (
	"crypto/tls"
	"database/sql"
	"math/rand"
	"net/http"
	"syscall"
	"time"

	"github.com/code-golf/code-golf/github"
	"github.com/code-golf/code-golf/middleware"
	"github.com/code-golf/code-golf/routes"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	db, err := sql.Open("postgres", "")
	if err != nil {
		panic(err)
	}

	var r = chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RedirectHost,
		middleware.Public,
		middleware.RedirectSlashes,
		middleware.Compress(5),
		middleware.DatabaseHandler(db),
		middleware.GolferHandler,
	)

	r.NotFound(routes.NotFound)

	r.Get("/", routes.Index)
	r.Get("/{hole}", routes.Hole)
	r.Get("/about", routes.About)
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AdminArea)
		r.Get("/", routes.Admin)
		r.Get("/solutions", routes.AdminSolutions)
		r.Get("/solutions/run", routes.AdminSolutionsRun)
	})
	r.Get("/assets/{asset}", routes.Asset)
	r.Get("/callback", routes.Callback)
	r.Get("/feeds/{feed}", routes.Feed)
	r.Route("/golfer", func(r chi.Router) {
		r.Use(middleware.GolferArea)
		r.Get("/settings", routes.GolferSettings)
	})
	r.Route("/golfers/{name}", func(r chi.Router) {
		r.Use(middleware.GolferInfoHandler)
		r.Get("/", routes.GolferTrophies)
		r.Get("/achievements", routes.GolferAchievements)
		r.Get("/holes", routes.GolferHoles)
	})
	r.Get("/ideas", routes.Ideas)
	r.Get("/log-out", routes.LogOut)
	r.Get("/random", routes.Random)
	r.Get("/recent", routes.Recent)
	r.Get("/recent/{lang}", routes.Recent)
	r.Get("/scores/{hole}/{lang}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/{suffix}", routes.Scores)
	r.Get("/sitemap.xml", routes.Sitemap)
	r.Post("/solution", routes.Solution)
	r.Get("/stats", routes.Stats)
	r.Get("/users/{name}", routes.User)

	certManager := autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"code-golf.io", "code.golf", "www.code-golf.io", "www.code.golf",
		),
	}

	server := &http.Server{
		Addr:    ":1443",
		Handler: r,
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // HTTP/2-required.
			},
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	var crt, key string
	if _, dev := syscall.Getenv("DEV"); dev {
		crt = "localhost.pem"
		key = "localhost-key.pem"
	} else {
		server.TLSConfig.GetCertificate = certManager.GetCertificate
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)

		for {
			<-ticker.C
			github.Ideas(db)
			github.PullRequests(db)
			github.Stars(db)
		}
	}()

	println("Listening…")

	// Redirect HTTP to HTTPS and handle ACME challenges.
	go func() {
		panic(http.ListenAndServe(":1080", certManager.HTTPHandler(nil)))
	}()

	// Serve HTTPS.
	panic(server.ListenAndServeTLS(crt, key))
}
