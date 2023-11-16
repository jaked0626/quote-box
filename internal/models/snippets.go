package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (id int, err error) {
	id = -1
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES (
			 $1,
			 $2,
			 CURRENT_TIMESTAMP,
			 CURRENT_TIMESTAMP + $3 * INTERVAL '1 day'
			 ) RETURNING id;`

	result := m.DB.QueryRow(stmt, title, content, expires)
	err = result.Scan(&id)
	if err != nil {
		return
	}

	return
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) List(length int) ([]*Snippet, error) {
	return nil, nil
}
