package handle

import (
	"bytes"
	myimg "ecom-be/app/image"
	"encoding/base64"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var errInternal = gin.H{
	"message": "internal",
}

func uploadImageHandle(img myimg.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		imgs, err := img.ParseImages(c)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		err = img.UploadImage(imgs)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.IndentedJSON(http.StatusAccepted, successMsg)
	}
}

func getImageHandle(img myimg.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		// no need to catch the error as the server
		// only match this route when there are image/:id
		id := c.Param("id")
		img, err := img.GetImage(id)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		//process the data field
		data := img.Data
		s := strings.Split(data, ",")
		if len(s) <= 1 {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": ("data url error: " + data),
			})
			return
		}

		data = s[1] //get the base64 string
		decodedData, err := base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		imageReader := bytes.NewReader(decodedData)
		serveImg, _, err := image.Decode(imageReader)

		buffer := new(bytes.Buffer)
		err = jpeg.Encode(buffer, serveImg, nil)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
		}

		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
			return
		}
		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		c.Header("Cache-Control", "public, max-age=15552000")
		if _, err := c.Writer.Write(buffer.Bytes()); err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
		}
	}
}

// TODO: ADD SERVICE
func imageHandle() func(c *gin.Context) {
	return func(c *gin.Context) {
		// TEST IMAGE
		// TODO: CHANGE TO DATABASE
		f, err := os.Open("./blue.jpg")
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
			return
		}

		defer f.Close()
		img, _, err := image.Decode(f)

		buffer := new(bytes.Buffer)
		err = jpeg.Encode(buffer, img, nil)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
		}

		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
			return
		}
		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		c.Header("Cache-Control", "public, max-age=15552000")
		if _, err := c.Writer.Write(buffer.Bytes()); err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
		}
	}
}
