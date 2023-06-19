package adminHelpers

import (
	"ProductManagement/helpers"
	"ProductManagement/helpers/userHelpers"
	"ProductManagement/responses"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Register(c echo.Context) error {
	claims := helpers.GetClaimsFromJwt(c)
	isAdmin := claims.Admin

	if isAdmin != true {
		return c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "You are not an Admin!", Data: &echo.Map{"data": ""}})
	}

	return userHelpers.Register(c, true)
}
