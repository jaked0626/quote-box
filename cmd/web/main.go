package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jaked0626/snippetbox/internal/config"
	"github.com/jaked0626/snippetbox/internal/db/dbutils"
	"github.com/jaked0626/snippetbox/internal/db/models"
	_ "github.com/lib/pq"
)

// define an application struct to hold application-wide dependencies
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	cache          map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	config := config.LoadConfig()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// database: only in main to save connection resources
	db, err := dbutils.OpenDB(config.DBDriver, config.DBSource)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// caching
	cache, err := newTemplateCache()
	if err != nil {
		errorLog.Printf("Cache cannot be initialized: %v", err)
	}

	// form decoding
	formDecoder := form.NewDecoder()

	// session management
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// application wide dependencies
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		cache:          cache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// tls: only use elliptic curves with assembly implementations for performance
	// tlsConfig := &tls.Config{
	// 	CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	// }

	// server
	srv := &http.Server{
		Addr:         config.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routeMux(),
		IdleTimeout:  2 * time.Minute, // verbose but for some cursed reason IdleTimeout defaults to the same value as ReadTimeout!?
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// TLSConfig:    tlsConfig,
	}

	infoLog.Printf("Starting server on %s", config.Addr)

	// err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
