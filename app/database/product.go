package database

import (
	"database/sql"
	"ecom-be/app/config"
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
		`SELECT id, name, p.categories, p.brand, p.price, 
				p.size, p.color, p.quantity, p.description, p.created_date::timestamp,
				p.quantity - COALESCE(sum(po.quantity), 0) as remain
		FROM Products as p
		LEFT JOIN ProductsOrder as po on po.productid = p.id
		WHERE p.id = $1
		GROUP BY id, name, categories, brand, p.price, P.size, p.color, 
		p.quantity, description, p.created_date::timestamp`,
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
		&p.Remain,
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

func (s *Storage) FetchAllProductsWithFilter(filter product.ProductFilter, sortIndex int) ([]product.ProductImage, error) {
	var plist []product.ProductImage
	sqlQuery := `
	SELECT 
		id, name, categories, brand, p.price, 
		p.size, p.color, p.quantity, description, p.created_date::timestamp,
		p.quantity - COALESCE(sum(po.quantity), 0) as remain
	FROM products as p
	LEFT JOIN productsorder as po on p.id = po.productid
	WHERE $1 <@ p.categories
		AND p.size && $2 
		AND p.color && $3
		AND (p.brand || ARRAY[]::varchar(50)[]) && $4
		AND p.price BETWEEN $5 AND $6
	GROUP BY id, name, categories, brand, p.price, 
		p.size, p.color, p.quantity, description, p.created_date::timestamp
		`

	// ADD FILTER TO QUERY
	if len(filter.IsAvailable) > 0 {
		if filter.IsAvailable[0] == "" {
			sqlQuery += ""
		} else if len(filter.IsAvailable) == 1 {
			if filter.IsAvailable[0] == "in" {
				sqlQuery += " HAVING p.quantity - COALESCE(sum(po.quantity), 0) <> 0"
			} else if filter.IsAvailable[0] == "out" {
				sqlQuery += " HAVING p.quantity - COALESCE(sum(po.quantity), 0) = 0"
			}
		} else if len(filter.IsAvailable) == 2 {
			sqlQuery += ""
		}
	}

	// ORDER THE QUERY
	sqlQuery += config.SortProductList[sortIndex]

	ctgsParam, sizeParam, colorsParam, brandsParam :=
		handleNullArray(filter.Categories, filter.Size, filter.Color, filter.Brand)

	priceStart, err := strconv.Atoi(filter.PriceStart)
	if err != nil {
		return nil, err
	}

	priceStop, err := strconv.Atoi(filter.PriceStop)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(sqlQuery,
		ctgsParam,
		sizeParam,
		colorsParam,
		brandsParam,
		priceStart,
		priceStop,
	)
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
			&p.Remain,
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

func (s *Storage) CountAllProductBySellerID(id int) (*int, error) {
	var count int
	sqlQuery := `
	SELECT 
	count(p.id) as count_id
	FROM PRODUCTS AS p
	JOIN ProductUser as pu on pu.productid = p.id
	WHERE pu.userid = $1`

	err := s.db.QueryRow(sqlQuery, id).Scan(&count)
	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		return nil, errors.New("count product by seller id error: " + errStr)
	}

	return &count, nil
}

func (s *Storage) FetchAllProductsWithOrderInfo(userid int, limit int, offset int) ([]product.ProductWithOrderInfo, error) {
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

	rows, err := s.db.Query(sqlQuery, userid, limit, offset)
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
			if err == sql.ErrNoRows {
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

func (s *Storage) SearchProductByName(name string) (*product.SearchProduct, error) {
	var result product.SearchProduct

	sqlQuery := `
	SELECT p.id, p.categories
	FROM Products as p
	WHERE UPPER(p.name) like '%' || $1 || '%'
	LIMIT 1
	`
	err := s.db.QueryRow(sqlQuery, name).
		Scan(&result.ID, pq.Array(&result.Categories))

	println(result.ID)
	if err, ok := err.(*pq.Error); ok {
		errStr := "pq error:" + err.Code.Name()
		println(errStr)
		return nil, errors.New("search product error: " + errStr)
	}

	return &result, nil
}
