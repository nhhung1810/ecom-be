package database

import (
	"database/sql"
	"errors"
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
