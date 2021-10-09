package models

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// sameple database
var database []Credential

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Credential struct {
	//TODO: Check for UNICODE compatible
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c Credential) CheckExisted() (*User, error) {
	data := database
	// var hashedPassword []byte
	for i := 0; i < len(data); i++ {
		if data[i].Email == c.Email {
			user := &User{
				ID:       i,
				Name:     data[i].Name,
				Email:    data[i].Email,
				Password: data[i].Password,
			}
			return user, nil
		}
	}
	return nil, errors.New("Error: Can not find the user")
}

func (c Credential) AddNewAccount() bool {
	//TODO: connect database
	password, _ := bcrypt.GenerateFromPassword([]byte(c.Password), 14)
	c.Password = string(password)
	database = append(database, c)
	return true
}

// DATABASE FUNCTION
func CountDB() int {
	return len(database)
}

func FetchAccount() []Credential {
	// TODO: connect database + querying
	return nil
}

func ParseCredential(c *gin.Context) Credential {
	var info Credential
	if err := c.BindJSON(&info); err != nil {
		return info
	}
	return info
}

func FindUser(id int) (*User, error) {
	data := database
	if len(data) < id {
		return nil, errors.New("error: Can not find the user")
	}
	user := &User{
		ID:       id,
		Name:     data[id].Name,
		Email:    data[id].Email,
		Password: data[id].Password,
	}

	return user, nil
}


