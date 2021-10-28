package database

import (
	"database/sql"
	"ecom-be/app/config"
	"ecom-be/app/product"
	"errors"
	"fmt"
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
				size, color, quantity, description, created_date::timestamp
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
		&p.CreatedDate,
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
				size, color, quantity, description, created_date::timestamp
		FROM Products as p
		JOIN ProductUser AS pu ON p.id = pu.productid
		WHERE pu.userid = $1`,
		id,
	)

	if err != nil {
		println(err.Error())
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
			&p.CreatedDate,
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

func (s *Storage) FetchAllProductsWithFilter(filter product.ProductFilter) ([]product.ProductImage, error) {
	var plist []product.ProductImage
	sqlQuery := `
		SELECT 
			id, name, categories, brand, price, 
			size, color, quantity, description, created_date::timestamp
		FROM products as p
		WHERE $1 <@ p.categories
		AND p.size && $2 
		AND p.color && $3
		AND (p.brand || ARRAY[]::varchar(50)[]) && $4
		`
	// THE OR-AND IS USE TO CANCEL THE
	// EFFECT WHEN THERE ARE NO ELEMENT IN OR-SELECTOR
	// IT WILL TAKE THE EFFECT OF THE AND-SELECTOR
	// (AS THE AND-SELECTOR IS SMALLER THAN THE
	// OR-SELECTOR, BUT IT WON'T RETURN 0 RESULT WHEN
	// THE ARRAY IS EMPTY )
	ctgsParam, sizeParam, colorsParam, brandsParam :=
		handleNullArray(filter.Categories, filter.Size, filter.Color, filter.Brand)

	fmt.Printf("brandsParam: %v\n", brandsParam)
	rows, err := s.db.Query(sqlQuery, ctgsParam, sizeParam, colorsParam, brandsParam)
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
			&p.CreatedDate,
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

func (s *Storage) FetchAllProductsWithOrderInfo(userid int) ([]product.ProductWithOrderInfo, error) {
	// TODO: ADD PAGING
	var plist []product.ProductWithOrderInfo
	sqlQuery :=
		`SELECT 
		p.id, p.name, p.price, p.quantity as capacity, 
		p.categories, p.created_date::timestamp, 
		COALESCE(sum(ps.quantity), 0) as sold
	FROM PRODUCTS AS p
	LEFT JOIN ProductsOrder as ps ON p.id = ps.productid
	JOIN ProductUser as pu on pu.productid = p.id
	WHERE pu.userid = $1
	GROUP BY p.id, p.name, p.quantity, p.categories, p.created_date
	ORDER BY p.created_date
	LIMIT $2
	OFFSET $3`

	rows, err := s.db.Query(sqlQuery, userid, 5, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p product.ProductWithOrderInfo
		rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Capacity,
			pq.Array(&p.Categories),
			&p.CreatedDate,
			&p.Sold,
		)

		if err != nil {
			println(err.Error())
			if err == sql.ErrNoRows {
				println(err.Error())
				return nil, errNotExistedProduct
			}
			return nil, errUnknown
		}
		plist = append(plist, p)
	}
	return plist, nil
}

// PARAM HANDLE
// PQ DP NOT HANDLE NIL VALUE OF THE STRING
// SO I HAVE TO SELECT MANUALLY
func handleNullArray(
	ctgs []string,
	sizes []string,
	colors []string,
	brands []string,
) (interface{}, interface{}, interface{}, interface{}) {
	var ctgsParam, sizeParam, colorsParam, brandsParam interface{}

	if (len(ctgs) == 1 && ctgs[0] == "") || (len(ctgs) == 0) {
		ctgsParam = "{}"
	} else {
		ctgsParam = pq.Array(ctgs)
	}

	if (len(sizes) == 1 && sizes[0] == "") || (len(sizes) == 0) {
		sizeParam = pq.Array(config.SizeArray)
	} else {
		sizeParam = pq.Array(sizes)
	}

	if (len(colors) == 1 && colors[0] == "") || (len(colors) == 0) {
		colorsParam = pq.Array(config.ColorArray)
	} else {
		colorsParam = pq.Array(colors)
	}

	if (len(brands) == 1 && brands[0] == "") || (len(brands) == 0) {
		brandsParam = pq.Array(config.BrandArray)
	} else {
		brandsParam = pq.Array(brands)
	}

	return ctgsParam, sizeParam, colorsParam, brandsParam
}
