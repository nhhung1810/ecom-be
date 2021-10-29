package handle

import (
	"ecom-be/app/auth"
	"ecom-be/app/config"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = config.DefaultConfig.SecretKey

var successMsg = gin.H{
	"message": "success",
}

func registerHandle(auth auth.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		newAccount, err := auth.ParseCredential(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		_, err = auth.CheckExisted(*newAccount)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "account existed",
			})
			return
		}

		err = auth.AddAccount(*newAccount)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, successMsg)
	}
}

func loginHandle(auth auth.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		reqAccount, err := auth.ParseCredential(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}

		user, err := auth.CheckExisted(*reqAccount)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqAccount.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authentication failed. Please check your password!",
			})
			return
		}

		expTime := time.Now().Add(time.Hour * 1)

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    strconv.Itoa(user.ID),
			ExpiresAt: expTime.Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Sign in fail. Can not login!",
			})
			return
		}

		c.SetCookie("jwt", token, int(expTime.Unix()), "/", "localhost", false, true)

		c.JSON(http.StatusOK, successMsg)
	}
}

func userHandle(auth auth.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := cookieAuth(c)

		_, err = auth.FindUserByID(*id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, successMsg)
	}
}

func logoutHandle(auth auth.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.SetCookie("jwt", "", int(time.Now().Add(-time.Hour).Unix()), "/", "localhost", false, true)
		c.JSON(http.StatusOK, successMsg)
	}
}

func cookieAuth(c *gin.Context) (*int, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not find cookies",
		})
		return nil, err
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	id, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": err.Error(),
		})
		return nil, err
	}

	return &id, nil
}
