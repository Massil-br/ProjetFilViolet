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
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://192.168.1.44:3000", // ajoute ici l'IP + port de ton frontend accessible sur le réseau local
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	routes.InitRoutes(e)
	routes.InitAuthRoutes(e)
	routes.InitUserRoutes(e)
	routes.InitFriendRoutes(e)
	routes.InitSaloonRoutes(e)

	e.Logger.Fatal(e.Start(":8081"))

}
