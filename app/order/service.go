package order

import (
	"github.com/gin-gonic/gin"
)

type Repository interface {
	AddOrder(order []ProductOrder) error
	FetchAllOrder() ([]ProductOrder, error)
}

type Service interface {
	AddOrder(order []ProductOrder) error
	FetchAllOrder() ([]ProductOrder, error)
	ParseOrder(c *gin.Context, userid int) ([]ProductOrder, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddOrder(order []ProductOrder) error {
	err := s.r.AddOrder(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FetchAllOrder() ([]ProductOrder, error) {
	return nil, nil
}

func (s *service) ParseOrder(c *gin.Context, userid int) ([]ProductOrder, error) {
	var ord []ProductOrder
	err := c.Bind(&ord)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ord); i++ {
		ord[i].UserID = userid
	}
	return ord, nil
}
