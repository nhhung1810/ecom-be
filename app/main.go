package main

import (
	"app/authentication"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("Hello pong")
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

	r.POST("/login", func(c *gin.Context) {
		var info authentication.Credential
		var a authentication.AuthHandle
		info = a.ParseCredential(c)
		// if err := c.BindJSON(&info); err != nil {
		// 	return
		// }

		fmt.Println(info.Username)
		fmt.Println(info.Password)
		c.IndentedJSON(http.StatusCreated, info.Username)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// type Credential struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }
