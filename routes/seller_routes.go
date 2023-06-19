package routes

import (
	"ProductManagement/controllers/productController"
	"ProductManagement/helpers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SellerRoute(e *echo.Echo) {
	s := e.Group("/seller")
	s.Use(echojwt.WithConfig(helpers.Config))

	s.POST("/createProduct", productController.CreateProduct)
	s.POST("/buyProduct", productController.BuyProduct)
}
