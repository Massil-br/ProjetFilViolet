package main

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	config.Init()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	routes.InitRoutes(e)
	routes.InitAuthRoutes(e)
	routes.InitUserRoutes(e)
	routes.InitFriendRoutes(e)
	routes.InitSaloonRoutes(e)

	e.Logger.Fatal(e.Start(":8081"))

}
