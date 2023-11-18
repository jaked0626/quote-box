package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/jaked0626/snippetbox/internal/models"
)

// Home handler method for application struct defined in main
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// display home page when url is exactly '/', otherwise
	// redirect to 404
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// use template.ParseFiles() to read template files into template set
	// path is set relative to root directory, meaning we must run go mod run in
	// the root of our project

	templateFiles := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	templateSet, err := template.ParseFiles(templateFiles...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// serve template set
	err = templateSet.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// validate input
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.badRequest(w)
		return
	} else if id < 1 {
		app.notFound(w)
		return
	}

	// get from db
	s, err := app.snippets.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// response
	res, err := json.Marshal(s)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return
}

func (app *application) snippetList(w http.ResponseWriter, r *http.Request) {
	// validate input
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		app.badRequest(w)
		return
	} else if limit < 1 {
		app.notFound(w)
		return
	}

	// get from db
	snippets, err := app.snippets.List(limit)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// response
	res, err := json.Marshal(snippets)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := r.URL.Query().Get("title")
	content := r.URL.Query().Get("content")
	expires, err := strconv.Atoi(r.URL.Query().Get("expires"))
	if err != nil {
		app.badRequest(w)
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
