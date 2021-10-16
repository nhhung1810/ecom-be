package database

import (
	"database/sql"
	"ecom-be/app/product"
	"errors"

	"github.com/lib/pq"
)

func (s *Storage) AddProduct(p product.Product) error {
	var id int

	// PHASE ADD PRODUCT
	err := s.db.QueryRow(
		`INSERT INTO Products 
			(name, categories, brand, price, size, color, quantity, description)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		p.Name,
		p.Categories,
		p.Brand,
		p.Price,
		p.Size,
		p.Color,
		p.Quantity,
		p.Description,
	).Scan(&id)

	//PHASE ADD USER

	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return errors.New("add product error: " + errStr)
	}
	return nil
}

func (s *Storage) FetchProduct(id int) (*product.Product, error) {
	var p product.Product
	rows := s.db.QueryRow(
		`SELECT id, name, categories, brand, price, 
				size, color, quantity, description
		FROM Products WHERE id = $1`,
		id,
	)

	err := rows.Scan(
		&p.ID, &p.Name, &p.Categories,
		&p.Brand, &p.Price, &p.Color,
		&p.Quantity, &p.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errNotExistedProduct
		}
		return nil, errUnknown
	}
	return &p, nil
}
