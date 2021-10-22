package database

import (
	"ecom-be/app/order"
	"errors"

	"github.com/lib/pq"
)

func (s *Storage) AddOrder(ord *order.Order) error {
	var orderid int
	orderQuery := `
	INSERT INTO Orders (userid) VALUES ($1) RETURNING id
	`
	err := s.db.QueryRow(orderQuery, ord.UserID).Scan(&orderid)
	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return errors.New("add order error: " + errStr)
	}

	// INSERT PRODUCT_ORDER
	sqlProd := `
	INSERT INTO ProductsOrder (orderid, productid, quantity, price, color, size)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING orderid
	`
	for i := 0; i < len(ord.Prods); i++ {
		err = s.db.QueryRow(sqlProd,
			orderid,
			ord.Prods[i].ProductID,
			ord.Prods[i].Quantity,
			ord.Prods[i].Price,
			ord.Prods[i].Color,
			ord.Prods[i].Size,
		).Scan(&orderid)
		if err, ok := err.(*pq.Error); ok {
			errStr := "pq error:" + err.Code.Name()
			return errors.New("add ProductOrder error: " + errStr)
		}
	}
	return nil
}
func (s *Storage) FetchAllOrder() (*order.Order, error) {
	return nil, nil
}
