package main

import (
	"html/template"
	"path/filepath"

	"github.com/jaked0626/snippetbox/internal/db/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// get paths of all template files
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateFiles := []string{
			"./ui/html/base.html",
			page,
		}
		templateSet, err := template.ParseFiles(templateFiles...)
		if err != nil {
			return nil, err
		}

		// add partials
		templateSet.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
