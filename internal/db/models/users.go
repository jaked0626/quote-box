package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int
	Name            string
	Email           string
	HashedPasssword []byte
	Created         time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name string, email string, password string) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	qry := `INSERT INTO users (name, email, hashed_password, created) VALUES (
		 	  $1,
		 	  $2,
		 	  $3,
		 	  CURRENT_TIMESTAMP);`

	_, err = m.DB.Exec(qry, name, email, string(hashedPassword))
	if err != nil {
		// Check if the error is due to a unique constraint violation
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email string, password string) (id int, err error) {
	var hashedPassword []byte

	qry := "SELECT id, hashed_password FROM users WHERE email = $1"
	row := m.DB.QueryRow(qry, email)
	err = row.Scan(&id, &hashedPassword)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrInvalidCredentials
		} else {
			return -1, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return -1, ErrInvalidCredentials
		} else {
			return -1, err
		}
	}

	return id, err
}

func (m *UserModel) Exists(id int) (exists bool, err error) {
	qry := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"
	row := m.DB.QueryRow(qry, id)
	err = row.Scan(&exists)

	return exists, err
}
