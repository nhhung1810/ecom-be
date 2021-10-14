package handle

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var errInternal = gin.H{
	"message": "internal",
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
