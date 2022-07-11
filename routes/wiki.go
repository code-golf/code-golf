package routes

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

type (
	wikiPage struct{ Slug, Section, Name string }
	wikiData struct {
		HTML       template.HTML
		Slug, Name string
		Pages      []wikiPage
	}
)

// GET /wiki
func wikiGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "wiki", wikiData{Pages: pages(r)}, "Wiki")
}

// GET /wiki/*
func wikiPageGET(w http.ResponseWriter, r *http.Request) {
	data := wikiData{Pages: pages(r), Slug: param(r, "*")}

	if err := session.Database(r).QueryRow(
		"SELECT html, name FROM wiki WHERE slug = $1", data.Slug,
	).Scan(&data.HTML, &data.Name); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	render(w, r, "wiki", data, "Wiki: "+data.Name)
}

func pages(r *http.Request) (pages []wikiPage) {
	pages = append(pages, wikiPage{Name: "Home"})

	rows, err := session.Database(r).Query(
		"SELECT slug, section, name FROM wiki ORDER BY section, name")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p wikiPage

		if err := rows.Scan(&p.Slug, &p.Section, &p.Name); err != nil {
			panic(err)
		}

		pages = append(pages, p)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return
}
