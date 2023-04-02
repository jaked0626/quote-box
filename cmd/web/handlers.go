package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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
	// extract value of id param from url query string
	// make sure id is a postive integer
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Displaying snippet with id %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
