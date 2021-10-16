package product

import (
	"github.com/gin-gonic/gin"
)

//provide function access to db
type Repository interface {
	AddProduct(p Product) error
	FetchProduct(id int) (*Product, error)
}

// Provide interface for product operation in handler
type Service interface {
	AddProduct(p Product) error
	FetchProduct(id int) (*Product, error)
	ParseProduct(g *gin.Context) (*Product, error)
}

// Abstract layer, implementing the service
type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddProduct(p Product) error {
	// TODO: VALIDATE ALL INPUT
	err := s.r.AddProduct(p)
	if err != nil {
		return nil
	}
	return nil
}

func (s *service) FetchProduct(id int) (*Product, error) {
	p, err := s.r.FetchProduct(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) ParseProduct(c *gin.Context) (*Product, error) {
	var product Product
	if err := c.BindJSON(&product); err != nil {
		return nil, err
	}
	return &product, nil
}
