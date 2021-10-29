package handle

import (
	"bytes"
	myimg "ecom-be/app/image"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var errInternal = gin.H{
	"message": "internal",
}

var errBadResquest = gin.H{
	"message": "bad request",
}

// BASE ONLY QUERY
func getImageHandle(img myimg.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var queryImg myimg.Image
		err := c.BindQuery(&queryImg)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"message": "check your query",
			})
			return
		}

		buffer, err := img.GetImage(queryImg.Index, queryImg.ProductID)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
		}

		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		c.Header("Cache-Control", "public, max-age=15552000")
		if _, err := c.Writer.Write(buffer.Bytes()); err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
			return
		}
	}
}

// ONLY FOR TESTING
func imageHandle() func(c *gin.Context) {
	return func(c *gin.Context) {
		// TEST IMAGE
		// TODO: CHANGE TO DATABASE
		f, err := os.Open("./blue.jpg")
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
			return
		}

		defer f.Close()
		img, _, err := image.Decode(f)

		buffer := new(bytes.Buffer)
		err = jpeg.Encode(buffer, img, nil)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
		}

		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
			return
		}
		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		c.Header("Cache-Control", "public, max-age=15552000")
		if _, err := c.Writer.Write(buffer.Bytes()); err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
		}

		c.JSON(http.StatusOK, successMsg)
	}
}

// NEW UPLOAD FUNCTION WITH MULTIPART/FORM-DATA REQUEST
func newUploadImage(img myimg.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		// call the interface to handle parsing
		images, err := img.ParseImages(c)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, errBadResquest)
		}

		// call the interface to handle new images
		img.UploadImage(images)

		c.JSON(http.StatusCreated, successMsg)
	}
}
