package adminController

import (
	"ProductManagement/helpers/adminHelpers"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	return adminHelpers.Register(c)
}
