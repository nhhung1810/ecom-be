package database

import (
	"database/sql"
	"ecom-be/app/product"
	"errors"

	"github.com/lib/pq"
)

func (s *Storage) AddProduct(p product.Product, userid int) (*int, error) {
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

	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return nil, errors.New("add product error: " + errStr)
	}

	//PHASE ADD USER
	err = s.db.QueryRow(
		`INSERT INTO ProductUser(productid, userid)
		VALUES($1, $2) RETURNING productid`,
		id,
		userid,
	).Scan(&id)

	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return nil, errors.New("add product-user error: " + errStr)
	}

	return &id, nil
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
		&p.ID,
		&p.Name,
		&p.Categories,
		&p.Brand,
		&p.Price,
		&p.Color,
		&p.Quantity,
		&p.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errNotExistedProduct
		}
		return nil, errUnknown
	}
	return &p, nil
}

func (s *Storage) FetchAllProductsByUser(id int) ([]product.ProductImage, error) {
	var plist []product.ProductImage
	rows, err := s.db.Query(
		`SELECT id, name, categories, brand, price, 
				size, color, quantity, description
		FROM Products as p
		JOIN ProductUser AS pu ON p.id = pu.productid
		WHERE pu.userid = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p product.Product
		var imageid []string

		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Categories,
			&p.Brand,
			&p.Price,
			&p.Size,
			&p.Color,
			&p.Quantity,
			&p.Description,
		)

		if err != nil {
			println(err.Error())
			if err == sql.ErrNoRows {
				println(err.Error())
				return nil, errNotExistedProduct
			}
			return nil, errUnknown
		}

		//IMAGE ID RETRIEVAL
		imgrows, err := s.db.Query(
			`SELECT imageid FROM ProductImages
			WHERE productid = $1`,
			p.ID,
		)
		if err != nil {
			return nil, err
		}
		defer imgrows.Close()
		for imgrows.Next() {
			var tmp string
			err = imgrows.Scan(&tmp)
			if err != nil {
				println(err.Error())
				if err == sql.ErrNoRows {
					println(err.Error())
					return nil, errNotExistedProduct
				}
				return nil, errUnknown
			}
			imageid = append(imageid, tmp)
		}

		plist = append(plist, product.ProductImage{
			Prod:    p,
			Imageid: imageid,
		})
	}

	return plist, nil
}
