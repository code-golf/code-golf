package routes

import (
	"encoding/xml"
	"net/http"
	"net/url"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /sitemap.xml
func sitemapGET(w http.ResponseWriter, r *http.Request) {
	type URL struct {
		Loc string `xml:"loc"`
	}

	type URLSet struct {
		XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
		URLs    []URL    `xml:"url"`
	}

	sitemap := URLSet{
		URLs: []URL{
			{"https://code.golf"},
			{"https://code.golf/about"},
			{"https://code.golf/ideas"},
			{"https://code.golf/stats"},
			{"https://code.golf/wiki"},
		},
	}

	rows, err := session.Database(r).Query(
		"SELECT 'https://code.golf/wiki/' || slug FROM wiki ORDER BY slug")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var url URL

		if err := rows.Scan(&url.Loc); err != nil {
			panic(err)
		}

		sitemap.URLs = append(sitemap.URLs, url)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	for _, hole := range config.HoleList {
		sitemap.URLs = append(
			sitemap.URLs, URL{"https://code.golf/" + url.PathEscape(hole.ID)})
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write([]byte(xml.Header))

	if err := xml.NewEncoder(w).Encode(sitemap); err != nil {
		panic(err)
	}
}
