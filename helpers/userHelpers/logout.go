package userHelpers

import (
	"ProductManagement/helpers"
	"ProductManagement/responses"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Logout(c echo.Context) error {
	//fmt.Println(userFromToken)
	claims := helpers.GetClaimsFromJwt(c)
	usernameFromClaims := claims.Username

	usernamefromURL := c.Param("username")

	if usernamefromURL != usernameFromClaims {
		return c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "unauthorised user", Data: &echo.Map{"data": ""}})
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1e9))
	//fmt.Println("logged out")

	return c.Redirect(http.StatusSeeOther, "/user/restricted/viewProfile/"+usernamefromURL)
}
