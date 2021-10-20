package database

import (
	"database/sql"
	"ecom-be/app/product"
	"errors"
	"io/ioutil"
	"strconv"

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

func (s *Storage) FetchProductByID(id int) (*product.ProductImage, error) {
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
		&p.Size,
		&p.Color,
		&p.Quantity,
		&p.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	//IMAGE FOLDER FETCH
	productid := strconv.FormatInt(int64(p.ID), 10)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir("./images/" + productid + "/")
	if err != nil {
		return nil, err
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() {
			count++
		}
	}

	return &product.ProductImage{
		Prod:       p,
		ImageCount: count,
	}, nil
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
		// FETCH PRODUCT
		var p product.Product
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

		//IMAGE FOLDER FETCH
		productid := strconv.FormatInt(int64(p.ID), 10)
		if err != nil {
			return nil, err
		}

		files, err := ioutil.ReadDir("./images/" + productid + "/")
		if err != nil {
			return nil, err
		}

		count := 0
		for _, file := range files {
			if !file.IsDir() {
				count++
			}
		}

		plist = append(plist, product.ProductImage{
			Prod:       p,
			ImageCount: count,
		})
	}

	return plist, nil
}
