package order

import "time"

type ProductOrder struct {
	OrderID     int       `json:"orderid"`
	UserID      int       `json:"-"`
	ProductID   int       `json:"productid"`
	Quantity    int       `json:"quantity"`
	Price       float32   `json:"price"`
	Color       string    `json:"color"`
	Size        string    `json:"size"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
}

type OrderBySeller struct {
	OrderID     int       `json:"orderid"`
	Name        string    `json:"name"`
	UserID      int       `json:"-"`
	ProductID   int       `json:"productid"`
	Quantity    int       `json:"quantity"`
	Price       float32   `json:"price"`
	Color       string    `json:"color"`
	Size        string    `json:"size"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
}
