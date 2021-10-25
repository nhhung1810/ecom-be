package database

import (
	"database/sql"
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

func (s *Storage) FetchAllOrderByProductID(id int) ([]order.ProductOrder, error) {
	var olist []order.ProductOrder
	sqlQuery := `
	SELECT 
		ps.orderid, ps.userid, ps.status, 
		ps.productid, ps.quantity, ps.price, 
		ps.color, ps.size, ps.created_date::timestamp
	FROM Productsorder as ps
	WHERE ps.productid = $1
	`
	rows, err := s.db.Query(sqlQuery, id)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var ord order.ProductOrder
		err = rows.Scan(
			&ord.OrderID,
			&ord.UserID,
			&ord.Status,
			&ord.ProductID,
			&ord.Quantity,
			&ord.Price,
			&ord.Color,
			&ord.Size,
			&ord.CreatedDate,
		)

		if err != nil {
			println(err.Error())
			if err == sql.ErrNoRows {
				println(err.Error())
				return nil, errNotExistedProduct
			}
			return nil, errUnknown
		}
		olist = append(olist, ord)
	}

	return olist, nil
}

func (s *Storage) FetchAllOrderBySellerID(id int, limit int, offset int) ([]order.OrderBySeller, error) {
	var olist []order.OrderBySeller
	sqlQuery := `
	SELECT 
		ps.orderid, p.name, ps.userid, ps.status, 
		ps.productid, ps.quantity, ps.price, 
		ps.color, ps.size, ps.created_date::timestamp
	FROM Productsorder as ps
	JOIN ProductUser as pu on ps.productid = pu.productid
	JOIN Products as p on p.id = ps.productid
	WHERE pu.userid = $1
	ORDER BY ps.orderid asc
	LIMIT $2
	OFFSET $3
	`
	rows, err := s.db.Query(sqlQuery, id, limit, offset)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var ord order.OrderBySeller
		err = rows.Scan(
			&ord.OrderID,
			&ord.Name,
			&ord.UserID,
			&ord.Status,
			&ord.ProductID,
			&ord.Quantity,
			&ord.Price,
			&ord.Color,
			&ord.Size,
			&ord.CreatedDate,
		)

		if err != nil {
			println(err.Error())
			if err == sql.ErrNoRows {
				println(err.Error())
				return nil, errNotExistedProduct
			}
			return nil, errUnknown
		}
		olist = append(olist, ord)
	}

	return olist, nil
}

func (s *Storage) CountAllOrderbySellerID(id int) (*int, error) {
	var count int
	sqlQuery := `
	SELECT 
		COUNT(ps.orderid) AS count_id
	FROM Productsorder as ps
	JOIN ProductUser as pu on ps.productid = pu.productid
	JOIN Products as p on p.id = ps.productid
	WHERE pu.userid = $1
	`
	err := s.db.QueryRow(sqlQuery, id).Scan(&count)
	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return nil, errors.New("count order by seller id error: " + errStr)
	}

	return &count, nil
}

func (s *Storage) UpdateStatusByOrder(id int, status string) error {
	println(status)
	println(id)
	var orderid int
	sqlQuery := `
	UPDATE ProductsOrder
	SET status = $1
	WHERE orderid = $2 
	RETURNING orderid
	`

	err := s.db.QueryRow(sqlQuery, status, id).Scan(&orderid)
	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return errors.New("update order status error: " + errStr)
	}
	return nil
}
