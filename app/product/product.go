package product

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Categories  string  `json:"categories"`
	Brand       string  `json:"brand"`
	Price       float32 `json:"price"`
	Size        string  `json:"size"`
	Color       string  `json:"color"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
}
