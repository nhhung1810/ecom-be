package handle

import (
	"ecom-be/app/product"
	"net/http"

	"github.com/gin-gonic/gin"
)

func uploadProduct(productService product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		println(*userid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
		}

		// Phase insert product
		prod, err := productService.ParseProduct(c)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		id, err := productService.AddProduct(*prod, *userid)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		println(*id)
		c.JSON(http.StatusCreated, gin.H{
			"message": "success",
			"id":      *id,
		})
	}
}

func getAllProducts(productService product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		println(*userid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
		}

		plist, err := productService.FetchAllProductsByUser(*userid)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "success",
			"products": plist,
		})

	}
}

func getProductByID(productService product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var queryProd product.Product
		err := c.BindQuery(&queryProd)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, errBadResquest)
			return
		}

		p, err := productService.FetchProductByID(queryProd.ID)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    p,
		})
	}
}

func getProductWithFilter(productService product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var filter product.ProductFilter
		type SortIndex struct {
			Index int `form:"sort"`
		}
		var sortIndex SortIndex
		err := c.ShouldBindQuery(&filter)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, errBadResquest)
			return
		}

		err = c.ShouldBindQuery(&sortIndex)
		if err != nil {
			println(err.Error())
			// sortIndex.index = 0
		}
		sortIndex.Index %= 4
		println(sortIndex.Index)

		p, err := productService.FetchAllProductsWithFilter(filter,
			sortIndex.Index)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    p,
		})
	}
}

func getProductWithOrder(productService product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		type Paging struct {
			Limit  int `form:"limit"`
			Offset int `form:"offset"`
		}

		var paging Paging
		err = c.ShouldBindQuery(&paging)
		if err != nil {
			paging.Limit = 5
			paging.Offset = 0
		}

		plist, err := productService.FetchAllProductsWithOrderInfo(*userid, paging.Limit, paging.Offset)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    plist,
		})

	}
}

func countAllProductBySellerID(productService product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errUnauthorized)
			return
		}

		count, err := productService.CountAllProductBySellerID(*userid)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"count":   &count,
		})
	}
}

func searchProductByName(productService product.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		type QueryHandle struct {
			Name string `form:"name"`
		}

		var qh QueryHandle
		err := c.ShouldBindQuery(&qh)
		if err != nil {
			c.JSON(http.StatusBadRequest, errBadResquest)
			return
		}
		println(qh.Name)

		p, err := productService.SearchProductByName(qh.Name)
		if err != nil {
			c.JSON(http.StatusNotFound, errNotFound)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    p,
		})
	}
}
