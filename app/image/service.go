package image

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Repository interface {
	UploadImage(images *[]Image) ([]string, error)
	GetImage(id string) (*Image, error)
}

type Service interface {
	UploadImage(images *[]Image) error
	GetImage(id string) (*Image, error)
	ParseImages(c *gin.Context) (*[]Image, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) UploadImage(images *[]Image) error {
	// TODO: VALIDATE
	_, err := s.r.UploadImage(images)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetImage(id string) (*Image, error) {
	img, err := s.r.GetImage(id)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (s *service) ParseImages(c *gin.Context) (*[]Image, error) {
	var images []Image
	err := c.BindJSON(&images)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(images); i++ {
		// RANDOM ID
		images[i].ID = xid.New().String()
	}
	return &images, nil
}
