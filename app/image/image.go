package image

import "mime/multipart"

type Image struct {
	// Don't use tag as this
	// will be parsed manually
	ProductID string                  `form:"productid" json:"-"`
	Index     string                  `form:"id" json:"-"`
	Data      []*multipart.FileHeader `json:"-"`
}
