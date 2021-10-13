package models

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	// mo hinh 3 lop
}

type Credential struct {
	//TODO: Check for UNICODE compatible
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c Credential) CheckExisted() (*User, error) {
	user, err := FindUserByEmail(c.Email)
	if err != nil {
		return nil, errors.New("error: Can not find the user")
	}

	return user, nil
}

func (c Credential) AddNewAccount() bool {
	//TODO: connect database
	password, _ := bcrypt.GenerateFromPassword([]byte(c.Password), 14)
	c.Password = string(password)
	user := &User{
		Name:     c.Name,
		Password: c.Password,
		Email:    c.Email,
	}

	AddUser(user)
	return true
}

// DATABASE FUNCTION
func CountDB() int {
	var users []User
	result := db.Find(&users)
	return int(result.RowsAffected)
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

func AddUser(user *User) {
	db.Create(user) 
	//context api
}

func FindUserByEmail(email string) (*User, error) {
	var user User
	db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return nil, errors.New("error: no such user")
	}
	return &user, nil
}

func FindUserByID(id int) (*User, error) {
	var user User
	db.First(&user, id)
	if user.ID == 0 {
		return nil, errors.New("error: no such user")
	}
	return &user, nil
}
