package routes

import (
	"ProjetFilViolet/backend/api/controllers"
	"ProjetFilViolet/backend/api/middleware"
	"ProjetFilViolet/backend/api/models"

	"github.com/labstack/echo/v4"
)

func InitAuthRoutes(e *echo.Echo) {
	e.POST("/api/register", controllers.CreateUser)
	e.POST("/api/login", controllers.Login)
	e.POST("/api/logout", controllers.Logout, middleware.AuthMiddleware(models.UserRole))
}
