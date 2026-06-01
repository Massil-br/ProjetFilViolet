package routes

import (
	"ProjetFilViolet/backend/api/controllers"
	"ProjetFilViolet/backend/api/middleware"
	"ProjetFilViolet/backend/api/models"

	"github.com/labstack/echo/v4"
)

func InitUserRoutes(e *echo.Echo) {
	e.GET("/api/users", controllers.GetAllUsers, middleware.AuthMiddleware(models.Moderator))
	e.GET("/api/users/:id", controllers.GetUserById, middleware.AuthMiddleware(models.Moderator))
	e.DELETE("/api/users/:id", controllers.DeleteUserById, middleware.AuthMiddleware(models.Admin))
}
