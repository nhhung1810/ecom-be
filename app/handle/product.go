package handle

import (
	"ecom-be/app/product"
	"net/http"

	"github.com/gin-gonic/gin"
)

func uploadProduct(pr product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		println(*userid)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
		}

		// Phase insert product
		prod, err := pr.ParseProduct(c)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		id, err := pr.AddProduct(*prod, *userid)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		println(*id)
		c.IndentedJSON(http.StatusCreated, gin.H{
			"message": "success",
			"id":      *id,
		})
	}
}

func getAllProducts(pr product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		println(*userid)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
		}

		plist, err := pr.FetchAllProductsByUser(*userid)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		}

		c.IndentedJSON(http.StatusAccepted, gin.H{
			"message":  "success",
			"products": plist,
		})

	}
}
