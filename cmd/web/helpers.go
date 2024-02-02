package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError: writes an error message and stack trace to errorLog,
// then sends 500 Internal Server Error response
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError: sends specific status code and corresponding description to user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound: convenience wrapper around clientError that sends a 404 Not Found response
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// badRequest: convenience wrapper around clientError that sends a 400 Bad Request response
func (app *application) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	templateSet, ok := app.cache[page]
	if !ok {
		err := fmt.Errorf("Error: the template %s does not exist in cache", page)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	err := templateSet.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}
