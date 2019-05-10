package routes

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/JRaspass/code-golf/cookie"
)

var tmpl = template.New("")

func init() {
	// Tests run from the package directory, walk upwards until we find views.
	for {
		if _, err := os.Stat("views"); os.IsNotExist(err) {
			os.Chdir("..")
		} else {
			break
		}
	}

	if err := filepath.Walk("views", func(path string, _ os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			if b, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				name := path[len("views/") : len(path)-len(".html")]
				tmpl = template.Must(tmpl.New(name).Parse(string(b)))
			}
		}

		return err
	}); err != nil {
		panic(err)
	}
}

// Render wraps common logic required for rendering a view to the user.
func Render(w http.ResponseWriter, r *http.Request, code int, name string, data interface{}) {
	header := w.Header()

	header["Content-Language"] = []string{"en"}
	header["Content-Type"] = []string{"text/html;charset=utf-8"}
	header["Referrer-Policy"] = []string{"no-referrer"}
	header["X-Content-Type-Options"] = []string{"nosniff"}
	header["X-Frame-Options"] = []string{"DENY"}
	header["Content-Security-Policy"] = []string{
		"base-uri 'none';" +
			"connect-src 'self';" +
			"default-src 'none';" +
			"form-action 'none';" +
			"frame-ancestors 'none';" +
			"img-src 'self' data: avatars.githubusercontent.com;" +
			"script-src 'self';" +
			"style-src 'self'",
	}

	args := struct {
		CommonCssPath, Login, LoginURL string
		Data                           interface{}
	}{
		CommonCssPath: commonCssPath,
		Data:          data,
	}

	if _, args.Login = cookie.Read(r); args.Login == "" {
		args.LoginURL = "//github.com/login/oauth/authorize?client_id=7f6709819023e9215205&scope=user:email&redirect_uri=https://code-golf.io/callback?redirect_uri%3D" + url.QueryEscape(url.QueryEscape(r.RequestURI))
	}

	w.WriteHeader(code)

	if err := tmpl.ExecuteTemplate(w, name, args); err != nil {
		panic(err)
	}
}
