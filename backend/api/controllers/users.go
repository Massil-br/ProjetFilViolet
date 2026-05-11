package controllers

import (
	"net/http"

	"github.com/Massil-br/GlobalWebsite/backend/config"
	"github.com/Massil-br/GlobalWebsite/backend/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllUsers(c echo.Context) error {
	var users []models.User
	err := config.DB.Find(&users).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Faild to fetch users",
		})
	}
	return c.JSON(http.StatusOK, users)
}

func DeleteUserById(c echo.Context) error {
	id := c.Param("id") // correspond à :id dans la route

	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is required"})
	}

	result := config.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete user"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}

func GetUserById(c echo.Context) error {
	id := c.Param("id")

	var user models.User

	err := config.DB.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch user"})
	}

	return c.JSON(http.StatusOK, user)
}
