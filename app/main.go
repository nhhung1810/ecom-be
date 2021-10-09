package main

import (
	"ecom-be/app/controllers/authHandle"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	routing()
}

func routing() {
	r := gin.Default()

	// The routing
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/register", authHandle.RegisterHandle)
	r.POST("/login", authHandle.LoginHandle)
	// r.POST("/logout", authHandle.LoginHandle)
	r.GET("/user", authHandle.UserHandle)

	r.Run()
}

func setEnvironment(key string, value string) {
	os.Setenv(key, value)
	// HOW TO GET THE ENV VARIABLE
	// os.Getenv(Key)
}
