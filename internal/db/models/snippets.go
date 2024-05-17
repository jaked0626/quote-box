package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID       int
	Title    string
	Author   string
	Work     string
	Content  string
	Created  time.Time
	Expires  time.Time
	UserID   int
	UserName string
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, author string, work string, content string, expires int, user_id int) (id int, err error) {
	qry := `INSERT INTO snippets (title, author, work, content, created, expires, user_id) VALUES (
			 $1,
			 $2,
			 $3,
			 $4,
			 CURRENT_TIMESTAMP,
			 CURRENT_TIMESTAMP + $5 * INTERVAL '1 day',
			 $6
			 ) RETURNING id;`

	row := m.DB.QueryRow(qry, title, author, work, content, expires, user_id)
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, err
}

func (m *SnippetModel) Get(id int) (s *Snippet, err error) {
	// query db
	qry := `SELECT 
		snippets.id,
		snippets.title,
		snippets.author,
		snippets.work,
		snippets.content,
		snippets.created,
		snippets.expires,
		snippets.user_id,
		users.name
	FROM snippets
	LEFT JOIN users ON snippets.user_id = users.id
	WHERE expires > CURRENT_TIMESTAMP AND snippets.id = $1;`
	row := m.DB.QueryRow(qry, id)

	// unmarshal
	s = &Snippet{}
	err = row.Scan(&s.ID, &s.Title, &s.Author, &s.Work, &s.Content, &s.Created, &s.Expires, &s.UserID, &s.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) List(limit int) (snippets []*Snippet, err error) {
	// query db
	qry := `SELECT *
	FROM snippets
	WHERE expires > CURRENT_TIMESTAMP
	ORDER BY created DESC
	LIMIT $1`
	rows, err := m.DB.Query(qry, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// unmarshal query results
	snippets = []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Author, &s.Work, &s.Content, &s.Created, &s.Expires, &s.UserID)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
