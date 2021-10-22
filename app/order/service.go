package order

import (
	"github.com/gin-gonic/gin"
)

type Repository interface {
	AddOrder(order *Order) error
	FetchAllOrder() (*Order, error)
}

type Service interface {
	AddOrder(order *Order) error
	FetchAllOrder() (*Order, error)
	ParseOrder(c *gin.Context, userid int) (*Order, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddOrder(order *Order) error {
	err := s.r.AddOrder(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FetchAllOrder() (*Order, error) {
	return nil, nil
}

func (s *service) ParseOrder(c *gin.Context, userid int) (*Order, error) {
	var ord Order
	err := c.Bind(&ord.Prods)
	if err != nil {
		return nil, err
	}
	ord.UserID = userid
	return &ord, nil
}
