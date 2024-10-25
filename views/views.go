package views

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/pretty"
)

//go:embed html/* svg/*
var views embed.FS

var tmpl = template.New("").Funcs(template.FuncMap{
	"bytes":     pretty.Bytes,
	"comma":     pretty.Comma,
	"dec":       func(i int) int { return i - 1 },
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"html":      func(html string) template.HTML { return template.HTML(html) },
	"inc":       func(i int) int { return i + 1 },
	"map":       func() map[string]bool { return map[string]bool{} },
	"ord":       pretty.Ordinal,
	"page":      func(i int) int { return i/pager.PerPage + 1 },
	"title":     pretty.Title,
	"time":      pretty.Time,

	"getSet": func(m map[string]bool, k string) bool {
		old := m[k]
		m[k] = true
		return old
	},

	"hasField": func(v any, name string) bool {
		s := reflect.ValueOf(v)
		return s.Kind() == reflect.Struct && s.FieldByName(name).IsValid()
	},

	"param": func(r *http.Request, key string) string {
		value, _ := url.QueryUnescape(r.PathValue(key))
		return value
	},

	"svg": func(name string) template.HTML {
		data, _ := views.ReadFile("svg/" + name + ".svg")
		return template.HTML(data)
	},
})

func init() {
	// Manually walk & parse HTML templates so we can support directories.
	if err := fs.WalkDir(views, "html", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			data, err := views.ReadFile(path)
			if err != nil {
				return err
			}

			// Trim the directory prefix and extension suffix.
			name := path[len("html/") : len(path)-len(".html")]
			tmpl = template.Must(tmpl.New(name).Parse(string(data)))
		}

		return err
	}); err != nil {
		panic(err)
	}
}

func Render(w http.ResponseWriter, name string, data any) error {
	return tmpl.ExecuteTemplate(w, name, data)
}
