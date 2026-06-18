package controllers
 
import (
	"net/http"

	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"
	"ProjetFilViolet/backend/api/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllUsers(c echo.Context) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

func DeleteUserById(c echo.Context) error {
	id := c.Param("id") // correspond à :id dans la route

	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is required"})
	}

	if err := services.DeleteUserByID(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete user"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}

func GetUserById(c echo.Context) error {
	id := c.Param("id")

	user, err := services.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch user"})
	}
	return c.JSON(http.StatusOK, user)
}

func GetMe(c echo.Context) error {
	user := c.Get("user").(*models.User)
	var freshUser models.User
	if err := config.DB.First(&freshUser, user.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to refresh user profile"})
	}
	freshUser.Password = ""
	return c.JSON(http.StatusOK, freshUser)
}
