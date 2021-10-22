package handle

import (
	"ecom-be/app/order"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errUnauthorized = gin.H{
	"message": "Unauthorized",
}

func uploadOrder(or order.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userid, err := cookieAuth(c)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusUnauthorized, errUnauthorized)
			return
		}
		println("userid", userid)
		ord, err := or.ParseOrder(c, *userid)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusBadRequest, errBadResquest)
			return
		}

		err = or.AddOrder(ord)
		if err != nil {
			println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, errInternal)
			return
		}

		c.IndentedJSON(http.StatusAccepted, successMsg)
	}
}
