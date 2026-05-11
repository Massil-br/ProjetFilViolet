package routes

import (
	"github.com/Massil-br/GlobalWebsite/backend/controllers"
	"github.com/Massil-br/GlobalWebsite/backend/middleware"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	InitGetRoutes(e)
	InitPostRoutes(e)
	InitDeleteRoutes(e)
	InitPutRoutes(e)

}

func InitGetRoutes(e *echo.Echo) {
	e.GET("/api", controllers.MainPage)
	e.GET("/api/users", controllers.GetAllUsers, middleware.AuthMiddleware("admin"))
	e.GET("/api/users/:id", controllers.GetUserById, middleware.AuthMiddleware("admin"))

	e.GET("/api/logged", controllers.LoggedTest,

		middleware.AuthMiddleware("user"),
	)




}

func InitPutRoutes(e *echo.Echo) {



}

func InitPostRoutes(e *echo.Echo) {

}

func InitDeleteRoutes(e *echo.Echo) {
	e.DELETE("/api/users/:id", controllers.DeleteUserById)
}
