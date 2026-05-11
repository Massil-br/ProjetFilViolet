package controllers

import (
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/Massil-br/GlobalWebsite/backend/config"
	"github.com/Massil-br/GlobalWebsite/backend/models"
	"github.com/Massil-br/GlobalWebsite/backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Role            string
}

func CreateUser(c echo.Context) error {
	var req CreateUserRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if req.ConfirmPassword != req.Password {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Password and confirmPassword are not the same",
		})
	}

	if len(req.Password) < 8 || !regexp.MustCompile("[0-9]").MatchString(req.Password) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Password must be at least 8 chars and contain a digit"})
	}

	var existingUser models.User
	// Recherche user même soft deleted
	err = config.DB.Unscoped().Where("email = ?", req.Email).First(&existingUser).Error

	if err == nil {
		// User existe (même soft deleted)
		if existingUser.DeletedAt.Valid {
			// Soft deleted => suppression hard avant recréation
			if err := config.DB.Unscoped().Delete(&existingUser).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to permanently delete user"})
			}
		} else {
			// User actif => refus création
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email already in use"})
		}
	} else if err != gorm.ErrRecordNotFound {
		// Erreur DB autre
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfully",
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c echo.Context) error {
	var req LoginRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	var user models.User

	err = config.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid creadentials"})
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid crendentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token", })
		
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   tokenString,
	})

}
