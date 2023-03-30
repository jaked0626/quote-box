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
	app.errorLog.Println(trace)

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
