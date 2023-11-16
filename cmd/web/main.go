package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jaked0626/snippetbox/internal/config"
)

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

func main() {
	// BEST PRACTICE: all fatal or panic error logs should be called from within main

	config := config.LoadConfig()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// application struct holds application wide variables
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// only open in main to save connection resources
	db, err := openDB(config.DBDriver, config.DBSource)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr:     config.Addr,
		ErrorLog: errorLog,
		Handler:  app.routeMux(),
	}

	infoLog.Printf("Starting server on %s", config.Addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
