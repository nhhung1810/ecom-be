package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ERROR DEFINITION

//Provide function access to database
type Repository interface {
	AddAccount(user User) error
	FindUserByEmail(email string) (*User, error)
	FindUserByID(id int) (*User, error)
}

// Provide interface for auth operation in handler
type Service interface {
	FindUserByID(id int) (*User, error)
	AddAccount(user Credential) error
	CheckExisted(user Credential) (*User, error)
	ParseCredential(g *gin.Context) (*Credential, error)
}

// Abstract layer, implementing the service
type service struct {
	r Repository
}

// Implementing Service interface with "service"
func NewService(r Repository) Service {
	return &service{r}
}

// Maybe change to return the ID of user? Not sure
func (s *service) AddAccount(req Credential) error {
	password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	req.Password = string(password)
	newUser := &User{
		// SHOULD WE ADD ID HERE ?
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}

	err := s.r.AddAccount(*newUser)
	return err
}

func (s *service) CheckExisted(req Credential) (*User, error) {
	result, err := s.r.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) FindUserByID(id int) (*User, error) {
	result, err := s.r.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) ParseCredential(c *gin.Context) (*Credential, error) {
	var info Credential
	if err := c.BindJSON(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
