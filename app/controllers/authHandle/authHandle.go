package authHandle

import (
	"ecom-be/app/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "EcaLf2vYAe1GtT369eD6jtfxA0iXC6HlPj1meCE/oro="

// REGISTER HANDLE
func RegisterHandle(c *gin.Context) {
	newAccount := models.ParseCredential(c)

	_, err := newAccount.CheckExisted()
	if err == nil {
		c.IndentedJSON(http.StatusBadRequest, "This account is existed!")
		return
	}

	newAccount.AddNewAccount()
	fmt.Println("Count the account: ", models.CountDB())
	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}

// LOGIN HANDLE
func LoginHandle(c *gin.Context) {
	existedAccount := models.ParseCredential(c)

	user, err := existedAccount.CheckExisted()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Can not find this user!",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(existedAccount.Password)); err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
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
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Sign in fail. Can not login!",
		})
		return
	}

	c.SetCookie("jwt", token, int(expTime.Unix()), "/", "localhost", false, true)

	c.IndentedJSON(http.StatusAccepted, gin.H{
		"message": "accept",
	})
}

func UserHandle(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "can not find cookies",
		})
		return
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	// TODO: Fix this
	if err != nil {
		fmt.Println(err) //json: cannot unmarshal object into Go value of type jwt.Claims
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	// TODO: connect db and return data
	id, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{
			"message": "id is not integer",
		})
		return
	}

	println("id: ", id)
	user, err := models.FindUser(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "can not find user",
		})
		return
	}

	fmt.Println(user)
	c.IndentedJSON(http.StatusAccepted, gin.H{
		"message": "success",
	})
}

func LogoutHandle(c *gin.Context) {
	c.SetCookie("jwt", "", int(time.Now().Add(-time.Hour).Unix()), "/", "localhost", false, true)
	c.IndentedJSON(http.StatusAccepted, gin.H{
		"message": "success",
	})
}
