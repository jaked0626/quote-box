package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/jaked0626/snippetbox/internal/db/models"
	"github.com/jaked0626/snippetbox/ui"
)

type templateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Toast           string
	IsAuthenticated bool
	CSRFToken       string // to allow form submissions from valid requests
}

func formatDate(rawDate *time.Time) (formattedDate string) {
	// Specify the desired timezone (Asia/Tokyo)
	japanTimezone, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		formattedDate = rawDate.Format("02 Jan 2006 at 15:04")
		return
	}
	japanTime := rawDate.In(japanTimezone)
	formattedDate = japanTime.Format("02 Jan 2006 at 15:04 JST")
	return
}

var functions = template.FuncMap{
	"formatDate": formatDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// get paths of all template files
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}
		templateSet, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
