package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"
)

const headerHSTS = "max-age=31536000;includeSubDomains;preload"

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()["Content-Language"] = []string{"en"}
	w.Header()["Date"] = nil

	defer func(start time.Time) {
		fmt.Printf(
			"%3dms %s %s %s\n",
			time.Since(start).Nanoseconds() / 1e6,
			r.Method,
			r.URL.Path,
			r.UserAgent(),
		)

		if r := recover(); r != nil {
			fmt.Fprint(os.Stderr, "<1>", r, "\n")
			http.Error(w, "500: It's Dead, Jim.", 500)
		}
	}(time.Now())

	switch r.Host {
	case "code-golf.io":
		codeGolf(w, r)
	case "raspass.me":
		raspass(w, r)
	case "www.code-golf.io", "www.raspass.me":
		http.Redirect(w, r, "https://"+r.Host[4:]+r.URL.String(), 301)
	}
}

func mustLoadX509KeyPair(certFile, keyFile string) tls.Certificate {
	if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		panic(err)
	} else {
		return cert
	}
}

func main() {
	server := &http.Server{
		Handler: &handler{},
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				mustLoadX509KeyPair(
					"/home/jraspass/dehydrated/certs/code-golf.io/fullchain.pem",
					"/home/jraspass/dehydrated/certs/code-golf.io/privkey.pem",
				),
				mustLoadX509KeyPair(
					"/home/jraspass/dehydrated/certs/raspass.me/fullchain.pem",
					"/home/jraspass/dehydrated/certs/raspass.me/privkey.pem",
				),
			},
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

	server.TLSConfig.BuildNameToCertificate()

	println("Listening…")

	// Redirect HTTP to HTTPS.
	go func() {
		panic(http.ListenAndServe(
			":80",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+r.Host+r.URL.String(), 301)
			}),
		))
	}()

	// Serve HTTPS.
	panic(server.ListenAndServeTLS("", ""))
}
