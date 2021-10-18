package image

type Image struct {
	ID        string `json:id`
	ProductID int    `json:productId` // for making the ProductTable
	Data      string `json:data`
}
