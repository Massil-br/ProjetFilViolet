package routes

import (
	"ProjetFilViolet/backend/api/controllers"
	"ProjetFilViolet/backend/api/middleware"
	"ProjetFilViolet/backend/api/models"

	"github.com/labstack/echo/v4"
)

func InitFriendRoutes(e *echo.Echo) {
	e.GET("/api/friends", controllers.GetFriends, middleware.AuthMiddleware(models.UserRole))
	e.POST("/api/friends/:id", controllers.AddFriendById, middleware.AuthMiddleware(models.UserRole))
	e.POST("/api/friends/:name", controllers.AddFriendByName, middleware.AuthMiddleware(models.UserRole))
	e.DELETE("/api/friends/:id", controllers.RemoveFriend, middleware.AuthMiddleware(models.UserRole))
}
