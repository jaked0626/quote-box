package db

import "database/sql"

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
