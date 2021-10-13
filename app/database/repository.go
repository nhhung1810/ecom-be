package database

import (
	"database/sql"
	"ecom-be/app/auth"
	"errors"

	"github.com/lib/pq"
)

var errExistingAccount = errors.New("error: account exist")
var errNotExisted = errors.New("error: account not exist")
var errUnknown = errors.New("error: unknown error")

type Storage struct {
	db *sql.DB
}

// make database
func NewStorage() (*Storage, error) {
	connStr := "host=localhost user=postgres password=admin dbname=ecom port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func (s *Storage) AddAccount(user auth.User) error {
	var testID int
	err := s.db.QueryRow(
		`INSERT INTO Users (name, email, password) 
		VALUES ($1, $2, $3) RETURNING id`, user.Name, user.Email, user.Password).Scan(&testID)

	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return errors.New("add account error: " + errStr)
	}

	println("INSERT ACCOUNT TEST:", err, testID)
	return nil
}

func (s *Storage) FindUserByEmail(email string) (*auth.User, error) {
	var user auth.User
	rows := s.db.QueryRow(`SELECT id, name, email, password 
							FROM users WHERE email = $1`, email)

	err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errNotExisted
		}
		return nil, errUnknown
	}
	return &user, nil
}

func (s *Storage) FindUserByID(id int) (*auth.User, error) {
	var user auth.User
	rows := s.db.QueryRow(`SELECT id, name, email, password 
							FROM users WHERE id = $1`, id)

	err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errNotExisted
		}
		return nil, errUnknown
	}
	return &user, nil
}
