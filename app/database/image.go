package database

import (
	"database/sql"
	image "ecom-be/app/image"
	"errors"

	"github.com/lib/pq"
)

func (s *Storage) UploadImage(images *[]image.Image) ([]string, error) {
	// upload won't change the original Data URLS
	// therefore, we must process the 'data' field before
	// any decoding

	var idlist []string

	for i := 0; i < len(*images); i++ {
		// TODO: Check for the existed product

		// Phase 1: add image
		img := (*images)[i]
		var id string
		err := s.db.QueryRow(
			`INSERT INTO Images(id, dat)
			VALUES ($1, $2) RETURNING id`,
			img.ID,
			img.Data,
		).Scan(&id)

		if err, ok := err.(*pq.Error); ok {
			// TODO: Consider change error handling to
			// not cancel all but continue to upload
			// other images
			errStr := "pq error:" + err.Code.Name()
			return nil, errors.New("add image error: " + errStr)
		}

		// Phase 2: binding to product
		productID := img.ProductID
		err = s.db.QueryRow(
			`INSERT INTO ProductImages(productid, imageid)
			VALUES ($1, $2) RETURNING imageid`,
			productID,
			img.ID,
		).Scan(&id)

		if err, ok := err.(*pq.Error); ok {
			errStr := "pq error:" + err.Code.Name()
			return nil, errors.New("add product-image error: " + errStr)
		}

		idlist = append(idlist, id)
	}

	// Use the id to temporary testing
	// can be remove later
	return idlist, nil
}

func (s *Storage) GetImage(id string) (*image.Image, error) {
	var img image.Image
	rows := s.db.QueryRow(
		`SELECT id, dat, pi.productid
		FROM Images as img
		JOIN ProductImages as pi on img.id = pi.imageid 
		WHERE id = $1`,
		id,
	)

	err := rows.Scan(&img.ID, &img.Data, &img.ProductID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errNotExistedAccount
		}
		return nil, errUnknown
	}
	return &img, nil
}
