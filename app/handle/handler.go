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

func Handler(auth auth.Service, img myimg.Service, productService product.Service, orderService order.Service) (*gin.Engine, error) {
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
	router.GET("/check", userHandle(auth))

	// IMAGE HANDLE
	// router.GET("/image", imageHandle())
	router.GET("/image", getImageHandle(img))
	router.POST("/image/upload", newUploadImage(img))

	// PROD HANDLE
	router.GET("/product", getAllProducts(productService))
	router.GET("/product/q", getProductWithFilter(productService))
	router.GET("/product/random", getRandomProduct(productService))
	router.GET("/product/search", searchProductByName(productService))

	// ONLY FOR QUERY USE
	router.GET("/product/info", getProductByID(productService))
	router.POST("/product/upload", uploadProduct(productService))

	// ORDER
	router.POST("/order/upload", uploadOrder(orderService))
	router.GET("/order", getAllOrderByProductID(orderService))

	// PRODUCT ORDER FOR SELLER DASHBOARD
	router.GET("/seller/product", getProductWithOrder(productService))
	router.GET("/seller/order", getAllOrderBySellerID(orderService))

	// COUNT FOR PAGINGATE
	router.GET("/order/count", countAllOrderbySellerID(orderService))
	router.GET("/product/count", countAllProductBySellerID(productService))
	router.GET("/product/count/all", countAllProducts(productService))

	// UPDATE STATUS
	router.PATCH("/order/status", updateOrderStatus(orderService))

	// ARCHIVE PRODUCTS
	router.PATCH("/product/archive", archiveProduct(productService))

	return router, nil
}
