package helpers

import (
	"ProductManagement/configs"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
	ID       string `json:"id"`
	jwt.RegisteredClaims
}

func GetClaimsFromJwt(c echo.Context) *JwtCustomClaims {
	userFromToken := c.Get("user").(*jwt.Token)
	//fmt.Println(userFromToken)
	return userFromToken.Claims.(*JwtCustomClaims)
}

var secretKey = configs.SecretKey()

var Config = echojwt.Config{
	//not necessary to understand
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(JwtCustomClaims)
	},
	SigningKey: []byte(secretKey),
}
