package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse commandline flags passed. -addr flag value will be assigned to addr variable.
	flag.Parse()

	// initialize servemux and map routes to handlers
	mux := http.NewServeMux()

	// HandleFunc coerces functions into interfaces that satisfy
	// the method ServeHTTP()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/view", snippetView)
	mux.HandleFunc("/create", snippetCreate)

	// file server for static files in ui/static/
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// use mux.Handle() function to register the file server as the
	// handler for all URL paths that start with "/static/". For matching paths,
	// we strip the "/static" prefix before the request reaches the file server
	// otherwise there will be two /statics
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// set up server. Pass in two parameters: TCP netweork address to listen to
	// (:4000) and servemux created.
	log.Printf("Starting server on %s", *addr)

	// mux is treated as a chained interface
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
