package database

import (
	"database/sql"
	"ecom-be/app/config"
	"errors"
)

var errExistingAccount = errors.New("error: user existed")
var errNotExistedAccount = errors.New("error: user not exist")
var errNotExistedProduct = errors.New("error: user not exist")
var errUnknown = errors.New("error: unknown error")

type Storage struct {
	db *sql.DB
}

// make database
func NewStorage() (*Storage, error) {
	connStr := config.DefaultConfig.ConnString
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}
