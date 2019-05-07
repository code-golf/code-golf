package main

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/JRaspass/code-golf/routes"
	brotli "github.com/cv-library/negroni-brotli"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

type err500 struct{}

func (*err500) FormatPanicError(w http.ResponseWriter, r *http.Request, _ *negroni.PanicInformation) {
	http.Error(w, "500: It's Dead, Jim.", 500)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if err := routes.ParseViews(); err != nil {
		panic(err)
	}

	certManager := autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"code-golf.io", "www.code-golf.io",
		),
	}

	logger := negroni.NewLogger()
	logger.ALogger = log.New(os.Stdout, "", 0)
	logger.SetFormat("{{.StartTime}} {{.Status}} {{.Method}} {{.Request.URL}} {{.Request.UserAgent}}")

	recovery := negroni.NewRecovery()
	recovery.Formatter = &err500{}
	recovery.Logger = log.New(os.Stderr, "<1>", 0)

	server := &http.Server{
		Handler: negroni.New(
			logger,
			brotli.New(5),
			recovery,
			negroni.Wrap(routes.Router),
		),
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

	_, dev := syscall.Getenv("DEV")

	if !dev {
		server.TLSConfig.GetCertificate = certManager.GetCertificate
	}

	server.TLSConfig.BuildNameToCertificate()

	go func() {
		ticker := time.NewTicker(time.Minute)

		for {
			<-ticker.C
			routes.Stars()
		}
	}()

	println("Listening…")

	// Redirect HTTP to HTTPS and handle ACME challenges.
	go func() { panic(http.ListenAndServe(":80", certManager.HTTPHandler(nil))) }()

	// Serve HTTPS.
	if dev {
		panic(server.ListenAndServeTLS("localhost.pem", "localhost-key.pem"))
	} else {
		panic(server.ListenAndServeTLS("", ""))
	}
}
