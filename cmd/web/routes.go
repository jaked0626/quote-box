package main

import "net/http"

func (app *application) routeMux() *http.ServeMux {
	// initialize servemux and map routes to handlers
	mux := http.NewServeMux()

	// file server for static files in ui/static/
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// use mux.Handle() function to register the file server as the
	// handler for all URL paths that start with "/static/". For matching paths,
	// we strip the "/static" prefix before the request reaches the file server
	// otherwise there will be two /statics
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// HandleFunc coerces functions into interfaces that satisfy
	// the method ServeHTTP()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/view", app.snippetView)
	mux.HandleFunc("/create", app.snippetCreate)

	return mux
}
