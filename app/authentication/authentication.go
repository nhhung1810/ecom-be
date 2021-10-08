package authentication

import "github.com/gin-gonic/gin"

type AuthHandle struct {
}

func (a AuthHandle) ParseCredential(c *gin.Context) Credential {
	var info Credential
	if err := c.BindJSON(&info); err != nil {
		return info
	}
	return info
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
