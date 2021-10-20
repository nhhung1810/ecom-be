package image

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

type Repository interface {
}

type Service interface {
	UploadImage(images *Image) error
	GetImage(index string, productid string) (*bytes.Buffer, error)
	ParseImages(c *gin.Context) (*Image, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetImage(index string, productid string) (*bytes.Buffer, error) {
	f, err := os.Open("./images/" + productid +
		"/" + index + ".jpg")
	if err != nil {
		println(err.Error())
		return nil, err
	}
	defer f.Close()
	serverImg, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, serverImg, nil)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (s *service) UploadImage(image *Image) error {
	// TODO: VALIDATE
	tmp := image.Data

	// Create the folder
	err := os.Mkdir("./images/"+image.ProductID, os.ModeAppend.Perm())
	if err != nil {
		println(err.Error())
	}

	for i, _ := range tmp {
		f, err := tmp[i].Open()
		if err != nil {
			println(err.Error())
			return err
		}
		defer f.Close()

		// Create the image in the created folder
		out, err := os.Create("./images/" + image.ProductID +
			"/" + tmp[i].Filename + ".jpg")
		if err != nil {
			println(err.Error())
			return err
		}
		defer out.Close()

		// Populate the file into placeholder
		_, err = io.Copy(out, f)
		if err != nil {
			println(err.Error())
			return err
		}

		println("success" + tmp[i].Filename)
	}

	// No need to handle the link with product
	// as the product id now become the
	// image folder id too

	return nil
}

func (s *service) ParseImages(c *gin.Context) (*Image, error) {
	c.Request.ParseMultipartForm(10 << 20)
	form, err := c.MultipartForm()
	if err != nil {
		println(err.Error())
		return nil, err
	}
	f := form.File["images[]"]
	tmp := form.Value["productid"]
	if len(tmp) > 1 {
		return nil, errors.New("error: there are multiple product id")
	}

	return &Image{
		ProductID: tmp[0],
		Data:      f,
	}, nil
}
