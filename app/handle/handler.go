package handle

import (
	"ecom-be/app/auth"
	myimg "ecom-be/app/image"
	"ecom-be/app/product"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler(auth auth.Service, img myimg.Service, pr product.Service) (*gin.Engine, error) {
	// TODO: Handle error here
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

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
	router.GET("/image", imageHandle())
	router.GET("/image/:id", getImageHandle(img))
	router.POST("/upload/image", uploadImageHandle(img))

	// PROD HANDLE
	router.GET("/product", getAllProducts(pr))
	router.POST("/upload/product", uploadProduct(pr))

	return router, nil
}
