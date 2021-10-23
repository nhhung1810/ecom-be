package product

import (
	"time"
)

type Product struct {
	ID          int       `form:"id" json:"id"`
	Name        string    `json:"name"`
	Categories  []string  `json:"categories"`
	Brand       string    `json:"brand"`
	Price       float32   `json:"price"`
	Size        []string  `json:"size"`
	Color       []string  `json:"colors"`
	Quantity    int       `json:"quantity"`
	Description string    `json:"description"`
	CreatedDate time.Time `json:"created_date"`
}

type ProductImage struct {
	Prod       Product
	ImageCount int
}
