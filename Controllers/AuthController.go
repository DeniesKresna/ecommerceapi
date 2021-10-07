package Controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/DeniesKresna/ecommerceapi/Configs"
	"github.com/DeniesKresna/ecommerceapi/Helpers"
	"github.com/DeniesKresna/ecommerceapi/Models"
	"github.com/DeniesKresna/ecommerceapi/Response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

type Auth struct {
	User      *Models.User
	TokenData string
}

func AuthLogin(c *gin.Context) {
	var err error

	apiSecret := os.Getenv("API_SECRET")

	var userLoginInput Models.UserLogin
	c.ShouldBindJSON(&userLoginInput)
	fmt.Print(userLoginInput.Email)
	v := validate.Struct(userLoginInput)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	var user Models.User
	err = Configs.DB.Preload("Role").Where("email = ?", userLoginInput.Email).First(&user).Error
	if err != nil {
		Response.Json(c, 404, "wrong email")
		return
	}

	err = Helpers.VerifyPassword(user.Password, userLoginInput.Password)
	if err != nil {
		Response.Json(c, 404, "wrong password")
		return
	}

	// Create the Claims

	exTime := time.Now().Add(2 * time.Hour).Unix()

	// Create the Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo":    "bar",
		"exp":    exTime,
		"userId": user.ID,
	})
	tokenString, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		Response.Json(c, 450, "cant create token")
		return
	}

	auth := Auth{&user, tokenString}

	Response.Json(c, 200, auth)
}
