package order

type ProductOrder struct {
	OrderID   int     `json:"-"`
	UserID    int     `json:"-"`
	ProductID int     `json:"productid"`
	Quantity  int     `json:"quantity"`
	Price     float32 `json:"price"`
	Color     string  `json:"color"`
	Size      string  `json:"size"`
	Status    string  `json:"-"`
}
