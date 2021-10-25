package handle

import (
	"ecom-be/app/auth"
	"ecom-be/app/config"
	myimg "ecom-be/app/image"
	"ecom-be/app/order"
	"ecom-be/app/product"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler(auth auth.Service, img myimg.Service, pr product.Service, or order.Service) (*gin.Engine, error) {
	router := gin.Default()

	corConfig := cors.DefaultConfig()
	corConfig.AllowOrigins = config.DefaultConfig.AllowOrigin
	corConfig.AllowCredentials = true
	router.Use(cors.New(corConfig))

	//TEST ROUTING
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// AUTH HANDLE
	router.POST("/register", registerHandle(auth))
	router.POST("/login", loginHandle(auth))
	router.POST("/logout", logoutHandle(auth))
	router.GET("/user", userHandle(auth))

	// IMAGE HANDLE
	// router.GET("/image", imageHandle())
	router.GET("/image", getImageHandle(img))
	router.POST("/image/upload", newUploadImage(img))

	// PROD HANDLE
	router.GET("/product", getAllProducts(pr))
	router.GET("/product/q", getProductWithFilter(pr))
	// ONLY FOR QUERY USE
	router.GET("/product/info", getProductByID(pr))
	router.POST("/product/upload", uploadProduct(pr))

	// ORDER
	router.POST("/order/upload", uploadOrder(or))
	router.GET("/order", getAllOrderByProductID(or))

	// PRODUCT ORDER FOR SELLER DASHBOARD
	router.GET("/seller/product", getProductWithOrder(pr))
	router.GET("/seller/order", getAllOrderBySellerID(or))

	// COUNT FOR PAGINGATE
	router.GET("/order/count", countAllOrderbySellerID(or))

	// UPDATE STATUS
	router.POST("/order/status", updateOrderStatus(or))

	return router, nil
}
