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

type ProductWithOrderInfo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Categories  []string  `json:"categories"`
	Price       float32   `json:"price"`
	Sold        int       `json:"sold"`
	Capacity    int       `json:"capacity"`
	CreatedDate time.Time `json:"created_date"`
}

type ProductFilter struct {
	Categories  []string `form:"categories" binding:"required"`
	Size        []string `form:"size"`
	Color       []string `form:"colors"`
	Brand       []string `form:"brands"`
	PriceStart  string   `form:"pstart"`
	PriceStop   string   `form:"pstop"`
	IsAvailable []string `form:"available"`
}
