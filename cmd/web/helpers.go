package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

	// write response to buffer first to see if it compiles without error
	buf := new(bytes.Buffer)
	err := templateSet.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// if the above succeeds, we can return a successful header status and write
	// to the http response writer
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) (data *templateData) {
	data = &templateData{
		CurrentYear: time.Now().Year(),
		Toast:       app.sessionManager.PopString(r.Context(), "toast"),
	}
	return
}
