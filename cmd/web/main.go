package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse commandline flags passed. -addr flag value will be assigned to addr variable.
	flag.Parse()

	// create logger for writing informational and error messages
	// include log.Lshortfile flag to include relevant file name and line number
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialize application struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// load config
	config := app.loadConfig("./")

	// open database connection pool (only in main to save connection resources)
	db, err := openDB(config.DBDriver, config.DBSource)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	// defer closing the database connection until main exits
	defer db.Close()

	// Initialize a new http.Server struct. Set the Addr and Handler fields
	// the same as before, but add errorLog to ErrorLog.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routeMux(), // returns mux
	}

	infoLog.Printf("Starting server on %s", *addr)

	// mux is treated as a chained interface
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// define an application struct to hold application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func openDB(DBDriver string, DBSource string) (*sql.DB, error) {
	db, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		return nil, err
	}
	// check if connection is still alive
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
