package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"math/rand"
	"net/http"
	"syscall"
	"time"

	"github.com/code-golf/code-golf/cookie"
	"github.com/code-golf/code-golf/github"
	"github.com/code-golf/code-golf/golfer"
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

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectHost)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Compress(5))
	r.Use(middleware.WithValue("db", db))

	r.NotFound(routes.NotFound)

	r.Get("/", routes.Index)
	r.Get("/{hole}", routes.GETHole)
	r.Get("/about", routes.About)
	r.Get("/assets/{asset}", routes.Asset)
	r.Get("/callback", routes.Callback)
	r.Get("/favicon.svg", routes.Asset)
	r.Get("/favicon16.png", routes.Asset)
	r.Get("/favicon32.png", routes.Asset)
	r.Get("/feeds/{feed}", routes.Feed)
	r.Get("/ideas", routes.Ideas)
	r.Get("/log-out", routes.LogOut)
	r.Get("/random", routes.Random)
	r.Get("/recent", routes.Recent)
	r.Get("/recent/{lang}", routes.Recent)
	r.Get("/robots.txt", routes.Robots)
	r.Get("/scores/{hole}/{lang}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/{suffix}", routes.Scores)
	r.Get("/settings", routes.Settings)
	r.Post("/solution", routes.Solution)
	r.Get("/stats", routes.Stats)
	r.Get("/users/{user}", routes.User)

	r.Route("/admin", func(r chi.Router) {
		// TODO Have previous middleware put Golfer into the context.
		// TODO Have an admin boolean on said Golfer.
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, login := cookie.Read(r); login == "JRaspass" {
					next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			})
		})

		r.Get("/", routes.Admin)
		r.Get("/solutions", routes.AdminSolutions)
	})

	r.Route("/golfers/{name}", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if info := golfer.GetInfo(db, chi.URLParam(r, "name")); info != nil {
					ctx := context.WithValue(r.Context(), "golferInfo", info)
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					routes.NotFound(w, r)
				}
			})
		})

		r.Get("/", routes.Golfer)
		r.Get("/holes", routes.GolferHoles)
	})

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
