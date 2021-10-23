package database

import (
	"ecom-be/app/order"
	"errors"

	"github.com/lib/pq"
)

func (s *Storage) AddOrder(ord []order.ProductOrder) error {
	sqlProd := `
	INSERT INTO ProductsOrder(userid, productid, quantity, price, color, size)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING orderid
	`
	for i := 0; i < len(ord); i++ {
		var orderid int
		err := s.db.QueryRow(sqlProd,
			ord[i].UserID,
			ord[i].ProductID,
			ord[i].Quantity,
			ord[i].Price,
			ord[i].Color,
			ord[i].Size,
		).Scan(&orderid)
		if err, ok := err.(*pq.Error); ok {
			errStr := "pq error:" + err.Code.Name()
			return errors.New("add ProductOrder error: " + errStr)
		}
	}
	return nil
}
func (s *Storage) FetchAllOrder() ([]order.ProductOrder, error) {
	return nil, nil
}
