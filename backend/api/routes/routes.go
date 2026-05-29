package routes

import (
	"ProjetFilViolet/backend/api/controllers"
	"ProjetFilViolet/backend/api/middleware"
	"ProjetFilViolet/backend/api/models"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	InitGetRoutes(e)
	InitDeleteRoutes(e)
}

func InitGetRoutes(e *echo.Echo) {
	e.GET("/api", controllers.MainPage)
	e.GET("/api/users", controllers.GetAllUsers, middleware.AuthMiddleware(models.Admin))
	e.GET("/api/users/:id", controllers.GetUserById, middleware.AuthMiddleware(models.Admin))

	e.GET("/api/logged", controllers.LoggedTest,

		middleware.AuthMiddleware(models.UserRole),
	)

}



func InitDeleteRoutes(e *echo.Echo) {
	e.DELETE("/api/users/:id", controllers.DeleteUserById)
}
