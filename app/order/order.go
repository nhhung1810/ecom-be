package order

type Order struct {
	ID     int
	UserID int
	Status int
	Prods  []ProductOrder
}

type ProductOrder struct {
	ProductID int     `json:"productid"`
	Quantity  int     `json:"quantity"`
	Price     float32 `json:"price"`
	Color     string  `json:"color"`
	Size      string  `json:"size"`
}
