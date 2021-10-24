package product

import (
	"github.com/gin-gonic/gin"
)

//provide function access to db
type Repository interface {
	AddProduct(p Product, userid int) (*int, error)
	FetchProductByID(id int) (*ProductImage, error)
	FetchAllProductsByUser(id int) ([]ProductImage, error)
	FetchAllProductsByCtg(ctg []string) ([]ProductImage, error)
	FetchAllProductsWithFilter(ctgs []string,
		sizes []string, colors []string) ([]ProductImage, error)
	FetchAllProductsWithOrderInfo(userid int) ([]ProductWithOrderInfo, error)
}

// Provide interface for product operation in handler
type Service interface {
	AddProduct(p Product, userid int) (*int, error)
	FetchProductByID(id int) (*ProductImage, error)
	ParseProduct(g *gin.Context) (*Product, error)
	FetchAllProductsByUser(id int) ([]ProductImage, error)
	FetchAllProductsByCtg(ctg []string) ([]ProductImage, error)
	FetchAllProductsWithFilter(ctgs []string,
		sizes []string, colors []string) ([]ProductImage, error)
	FetchAllProductsWithOrderInfo(userid int) ([]ProductWithOrderInfo, error)
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
		println(err.Error())
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

func (s *service) FetchProductByID(id int) (*ProductImage, error) {
	p, err := s.r.FetchProductByID(id)
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

func (s *service) FetchAllProductsByCtg(ctg []string) ([]ProductImage, error) {
	s.r.FetchAllProductsByCtg(ctg)
	return nil, nil
}

func (s *service) FetchAllProductsWithFilter(ctgs []string,
	sizes []string, colors []string) ([]ProductImage, error) {
	p, err := s.r.FetchAllProductsWithFilter(ctgs, sizes, colors)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) FetchAllProductsWithOrderInfo(userid int) ([]ProductWithOrderInfo, error) {
	p, err := s.r.FetchAllProductsWithOrderInfo(userid)
	if err != nil {
		return nil, err
	}
	return p, nil
}
