package routes

import (
	"ProductManagement/controllers/userController"
	"ProductManagement/helpers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UserRoute(e *echo.Echo) {
	e.GET("/user", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Welcome to main page")
	})
	e.POST("/user/register", userController.Register)
	e.POST("/user/login", userController.Login)

	r := e.Group("/user/restricted")
	r.Use(echojwt.WithConfig(helpers.Config))

	r.GET("/viewProfile/:username", userController.ViewProfile)
	r.PUT("/editProfile/:username", userController.EditAProfile)
	r.DELETE("/deleteProfile/:username", userController.DeleteAUser)
	r.POST("/logout/:username", userController.Logout)
}
