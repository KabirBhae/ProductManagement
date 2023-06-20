package productController

import (
	"ProductManagement/helpers/productHelpers"
	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	return productHelpers.CreateProduct(c)
}

func BuyProduct(c echo.Context) error {
	return productHelpers.BuyProduct(c)
}
func ViewAllProducts(c echo.Context) error {
	return productHelpers.ViewAllProducts(c)
}

func ViewOwn(c echo.Context) error {
	return productHelpers.ViewOwn(c)
}

func DeleteProduct(c echo.Context) error {
	return productHelpers.DeleteProduct(c)
}
