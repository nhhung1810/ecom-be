package handle

import (
	"ecom-be/app/config"
	"ecom-be/app/order"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errUnauthorized = gin.H{
	"message": "Unauthorized",
}

var errNotFound = gin.H{
	"message": "not found",
}

func uploadOrder(orderService order.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusUnauthorized, errUnauthorized)
			return
		}
		println("userid", userid)
		ord, err := orderService.ParseOrder(c, *userid)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, errBadResquest)
			return
		}

		err = orderService.AddOrder(ord)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, errInternal)
			return
		}

		c.JSON(http.StatusOK, successMsg)
	}
}

func getAllOrderByProductID(orderService order.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		type queryHandle struct {
			ProductID int `form:"productid" binding:"required"`
		}

		var qh queryHandle
		err := c.ShouldBindQuery(&qh)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, errBadResquest)
			return
		}

		orders, err := orderService.FetchAllOrderByProductID(qh.ProductID)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusNotFound, errNotFound)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    orders,
		})
	}
}

func getAllOrderBySellerID(orderService order.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errUnauthorized)
			return
		}

		type Paging struct {
			Limit  int `form:"limit"`
			Offset int `form:"offset"`
		}

		var paging Paging
		err = c.ShouldBindQuery(&paging)
		if err != nil {
			paging.Limit = 10
			paging.Offset = 0
		}

		orders, err := orderService.FetchAllOrderBySellerID(*userid, paging.Limit, paging.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errInternal)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    orders,
		})
	}
}

func countAllOrderbySellerID(orderService order.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errUnauthorized)
			return
		}

		count, err := orderService.CountAllOrderbySellerID(*userid)
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

func updateOrderStatus(orderService order.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		_, err := cookieAuth(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errUnauthorized)
			return
		}

		type QueryHandle struct {
			OrderId    int `form:"orderid"`
			StatusCode int `form:"status"`
		}

		var qh QueryHandle
		err = c.ShouldBindQuery(&qh)
		if err != nil {
			c.JSON(http.StatusBadRequest, errBadResquest)
			return
		}
		println(qh.OrderId)
		println(qh.StatusCode)
		err = orderService.UpdateStatusByOrder(qh.OrderId,
			config.OrderStatus[qh.StatusCode%3])
		if err != nil {
			c.JSON(http.StatusNotFound, errNotFound)
			return
		}

		c.JSON(http.StatusOK, successMsg)
	}
}
