package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jaked0626/snippetbox/internal/config"
	"github.com/jaked0626/snippetbox/internal/db/dbutils"
	"github.com/jaked0626/snippetbox/internal/db/models"
)

// define an application struct to hold application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
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
	infoLog := log.New(os.Stdout, "API INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "API ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// only open in main to save connection resources
	db, err := dbutils.OpenDB(config.DBDriver, config.DBSource)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// application struct holds application wide variables
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     config.ApiAddr,
		ErrorLog: errorLog,
		Handler:  app.routeMux(),
	}

	infoLog.Printf("Starting API server on %s", config.ApiAddr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
