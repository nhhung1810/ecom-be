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
	// PRE-PROCESS ARRAY STRING

	// PHASE ADD PRODUCT
	err := s.db.QueryRow(
		`INSERT INTO Products 
			(name, categories, brand, price, size, color, quantity, description)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		p.Name,
		pq.Array(p.Categories),
		p.Brand,
		p.Price,
		pq.Array(p.Size),
		pq.Array(p.Color),
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
		pq.Array(&p.Categories),
		&p.Brand,
		&p.Price,
		pq.Array(&p.Size),
		pq.Array(&p.Color),
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
			pq.Array(&p.Categories),
			&p.Brand,
			&p.Price,
			pq.Array(&p.Size),
			pq.Array(&p.Color),
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

func (s *Storage) FetchAllProductsByCtg(ctgs []string) ([]product.ProductImage, error) {
	return nil, nil
}

func (s *Storage) FetchAllProductsWithFilter(ctgs []string, sizes []string, colors []string) ([]product.ProductImage, error) {
	var plist []product.ProductImage
	println("LENs OF ", len(ctgs), len(sizes), len(colors))
	println(ctgs[0], sizes[0], colors[0])
	sqlQuery := `
		SELECT 
			id, name, categories, brand, price, 
			size, color, quantity, description
		FROM products as p
		WHERE $1 <@ p.categories
		AND $2 <@ p.size
		AND $3 <@ p.color`
	ctgsParam, sizeParam, colorsParam := handleNullArray(ctgs, sizes, colors)

	rows, err := s.db.Query(sqlQuery, ctgsParam, sizeParam, colorsParam)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p product.Product
		err = rows.Scan(
			&p.ID,
			&p.Name,
			pq.Array(&p.Categories),
			&p.Brand,
			&p.Price,
			pq.Array(&p.Size),
			pq.Array(&p.Color),
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

// PARAM HANDLE
// PQ DP NOT HANDLE NIL VALUE OF THE STRING
// SO I HAVE TO SELECT MANUALLY
func handleNullArray(ctgs []string,
	sizes []string, colors []string) (interface{}, interface{}, interface{}) {
	var ctgsParam, sizeParam, colorsParam interface{}
	if len(ctgs) > 0 {
		ctgsParam = pq.Array(ctgs)
	} else {
		ctgsParam = "{}"
	}
	if sizes[0] != "" {
		sizeParam = pq.Array(sizes)
	} else {
		sizeParam = "{}"
	}
	if colors[0] != "" {
		colorsParam = pq.Array(colors)
	} else {
		colorsParam = "{}"
	}
	return ctgsParam, sizeParam, colorsParam
}
