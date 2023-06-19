package routes

import (
	"ProductManagement/controllers/productController"
	"ProductManagement/helpers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SellerRoute(e *echo.Echo) {
	e.GET("/product/viewProducts", productController.ViewAllProducts)
	e.GET("/product/viewProducts/:sellerUsername", productController.ViewOwn)

	s := e.Group("/product")
	s.Use(echojwt.WithConfig(helpers.Config))

	s.POST("/createProduct", productController.CreateProduct)
	s.POST("/buyProduct", productController.BuyProduct)
}
