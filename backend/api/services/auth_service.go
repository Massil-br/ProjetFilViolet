package services

import (
	"errors"
	"os"
	"time"

	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"
	"ProjetFilViolet/backend/api/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	ErrEmailInUse           = errors.New("email already in use")
	ErrNicknameInUse        = errors.New("nickname already in use")
	ErrFailedRemovePrevious = errors.New("failed to remove previous deleted user")
)

// CreateUser creates a new user. It handles DB checks, deletes previous soft-deleted
// accounts with the same email, hashes the password and persists the user.
func CreateUser(firstName, lastName, nickName, email, password string) (*models.User, error) {
	// Check existing nickname
	var existingNickname models.User
	if err := config.DB.Where("nick_name = ?", nickName).First(&existingNickname).Error; err == nil {
		return nil, ErrNicknameInUse
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Check existing user (including soft deleted)
	var existing models.User
	if err := config.DB.Unscoped().Where("email = ?", email).First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			if err := config.DB.Unscoped().Delete(&existing).Error; err != nil {
				return nil, ErrFailedRemovePrevious
			}
		} else {
			return nil, ErrEmailInUse
		}
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Hash password
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName:   firstName,
		LastName:    lastName,
		NickName:    nickName,
		Email:       email,
		Password:    hashed,
		Role:        models.UserRole,
		Money:       config.InitialUserMoney,
		IsVerified:  false,
		Status:      models.MainMenu,
		IsConnected: false,
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func ConnectUser(userID uint) error {
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_connected", true).Error; err != nil {
		return err
	}
	return nil
}

func DisconnectUser(userID uint) error {
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_connected", false).Error; err != nil {
		return err
	}
	return nil
}

// GenerateTokenForUserID returns a signed JWT containing only the user id.
func GenerateTokenForUserID(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

// GetUserByEmail returns a user by email.
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
