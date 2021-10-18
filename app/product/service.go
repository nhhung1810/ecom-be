package product

import (
	"github.com/gin-gonic/gin"
)

//provide function access to db
type Repository interface {
	AddProduct(p Product, userid int) (*int, error)
	FetchProduct(id int) (*Product, error)
	FetchAllProductsByUser(id int) ([]ProductImage, error)
}

// Provide interface for product operation in handler
type Service interface {
	AddProduct(p Product, userid int) (*int, error)
	FetchProduct(id int) (*Product, error)
	ParseProduct(g *gin.Context) (*Product, error)
	FetchAllProductsByUser(id int) ([]ProductImage, error)
}

// Abstract layer, implementing the service
type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) ParseProduct(c *gin.Context) (*Product, error) {
	var product Product
	if err := c.BindJSON(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *service) AddProduct(p Product, userid int) (*int, error) {
	// TODO: VALIDATE ALL INPUT
	id, err := s.r.AddProduct(p, userid)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *service) FetchProduct(id int) (*Product, error) {
	p, err := s.r.FetchProduct(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) FetchAllProductsByUser(id int) ([]ProductImage, error) {
	plist, err := s.r.FetchAllProductsByUser(id)
	if err != nil {
		return nil, err
	}
	return plist, nil
}
