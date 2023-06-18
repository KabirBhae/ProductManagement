package routes

import (
	"ProductManagement/controllers/adminController"
	"ProductManagement/helpers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AdminRoute(e *echo.Echo) {
	a := e.Group("/admin")
	a.Use(echojwt.WithConfig(helpers.Config))

	a.POST("/register", adminController.Register)

}
