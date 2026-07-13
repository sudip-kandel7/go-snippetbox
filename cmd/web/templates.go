package main

import (
	"html/template"
	"path/filepath"
	"time"
	"go-snippetbox.kandel.net/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 03:04 PM")
}

var funcs = template.FuncMap {
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error ){
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(funcs).ParseFiles("./ui/html/base.tmpl")
		
		if err != nil {
			return nil , err
		}

		// files := []string {
		// 	"./ui/html/base.tmpl",
		// 	"./ui/html/partials/nav.tmpl",
		// 	page,
		// }

		ts, err =  ts.ParseGlob("./ui/html/partials/*.tmpl")

		if err != nil {
			return  nil, err
		}

		ts, err = ts.ParseFiles(page)

		if err != nil {
			return  nil, err
		}

		cache[name] = ts

	}

	return cache, nil

}