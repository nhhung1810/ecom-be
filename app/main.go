package main

import (
	"ecom-be/app/controllers/authHandle"
	"ecom-be/app/models"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	models.Connect()
	routing()
}

func routing() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// The routing
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/register", authHandle.RegisterHandle)
	r.POST("/login", authHandle.LoginHandle)
	r.POST("/logout", authHandle.LogoutHandle)
	r.GET("/user", authHandle.UserHandle)

	r.Run()
}

func setEnvironment(key string, value string) {
	os.Setenv(key, value)
	// HOW TO GET THE ENV VARIABLE
	// os.Getenv(Key)
}
