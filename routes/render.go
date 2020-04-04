package routes

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/code-golf/code-golf/cookie"
	"github.com/code-golf/code-golf/pretty"
)

func colour(i int) string {
	if i <= 1 {
		return "yellow"
	}
	if i <= 2 {
		return "orange"
	}
	if i <= 3 {
		return "red"
	}
	if i <= 10 {
		return "purple"
	}
	if i <= 100 {
		return "blue"
	}
	return "green"
}

var tmpl = template.New("").Funcs(template.FuncMap{
	"colour":    colour,
	"comma":     pretty.Comma,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"ord":       pretty.Ordinal,
	"symbol": func(name string) template.HTML {
		if svg, err := ioutil.ReadFile("views/" + name + ".svg"); err != nil {
			panic(err)
		} else {
			return template.HTML(bytes.ReplaceAll(svg, []byte("svg"), []byte("symbol")))
		}
	},
	"title": strings.Title,
	"time":  pretty.Time,
})

func init() {
	// Tests run from the package directory, walk upwards until we find views.
	for {
		if _, err := os.Stat("views"); os.IsNotExist(err) {
			os.Chdir("..")
		} else {
			break
		}
	}

	if err := filepath.Walk("views", func(file string, _ os.FileInfo, err error) error {
		if ext := path.Ext(file); ext == ".html" || ext == ".svg" {
			if b, err := ioutil.ReadFile(file); err != nil {
				return err
			} else {
				name := file[len("views/") : len(file)-len(ext)]
				tmpl = template.Must(tmpl.New(name).Parse(string(b)))
			}
		}

		return err
	}); err != nil {
		panic(err)
	}
}

func render(
	w http.ResponseWriter,
	r *http.Request,
	code int,
	name, title string,
	data interface{},
) {
	// The generated value SHOULD be at least 128 bits long (before encoding),
	// and SHOULD be generated via a cryptographically secure random number
	// generator - https://w3c.github.io/webappsec-csp/#security-nonces
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	args := struct {
		CommonCssPath, Login, LogInURL, Nonce, Path, Title string
		Data                                               interface{}
		Request                                            *http.Request
	}{
		CommonCssPath: commonCssPath,
		Data:          data,
		Nonce:         base64.StdEncoding.EncodeToString(nonce),
		Path:          r.URL.Path,
		Request:       r,
		Title:         title,
	}

	header := w.Header()

	header.Set("Content-Language", "en")
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("X-Frame-Options", "DENY")
	header.Set("Content-Security-Policy",
		"base-uri 'none';"+
			"connect-src 'self';"+
			"default-src 'none';"+
			"form-action 'none';"+
			"font-src 'self';"+
			"frame-ancestors 'none';"+
			"img-src 'self' data: avatars.githubusercontent.com;"+
			"script-src 'self' 'nonce-"+args.Nonce+"';"+
			"style-src 'self' 'nonce-"+args.Nonce+"'",
	)

	if _, args.Login = cookie.Read(r); args.Login == "" {
		// Shallow copy because we want to modify a string.
		config := config

		config.RedirectURL = "https://code-golf.io/callback?redirect_uri=" +
			url.QueryEscape(r.RequestURI)

		// TODO State is a token to protect the user from CSRF attacks.
		args.LogInURL = config.AuthCodeURL("")
	}

	w.WriteHeader(code)

	if err := tmpl.ExecuteTemplate(w, name, args); err != nil {
		panic(err)
	}
}
