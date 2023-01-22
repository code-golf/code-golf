package github

import (
	"bytes"
	"io"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// TODO Make our own minimal wiki that doesn't use GitHub. Reward contrib.
func Wiki(db *sqlx.DB) error {
	fs := memfs.New()

	if _, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		Depth: 1,
		URL:   "https://github.com/code-golf/code-golf.wiki",
	}); err != nil {
		return err
	}

	files, err := fs.ReadDir("/")
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("TRUNCATE wiki"); err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO wiki VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), ".md")
		slug := config.ID(name)
		section := ""

		if lang := config.LangByID[slug]; lang != nil {
			slug = "langs/" + slug
			section = "Languages"
			name = lang.Name
		} else if slug == "hole-specific-tips" {
			name = "Hole Specific Tips"
		} else if slug == "other-sites" {
			name = "Other Sites"
		} else {
			continue
		}

		f, err := fs.Open(file.Name())
		if err != nil {
			return err
		}

		md, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		var html bytes.Buffer
		if err := markdown.Convert(md, &html); err != nil {
			return err
		}

		if _, err := stmt.Exec(slug, section, name, html.String()); err != nil {
			return err
		}
	}

	return tx.Commit()
}
