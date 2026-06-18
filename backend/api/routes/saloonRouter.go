package routes

import (
	"ProjetFilViolet/backend/api/controllers"
	"ProjetFilViolet/backend/api/middleware"
	"ProjetFilViolet/backend/api/models"

	"github.com/labstack/echo/v4"
)

func InitSaloonRoutes(e *echo.Echo) {
	e.GET("/api/saloons", controllers.GetAllSaloons, middleware.AuthMiddleware(models.UserRole))
	e.GET("/api/saloons/:id", controllers.GetSaloonByID, middleware.AuthMiddleware(models.UserRole))
	e.GET("/api/saloons/name/:name", controllers.GetSaloonByName, middleware.AuthMiddleware(models.UserRole))
	e.POST("/api/saloons", controllers.CreateSaloon, middleware.AuthMiddleware(models.Moderator))	
	e.PUT("/api/saloons/:id", controllers.UpdateSaloon, middleware.AuthMiddleware(models.Moderator))
	e.DELETE("/api/saloons/:id", controllers.DeleteSaloonByID, middleware.AuthMiddleware(models.Moderator))
	e.DELETE("/api/saloons/name/:name", controllers.DeleteSaloonByName, middleware.AuthMiddleware(models.Moderator))
	e.GET("/api/tables", controllers.GetAllTables, middleware.AuthMiddleware(models.UserRole))
}
