package userController

import (
	"ProductManagement/helpers/userHelpers"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	return userHelpers.Register(c, false)
}

func Login(c echo.Context) error {
	return userHelpers.Login(c)
}

func ViewProfile(c echo.Context) error {
	return userHelpers.ViewProfile(c)
}

func EditAProfile(c echo.Context) error {
	return userHelpers.EditAProfile(c)
}

func DeleteAUser(c echo.Context) error {
	return userHelpers.DeleteAUser(c)
}

func Logout(c echo.Context) error {
	return userHelpers.Logout(c)
}
