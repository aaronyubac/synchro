package main

import (
	"path/filepath"
	"synchro/internal/models"
	"text/template"
)

type templateData struct {
	Form any
	User models.User
	Event models.Event
	Events []models.Event
	EventUnavailabilities []models.Unavailability
	Flash string
	IsAuthenticated bool
}


func newTemplateCache() (map[string]*template.Template, error) {

	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		
		ts, err := template.New(name).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

