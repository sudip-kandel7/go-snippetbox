package main

import (
	"html/template"
	"io/fs"
	"path"
	"time"

	"go-snippetbox.kandel.net/internal/models"
	"go-snippetbox.kandel.net/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 03:04 PM")
}

var funcs = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := path.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(funcs).ParseFS(ui.Files, patterns...)

		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}

	return cache, nil

}
