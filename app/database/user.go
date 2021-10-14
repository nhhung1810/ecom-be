package database

import (
	"database/sql"
	"ecom-be/app/auth"
	"errors"

	"github.com/lib/pq"
)

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
