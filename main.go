package main

import (
	"ProductManagement/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	routes.UserRoute(e)
	routes.AdminRoute(e)

	e.Logger.Fatal(e.Start(":6000"))
}
